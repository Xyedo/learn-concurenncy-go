package main

import (
	"database/sql"
	"encoding/gob"
	"final-project/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = 80

func main() {

	db := initDB()
	db.Ping()

	session := initSession()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	wg := sync.WaitGroup{}

	app := Config{
		Session:       session,
		DB:            db,
		Wait:          &wg,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Models:        data.New(db),
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
	}
	app.Mailer = app.createMail()
	go app.listenForMail()
	go app.listenForErrors()
	go app.listenForShutdown()
	app.serve()

}

func (app *Config) listenForErrors() {
	for {
		select {
		case err, ok := <-app.ErrorChan:
			if ok {
				app.ErrorLog.Println(err)
			}
		case <-app.ErrorChanDone:
			return

		}
	}
}

func (app *Config) createMail() Mail {
	return Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  None,
		FromName:    "Info",
		FromAddress: "info@mycompany.com",
		Wait:        app.Wait,
		ErrorChan:   make(chan error),
		MailerChan:  make(chan Message, 100),
		DoneChan:    make(chan bool),
	}

}
func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
func initSession() *scs.SessionManager {
	gob.Register(data.User{})
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	return session
}
func initRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}

}

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("cant connect to database")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postges not yet ready...", err)

		} else {
			log.Print("connected to database")
			return connection
		}
		if counts > 10 {
			return nil
		}
		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	app.InfoLog.Println("would run cleanup tasks...")
	app.Wait.Wait()

	app.Mailer.DoneChan <- true
	app.ErrorChanDone <- true
	app.InfoLog.Println("closing channels and shutting down app")
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)
	close(app.ErrorChan)
	close(app.ErrorChanDone)
}

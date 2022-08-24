package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type app struct {
	seatingCapacity int
	arrivalRate     int
	cutDuration     time.Duration
	timeOpen        time.Duration
}

func main() {
	app := &app{
		seatingCapacity: 10,
		arrivalRate:     100,
		cutDuration:     1000 * time.Millisecond,
		timeOpen:        10 * time.Second,
	}
	rand.Seed(time.Now().UnixNano())

	color.Yellow("The sleeping Barber Problem")
	color.Yellow("---------------------------")
	clientChan := make(chan string, app.seatingCapacity)
	doneChan := make(chan bool)

	saloon := &barberShop{
		shopCapacity:    app.seatingCapacity,
		HairCutDuration: app.cutDuration,
		numberOfBarbers: 0,
		clientsChan:     clientChan,
		barbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day")

	saloon.addBarber("Frank")
	saloon.addBarber("Hafid")
	saloon.addBarber("Tanri")
	saloon.addBarber("Golang")
	saloon.addBarber("Is")
	saloon.addBarber("Cool")

	saloonClosing := make(chan bool)
	closed := make(chan bool)
	go func() {
		<-time.After(app.timeOpen)
		saloonClosing <- true
		saloon.closeShopForDay()
		closed <- true
	}()

	clientNum := 1
	go func() {
		for {
			randomMilis := rand.Int() % (2 * app.arrivalRate)
			select {
			case <-saloonClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilis)):
				saloon.addClient(fmt.Sprintf("Client#%d", clientNum))
				clientNum++

			}
		}
	}()
	<-closed
}

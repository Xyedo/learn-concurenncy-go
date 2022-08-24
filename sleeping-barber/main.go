package main

import (
	"fmt"
	"math/rand"
	"time"
)

type app struct {
	seatNum         int
	durrOpen        time.Duration
	barberWork      time.Duration
	custArivingRate int
	saloonClosing   chan bool
	canFinish       chan bool
}

func main() {
	rand.Seed(time.Now().UnixNano())
	app := &app{
		seatNum:         10,
		durrOpen:        10 * time.Second,
		barberWork:      1000 * time.Millisecond,
		custArivingRate: 100,
		saloonClosing:   make(chan bool),
		canFinish:       make(chan bool),
	}
	fmt.Println("Sleeping Barber Problem")
	fmt.Println("-----------------------")

	defer close(app.saloonClosing)
	defer close(app.canFinish)
	orders := make(chan int, app.seatNum)
	barbershop := &saloon{
		barberWork:   app.barberWork,
		numBarber:    0,
		barberFinish: make(chan bool),
	}
	barbershop.addBarber("Hafid", orders)
	barbershop.addBarber("Tanri", orders)
	barbershop.addBarber("Ray", orders)
	barbershop.addBarber("Golang", orders)
	go func() {
		<-time.After(app.durrOpen)
		app.saloonClosing <- true
		close(orders)
		barbershop.barberWorking()
		app.canFinish <- true
	}()
	go app.customer(orders)

	<-app.canFinish
}

func (a *app) customer(orders chan<- int) {
	orderNum := 1
	for {
		randomMilis := rand.Int() % (2 * a.custArivingRate)
		<-time.After(time.Millisecond * time.Duration(randomMilis))
		fmt.Println("Customer Coming!")
		select {
		case orders <- orderNum:

			fmt.Println("there is a seat, customer order ", orderNum)

			orderNum++

		case isClose, ok := <-a.saloonClosing:
			if ok && isClose {
				fmt.Println("Saloon closed, customer leaving")
				return
			}
		default:
			fmt.Println("Saloon is full, customer leaving with anger")

		}

	}

}

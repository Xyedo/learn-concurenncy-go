package main

import (
	"time"

	"github.com/fatih/color"
)

type barberShop struct {
	shopCapacity    int
	HairCutDuration time.Duration
	numberOfBarbers int
	barbersDoneChan chan bool
	clientsChan     chan string
	Open            bool
}

func (shop *barberShop) addBarber(barber string) {
	shop.numberOfBarbers++

	go func() {
		isSleeping := false

		color.Yellow("%s goes to the waiting room to check for clients.\n", barber)

		for {
			if len(shop.clientsChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap", barber)
				isSleeping = true
			}
			client, shopOpen := <-shop.clientsChan

			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up", client, barber)
					isSleeping = false
				}
				shop.cutHair(barber, client)

			} else {
				shop.sendBarberHome(barber)
				return
			}

		}
	}()
}

func (shop *barberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair\n", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finish cutting %s's hair\n", barber, client)
}
func (shop *barberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home\n", barber)
	shop.barbersDoneChan <- true
}

func (shop *barberShop) closeShopForDay() {
	color.Cyan("Closing shop for the day")

	close(shop.clientsChan)
	shop.Open = false

	for a := 1; a <= shop.numberOfBarbers; a++ {
		<-shop.barbersDoneChan
	}
	close(shop.barbersDoneChan)

	color.Green("---------------------------------------------------------------------")
	color.Green("The barbershop is now closed for the day, and everyone has gone home")
}

func (shop *barberShop) addClient(client string) {
	color.Green("*** %s arrives!", client)

	if shop.Open {
		select {
		case shop.clientsChan <- client:
			color.Yellow("%s takes a seat in a waiting room", client)
		default:
			color.Red("The waiting room is full, so %s leaves angrily😠", client)
		}

	} else {
		color.Red("The shop is already closed ", client)
	}
}

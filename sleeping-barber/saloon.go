package main

import (
	"fmt"
	"time"
)

type saloon struct {
	barberWork   time.Duration
	barberFinish chan bool
	numBarber    int
}

func (saloon *saloon) addBarber(barber string, orders <-chan int) {
	saloon.numBarber++
	go func() {
		isSleeping := false
		for {
			if len(orders) == 0 {
				fmt.Printf("There is nothing to do, so %s takes a nap\n", barber)
				isSleeping = true
			}

			cust, ok := <-orders
			if ok {
				if isSleeping {
					fmt.Printf("cust no #%d wakes %s up\n", cust, barber)
					isSleeping = false
				}
				fmt.Printf("order #%d\n handled by %s\n", cust, barber)
				<-time.After(saloon.barberWork)
				fmt.Printf("order #%d finishes\n", cust)
			} else {
				fmt.Printf("the work of %s is finish, %s going home\n", barber, barber)
				saloon.barberFinish <- true
				return
			}

		}
	}()

}

func (saloon *saloon) barberWorking() {
	for i := 0; i < saloon.numBarber; i++ {
		<-saloon.barberFinish
	}
	close(saloon.barberFinish)
	fmt.Println("all barber going home! app finishes")

}

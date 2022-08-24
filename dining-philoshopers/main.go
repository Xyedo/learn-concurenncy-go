package main

import (
	"fmt"
	"sync"
	"time"
)

var philosopher = []string{
	"Plato",
	"Socrates",
	"Aristotle",
	"Pascal",
	"Locke",
}
var wg sync.WaitGroup

var sleepTime = 1 * time.Second
var hunger = 3

var eatTime = 3 * time.Second
var thinkTime = 1 * time.Second
var order = make([]string, 0, len(philosopher))

func main() {
	fmt.Println("The dining philosophers Problem")
	fmt.Println("--------------------------------")
	wg.Add(len(philosopher))

	finished := make(chan string)

	forkLeft := &sync.Mutex{}
	for _, v := range philosopher {
		forkRight := &sync.Mutex{}
		go diningProblem(v, forkLeft, forkRight, finished)
		forkLeft = forkRight

	}
	go func() {

		wg.Wait()
		close(finished)
	}()

	fmt.Println("The table is empty")

	for sufi := range finished {
		fmt.Printf("%s is finished eating\n", sufi)
		order = append(order, sufi)
	}
	fmt.Println("the order of who finishes is", order)

}

func diningProblem(sufi string, dominantHand, otherHand *sync.Mutex, finished chan string) {
	defer wg.Done()

	fmt.Println(sufi, "is seated")
	time.Sleep(sleepTime)
	for i := hunger; i > 0; i-- {
		fmt.Printf("%s is hungry\n", sufi)
		time.Sleep(sleepTime)
		dominantHand.Lock()
		fmt.Printf("\t%s is getting a left fork\n", sufi)
		otherHand.Lock()
		fmt.Printf("\t%s is getting a right fork\n", sufi)

		fmt.Printf("%s has both fork and now eating\n", sufi)
		time.Sleep(eatTime)

		fmt.Printf("%s is thinking\n", sufi)
		time.Sleep(thinkTime)

		dominantHand.Unlock()
		fmt.Printf("\t%s put down the fork on his right\n", sufi)
		otherHand.Unlock()
		fmt.Printf("\t%s put down the fork on his left\n", sufi)
		time.Sleep(sleepTime)
	}

	fmt.Printf("%s is satisfied.\n", sufi)
	time.Sleep(sleepTime)
	fmt.Printf("%s has left the table.\n", sufi)
	finished <- sufi

}

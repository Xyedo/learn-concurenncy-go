package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Opening a program...")
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	var userInput string
	for {
		fmt.Println("type any word and Press enter or press q to quit")
		fmt.Scanln(&userInput)
		if strings.Compare(strings.ToLower(userInput), "q") == 0 {
			break
		}
		ping <- userInput
		fmt.Println(<-pong)
	}
	fmt.Println("program end. closing a chanels")
	close(ping)
	close(pong)

}

func shout(ping <-chan string, pong chan<- string) {
	for {
		shout := <-ping
		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(shout))
	}
}

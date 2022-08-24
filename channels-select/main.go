package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Select with channels")
	fmt.Println("---------------------")
	chan1 := make(chan string)
	chan2 := make(chan string)
	go server1(chan1)
	go server2(chan2)

	for {
		select {
		case s1 := <-chan1:
			fmt.Println("Case one:", s1)
		case s2 := <-chan1:
			fmt.Println("Case two:", s2)
		case s3 := <-chan2:
			fmt.Println("Case three:", s3)
		case s4 := <-chan2:
			fmt.Println("case four:", s4)

		}
	}
}
func server1(ch chan<- string) {
	for {
		time.Sleep(6 * time.Second)
		ch <- "This is from server 1"
	}
}

func server2(ch chan<- string) {
	for {
		time.Sleep(3 * time.Second)
		ch <- "THis is from server 2"
	}
}

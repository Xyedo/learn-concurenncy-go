package main

import (
	"fmt"
	"time"
)

var philosophers = []string{
	"Plato",
	"Socrates",
	"Aristotle",
	"Pascal",
	"Locke",
}

type App struct {
	Hunger    int
	EatTime   time.Duration
	ThinkTime time.Duration
	SleepTime time.Duration
}

var order = make([]string, 0, len(philosophers))

func main() {
	app := &App{
		Hunger:    3,
		SleepTime: 1 * time.Second,
		EatTime:   3 * time.Second,
		ThinkTime: 1 * time.Second,
	}

	fmt.Println("The dining philosophers Problem")
	fmt.Println("--------------------------------")

	fork := make(chan string, 5)
	finished := make(chan string)
	for k, philo := range philosophers {
		go diningTable(philo, fork, finished)
	}

}

func (a *App) diningTable(philo string, fork chan string, finished chan string) {
	fmt.Println(philo, "is seated")
	time.Sleep(sleepTime)
}

package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func main() {
	var bankBalance int
	var balance sync.Mutex
	fmt.Printf("initial account balance :%d.00", bankBalance)
	fmt.Println()

	incomes := []Income{
		{Source: "Main job", Amount: 500},
		{Source: "Gifts", Amount: 10},
		{Source: "Part time job", Amount: 50},
		{Source: "Investments", Amount: 100},
	}

	for i, income := range incomes {
		wg.Add(1)
		go func(i int, income Income) {
			defer wg.Done()
			for week := 1; week <= 52; week++ {
				balance.Lock()
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				balance.Unlock()

				fmt.Printf("ON week %d, you have earned $%d.00 from %s\n", week, income.Amount, income.Source)
			}
		}(i, income)
	}
	wg.Wait()

	fmt.Printf("Final bank balance : $%d.00 \n", bankBalance)
}

// var msg string
// var wg sync.WaitGroup

// func updateMessage(s string) {
// 	defer wg.Done()

// 	msg = s

// }
// func main() {
// 	msg = "Hello, world!"
// 	wg.Add(2)
// 	go updateMessage("hello universe!")
// 	go updateMessage("hello cosmos!")
// 	wg.Wait()

// 	fmt.Println(msg)
// }

// func updateMessage(s string, m *sync.Mutex) {
// 	defer wg.Done()

// 	m.Lock()
// 	msg = s
// 	m.Unlock()
// }
// func main() {
// 	msg = "Hello, world!"
// 	var mutex sync.Mutex
// 	wg.Add(2)
// 	go updateMessage("hello universe!", &mutex)
// 	go updateMessage("hello cosmos!", &mutex)
// 	wg.Wait()

// 	fmt.Println(msg)
// }

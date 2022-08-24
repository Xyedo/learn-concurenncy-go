package main

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

func (p *Producer) close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

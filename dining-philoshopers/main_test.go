package main

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	eatTime = 0 * time.Second
	sleepTime = 0 * time.Second
	thinkTime = 0 * time.Second

	for i := 0; i < 100; i++ {
		main()
		if len(order) != 5 {
			t.Errorf("expected to 5 but get %d", len(order))
		}
		order = make([]string, 0, len(philosopher))
	}
}

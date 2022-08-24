package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		m := strconv.Itoa(i)
		wg.Add(1)
		go updateMessage(m, &wg)
		wg.Wait()
		if strings.Compare(m, msg) != 0 {
			t.Errorf("expected %s, yet its not there", m)
		}
	}

}

func Test_printMessage(t *testing.T) {
	msg = "Hello smuanya"
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("pipe error")
	}
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	printMessage()
	w.Close()
	os.Stdout = old
	output := <-outC

	if !strings.Contains(output, msg) {
		t.Errorf("expected %s, yet got %s ", msg, output)
	}

}

func Test_main(t *testing.T) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("pipe error")
	}
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	main()
	w.Close()
	os.Stdout = old
	output := <-outC
	test := []string{
		"Hello, universe!",
		"Hello, cosmos!",
		"Hello, world!",
	}
	for _, v := range test {
		if !strings.Contains(output, v) {
			t.Errorf("expected %s, yet got %s ", v, output)
		}
		output = strings.ReplaceAll(output, v, "")

	}

}

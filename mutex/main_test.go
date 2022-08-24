package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	os.Stdout = w

	// print() // fails if called here "fatal error: all goroutines are asleep - deadlock"

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	main()

	// back to normal state
	err = w.Close()
	if err != nil {
		t.Error(err)
	}
	os.Stdout = old // restoring the real stdout
	out := <-outC
	if !strings.Contains(out, "$34320.00") {
		t.Error("not found")
	}
	// reading our temp stdout

	//fmt.Print(out)

}

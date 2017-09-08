package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {

	fmt.Printf(os.Args[1])

	i, err := strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var (
		pin = rpio.Pin(i)
	)

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	// Toggle pin 20 times
	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second / 5)
	}
}

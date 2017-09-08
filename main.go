package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/stianeikeland/go-rpio"
)

var (
	pin rpio.Pin
)

func main() {

	fmt.Printf(os.Args[1])

	i, err := strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pin = rpio.Pin(i)

	router := chi.NewRouter()
	router.HandleFunc("/toggle", toggleGate)

	http.ListenAndServe(":8080", router)
}

func toggleGate(w http.ResponseWriter, req *http.Request) {

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

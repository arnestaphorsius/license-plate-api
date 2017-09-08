package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
		log.Fatalf("error reading argument from cmd: %s", err.Error())
	}

	// open memory range for GPIO access in /dev/mem
	if err := rpio.Open(); err != nil {
		log.Fatalf("error opening memory range: %s", err.Error())
	}
	// unmap gpio memory when done
	defer rpio.Close()

	// select the pin specified as cli argument
	pin = rpio.Pin(i)

	router := chi.NewRouter()
	router.HandleFunc("/toggle", toggleGate)

	http.ListenAndServe(":8080", router)
}

func toggleGate(w http.ResponseWriter, req *http.Request) {
	// Set pin to output mode
	pin.Output()

	// Toggle the pin on the Raspberry Pi
	pin.Toggle()
}

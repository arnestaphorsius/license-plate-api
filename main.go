package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/stianeikeland/go-rpio"
)

var (
	port int      = 8080
	pin  rpio.Pin = rpio.Pin(4)
)

func main() {
	// read the preferred pin number from the commandline
	i, err := strconv.Atoi(os.Args[1])

	if err != nil {
		log.Fatalf("error reading argument from cmd: %s", err.Error())
	}

	// open memory range for GPIO access in /dev/mem
	if err := rpio.Open(); err != nil {
		log.Fatalf("error opening memory range: %s", err.Error())
	}
	// unmap GPIO memory when done
	defer rpio.Close()

	// select the pin specified
	pin = rpio.Pin(i)

	router := chi.NewRouter()
	router.HandleFunc("/toggle", toggleGate)

	fmt.Printf("starting service listening  on port [%s]", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func toggleGate(w http.ResponseWriter, req *http.Request) {

	go func() {
		// Set pin to output mode
		pin.Output()

		// Set the pin to High to manually open the gate.
		pin.High()

		// Wait three seconds to give the gate some time to adjust.
		time.Sleep(time.Second * 3)

		// Set the pin to Low to put the gate back in it's normal operating mode.
		pin.Low()
	}()
}

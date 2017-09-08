package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/stianeikeland/go-rpio"
)

var (
	port = 8080
	pin  = rpio.Pin(4)
)

var tokenAuth *jwtauth.JwtAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("my_secret_keyy"), nil)
}

func main() {
	// read the preferred pin number from the commandline
	i, err := strconv.Atoi(os.Args[1])

	if err != nil {
		log.Printf("error reading argument from cmd: %s", err.Error())
	} else {
		// select the specified pin
		pin = rpio.Pin(i)
	}

	// open memory range for GPIO access in /dev/mem
	if err := rpio.Open(); err != nil {
		log.Fatalf("error opening memory range: %s", err.Error())
	}
	// unmap GPIO memory when done
	defer rpio.Close()

	router := chi.NewRouter()

	// Protected routes
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/validate", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["beam"])))
		})
	})

	router.HandleFunc("/toggle", toggleGate)

	fmt.Printf("starting service listening  on port [%d]", port)
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

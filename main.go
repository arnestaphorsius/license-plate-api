package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/stianeikeland/go-rpio"
)

var (
	port       = 8080
	pin        = rpio.Pin(4)
	hmacSecret = []byte("my_secret_keyy")
)

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
	router.HandleFunc("/toggle", toggleGate)
	router.HandleFunc("/validate", validateJWT)

	fmt.Printf("starting service listening  on port [%d]", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func valid(req *http.Request) bool {
	headToken := req.Header.Get("Authorization")
	var tokenString string

	if strings.Contains(headToken, "Bearer ") {
		tokenString = headToken[7:]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Only accept HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return hmacSecret, nil
		})

		fmt.Println(err)

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			beam, _ := claims["beam"].(string)
			if beam == "yes" {
				return true
			}
		}
	}
	return false
}

func validateJWT(w http.ResponseWriter, req *http.Request) {
	if valid(req) {
		w.Write([]byte("Proceed"))
	} else {
		w.WriteHeader(401)
		w.Write([]byte("Please authenticate"))
	}
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

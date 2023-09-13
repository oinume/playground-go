package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
)

// output:
// Operation(): n = 1
// Notifiy(): err = error!, d = 3s
// Operation(): n = 2
// Notifiy(): err = error!, d = 7.5s
// Operation(): n = 3
// Notifiy(): err = error!, d = 15s
// Operation(): n = 4
func main() {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.Multiplier = 2.5
	expBackoff.InitialInterval = 3000 * time.Millisecond
	expBackoff.RandomizationFactor = 0.0
	expBackoff.MaxElapsedTime = 15 * time.Second
	expBackoff.MaxInterval = 15 * time.Second

	n := 0
	err := backoff.RetryNotify(func() error {
		time.Sleep(time.Millisecond * 500)
		n++
		fmt.Printf("Operation(): n = %v\n", n)
		return fmt.Errorf("error!")
	}, expBackoff, func(err error, d time.Duration) {
		fmt.Printf("Notifiy(): err = %v, d = %v\n", err, d)
	})
	if err != nil {
		log.Fatal(err)
	}
}

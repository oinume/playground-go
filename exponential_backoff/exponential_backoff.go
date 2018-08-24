package main

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
)

func main() {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.Multiplier = 2.5
	expBackoff.InitialInterval = 3000 * time.Millisecond
	expBackoff.RandomizationFactor = 0.0
	expBackoff.MaxElapsedTime = 20 * time.Second
	expBackoff.MaxInterval = 15 * time.Second
	backoff.RetryNotify(func() error {
		fmt.Printf("Operation()\n")
		return fmt.Errorf("error!")
	}, expBackoff, func(err error, d time.Duration) {
		fmt.Printf("Notifiy(): err = %v, d = %v\n", err, d)
	})
}

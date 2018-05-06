package main

import (
	"crypto/rand"
	"math/big"
	"log"
	"fmt"
)

func main() {
	r, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("random int = %v\n", r.Int64())
}
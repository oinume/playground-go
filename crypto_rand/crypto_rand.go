package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
)

func main() {
	v, err := generateInt(1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("random int = %v\n", v)

	s, err := generateString(10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("random string = %v\n", s)
}

func generateInt(n int64) (int64, error) {
	r, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		return 0, err
	}
	return r.Int64(), nil
}

func generateString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

package main

import (
	"math/rand"
	"time"
	"fmt"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Printf("rand = %v\n", r.Int31())
}

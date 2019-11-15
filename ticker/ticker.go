package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for i := 0; i < 10; i++ {
		select {
		case <-ticker.C:
			doSomething(i)
		}
	}
}

func doSomething(v int) {
	fmt.Printf("%d\n", v)
}

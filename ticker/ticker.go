package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for i := 0; i < 10; i++ {
		//nolint:gosimple
		select {
		case t := <-ticker.C:
			doSomething(i, t)
		}
	}
}

func doSomething(v int, t time.Time) {
	fmt.Printf("%d at %v\n", v, t)
}

package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Printf("i = %02d, fibonacci = %04d\n", i, fibonacci(i))
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("Done")
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// https://qiita.com/sudix/items/67d4cad08fe88dcb9a6d
package main

import (
	"fmt"
	"sync"
)

func main() {
	values := []int{1, 2, 3}
	wg := sync.WaitGroup{}
	for _, v := range values {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			fmt.Printf("v = %d\n", v)
		}(v)
	}
	wg.Wait()
}

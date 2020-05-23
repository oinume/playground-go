package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("Started\n")
	ch := make(chan int, 100)
	wg := new(sync.WaitGroup)
	for i := 0; i < 2; i++ {
		go worker(ctx, i, ch, wg)
	}
	for i := 1; i < 100; i++ {
		ch <- i
	}
	close(ch)
	time.Sleep(10 * time.Second)
	log.Printf("Call parent context cancel\n")
	cancel()
	log.Printf("After Call parent context cancel\n")
	wg.Wait()
}

func worker(ctx context.Context, n int, ch chan int, wg *sync.WaitGroup) {
	defer func() {
		log.Printf("[%d] worker done in defer\n", n) // なぜ呼ばれない？
		wg.Done()
	}()
	//for v := range ch {
	//	func() {
	//		log.Printf("[%d] %d from channel\n", n, v)
	//		_, cancel := context.WithTimeout(ctx, 10*time.Second)
	//		defer func() {
	//			log.Printf("[%d] child context canceled\n", n)
	//			cancel()
	//		}()
	//		time.Sleep(15 * time.Second)
	//		fmt.Printf("[%d] worker done\n", n)
	//	}()
	//}
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			func() {
				log.Printf("[%d] %d from channel\n", n, v)
				_, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer func() {
					log.Printf("[%d] child context canceled\n", n)
					cancel()
				}()
				time.Sleep(15 * time.Second)
				log.Printf("[%d] worker done\n", n)
			}()
		case <-ctx.Done():
			log.Printf("[%d] worker done because of ctx.Done\n", n)
			return
		}
	}
}

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
		wg.Add(1)
		go worker(ctx, i, ch, wg)
	}
	for i := 1; i < 10; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	time.Sleep(7 * time.Second)
	log.Printf("Call parent context cancel()\n")
	cancel()
	log.Printf("After Call parent context cancel\n")
	time.Sleep(3 * time.Second)
}

func worker(ctx context.Context, n int, ch chan int, wg *sync.WaitGroup) {
	defer func() {
		log.Printf("[%d] worker done in defer\n", n)
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
				log.Printf("[%d] worker done because of channel closed\n", n)
				return
			}
			func() {
				log.Printf("[%d] '%d' from channel\n", n, v)
				childCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
				defer func() {
					log.Printf("[%d] Call child context cancel() in defer\n", n)
					cancel()
				}()
				time.Sleep(5 * time.Second)
				log.Printf("[%d] worker done for '%d'\n", n, v)
				select { //nolint:gosimple
				case <-childCtx.Done():
					log.Printf("[%d] worker done because of childCtx.Done: %v\n", n, childCtx.Err())
					return
				}
			}()
		case <-ctx.Done():
			log.Printf("[%d] worker done because of ctx.Done: %v\n", n, ctx.Err())
			return
		}
	}
}

package main

import (
	"fmt"
	"net/http"
)

type result struct {
	url    string
	status int
	err    error
}

func main() {
	urls := make(chan string, 3)
	results := make(chan result, 3)

	// Consumer
	go httpGet(urls, results)

	// Producer
	go func() {
		targetURLs := []string{
			"https://www.lekcije.com/",
			"https://journal.lampetty.net/",
			"https://journal.lampetty.net/entry/concurrency-in-go-goroutines",
			"https://journal.lampetty.net/entry/concurrency-in-go-goroutines",
			"https://journal.lampetty.net/entry/concurrency-in-go-goroutines",
			"https://journal.lampetty.net/entry/concurrency-in-go-goroutines",
			"https://journal.lampetty.net/entry/concurrency-in-go-goroutines",
			"https://journal.lampetty.net/entry/concurrency-in-go-goroutines",
			"https://journal.lampetty.net/entry/concurrency-in-go-goroutines",
		}
		for _, url := range targetURLs {
			urls <- url
		}
		close(urls)
	}()

	for r := range results {
		fmt.Printf("url = %v, status = %v, err = %v\n", r.url, r.status, r.err)
	}
}

func httpGet(urls <-chan string, results chan result) {
	for url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			results <- result{
				url: url,
				err: err,
			}
			return
		}
		results <- result{
			url:    url,
			status: resp.StatusCode,
		}
		_ = resp.Body.Close()
	}
	close(results)
}

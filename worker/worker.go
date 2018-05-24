package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"regexp"
	"sync"
)

func main() {
	results := make(chan string)
	errs := make(chan error)

	go producer(results, errs)

	for {
		select {
		case r, ok := <- results:
			if !ok {
				fmt.Printf("Channel is closed\n")
				return
			}
			fmt.Printf("result = %v\n", r)
		case err := <- errs:
			fmt.Printf("err = %v\n", err)
		}
	}
}

func producer(results chan string, errs chan error) {
	urls := []string{
		"https://oinume.hatenablog.com/entry/paid-mac-apps",
		"https://oinume.hatenablog.com/entry/generating-an-unpredictable-random-value-in-go",
		"https://oinume.hatenablog.com/entry/review-2018-04",
	}
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.DefaultClient.Get(url)
			//fmt.Printf("GET %v\n", url)
			if err != nil {
				errs <- err
				return
			}
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errs <- err
				return
			}
			defer resp.Body.Close()
			title := extractTitle(string(b))
			results <- title
		}(url)
	}
	wg.Wait()
	close(results)
	//close(errs)
}

func extractTitle(s string) string {
	title := regexp.MustCompile("<title>(.*)</title>")
	group := title.FindStringSubmatch(s)
	return group[1]
}

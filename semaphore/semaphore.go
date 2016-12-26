package main

import (
	"sync"
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"os"
	"flag"
)

var (
	concurrency = flag.Int("c", 1, "num of concurrency")
)

func main() {
	flag.Parse()
	semaphore := make(chan struct{}, *concurrency)
	wg := &sync.WaitGroup{}
	urls := []string{
		"http://oinume.hatenablog.com/",
		"http://oinume.hatenablog.com/entry/what-i-like-about-heroku",
		"http://oinume.hatenablog.com/entry/e2e-test-with-agouti-in-go",
		"http://oinume.hatenablog.com/entry/heroku-custom-clock-processes",
		"http://oinume.hatenablog.com/entry/mac-settings-on-sierra",
		"http://oinume.hatenablog.com/entry/mysqldump-option-where",
		"http://oinume.hatenablog.com/entry/introducing-lekcije",
		"http://oinume.hatenablog.com/entry/intellij-shortcuts-for-reading-source-code",
		"http://oinume.hatenablog.com/entry/introducing-dead-mans-snitch",
	}
	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			fetch(semaphore, u)
		}(u)
	}
	wg.Wait()
	os.Exit(0)
}

var r = regexp.MustCompile(`<title>(.*)</title>`)

func fetch(semaphore chan struct{}, url string) {
	semaphore<-struct{}{}
	defer func() {
		<-semaphore
	}()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}
	body := string(bytes)
	if group := r.FindStringSubmatch(body); len(group) > 0 {
		fmt.Printf("title = %v\n", group[1])
	}
}
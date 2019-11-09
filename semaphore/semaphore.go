package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

var (
	concurrency = flag.Int("c", 1, "num of concurrency")
)

func main() {
	flag.Parse()
	semaphore := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup
	urls := []string{
		"https://journal.lampetty.net/",
		"https://journal.lampetty.net/entry/what-i-like-about-heroku",
		"https://journal.lampetty.net/entry/e2e-test-with-agouti-in-go",
		"https://journal.lampetty.net/entry/heroku-custom-clock-processes",
		"https://journal.lampetty.net/entry/mac-settings-on-sierra",
		"https://journal.lampetty.net/entry/mysqldump-option-where",
		"https://journal.lampetty.net/entry/introducing-lekcije",
		"https://journal.lampetty.net/entry/intellij-shortcuts-for-reading-source-code",
		"https://journal.lampetty.net/entry/introducing-dead-mans-snitch",
	}
	for _, u := range urls {
		wg.Add(1)
		u := u
		go func() {
			defer wg.Done()
			fetch(semaphore, u)
		}()
	}
	wg.Wait()
}

var r = regexp.MustCompile(`<title>(.*)</title>`)

func fetch(semaphore chan struct{}, url string) {
	semaphore <- struct{}{}
	defer func() {
		<-semaphore
	}()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}
	body := string(bytes)
	if group := r.FindStringSubmatch(body); len(group) > 0 {
		fmt.Printf("%v\n", group[1])
	}
}

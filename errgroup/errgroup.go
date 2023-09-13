package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type transport struct {
	count int
	sync.Mutex
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.Lock()
	t.count++
	t.Unlock()

	code := http.StatusOK
	status := "200 OK"
	if t.count%3 == 0 {
		code = http.StatusInternalServerError
		status = fmt.Sprintf("500 error: count=%v", t.count)
	}
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    r,
		StatusCode: code,
		Status:     status,
		Body:       io.NopCloser(strings.NewReader("hoge")),
	}
	resp.Header.Set("Content-Type", "text/html; charset=UTF-8")

	return resp, nil
}

var client *http.Client

func main() {
	client = &http.Client{
		Timeout:   5 * time.Second,
		Transport: &transport{},
	}
	urls := []string{
		"https://eikaiwa.dmm.com/teacher/index/6210/",
		"https://eikaiwa.dmm.com/teacher/index/5654/",
		"https://eikaiwa.dmm.com/teacher/index/5616/",
		"https://eikaiwa.dmm.com/teacher/index/3923/",
		"https://eikaiwa.dmm.com/teacher/index/5412/",
		"https://eikaiwa.dmm.com/teacher/index/9848/",
		"https://eikaiwa.dmm.com/teacher/index/6122/",
		"https://eikaiwa.dmm.com/teacher/index/3370/",
		"https://eikaiwa.dmm.com/teacher/index/13786/",
		"https://eikaiwa.dmm.com/teacher/index/3133/",
	}
	eg := errgroup.Group{}
	for _, url := range urls {
		url := url
		eg.Go(func() error {
			return getRequest(url)
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

func getRequest(url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statusCode:%v, status:%v", resp.StatusCode, resp.Status)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("GET: %v\n", url)
	return nil
}

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Printf("err happened: %v\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	go func() {
		if err := startServer(8080); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)

	c := http.DefaultClient
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", "http://localhost:8080?duration=5s", nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := c.Do(req)
	if err != nil {
		//fmt.Printf("c.Do() failed: err=%v\n", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("body = %v\n", string(body))

	return nil
}

type server struct {
	innerServer *httptest.Server
}

func newServer() *server {
	is := httptest.NewServer()
	return &server{

	}
}

func startServer(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", sleep)
	return http.ListenAndServe(fmt.Sprintf("localhost:%d", port), mux)
}

func sleep(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("/sleep\n")
	duration := r.URL.Query().Get("duration")
	d, err := time.ParseDuration(duration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	go func() {
		s := newInnerServer()
		s.URL
	}

	time.Sleep(d)
	//w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "ok")
}

func newInnerServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(inner))
}

func inner(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("/inner\n")
	duration := r.URL.Query().Get("duration")
	d, err := time.ParseDuration(duration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	time.Sleep(d)
	//w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "ok")
}

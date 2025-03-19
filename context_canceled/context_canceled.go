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
	server := newServer()
	go func() {
		if err := server.start(8080); err != nil {
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
		fmt.Printf("c.Do() failed: err=%v\n", err)
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
	mux         *http.ServeMux
	innerServer *httptest.Server
}

func newServer() *server {
	is := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("/inner\n")
		duration := r.URL.Query().Get("duration")
		d, err := time.ParseDuration(duration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		time.Sleep(d)
		//w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "ok")
	}))
	s := &server{
		innerServer: is,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.outer)
	s.mux = mux

	return s
}

func (s *server) start(port int) error {
	return http.ListenAndServe(fmt.Sprintf("localhost:%d", port), s.mux)
}

func (s *server) outer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("/outer\n")
	duration := r.URL.Query().Get("duration")
	d, err := time.ParseDuration(duration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	c := http.DefaultClient
	req, err := http.NewRequest("GET", s.innerServer.URL+"/inner?duration=10s", nil)
	if err != nil {
		panic(err)
	}
	req.WithContext(r.Context())
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("body = %+v\n", string(body))

	time.Sleep(d)
	_, _ = fmt.Fprint(w, "ok:")
}

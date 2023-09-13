package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
)

func main() {
	server := server(routes(), 8080)
	server.Close()
}

func testServer() { //nolint:unused
	server := httptest.NewServer(nil)
	server.Close()
}

func server(handler http.Handler, port int) *httptest.Server {
	return &httptest.Server{
		Listener: newLocalListener(port),
		Config:   &http.Server{Handler: handler},
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()
	//mux.HandleFunc("/", func () {
	//
	//})
	return mux
}

func newLocalListener(port int) net.Listener {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		if l, err = net.Listen("tcp6", fmt.Sprintf("[::1]:%d", port)); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	return l
}

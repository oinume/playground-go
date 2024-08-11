package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.DefaultClient.Get("https://github.com/golang/go")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

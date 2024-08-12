package main

import (
	"fmt"
	"io/ioutil" //nolint:staticcheck
	"net/http"
)

func main() {
	resp, err := http.DefaultClient.Get("https://github.com/golang/go")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

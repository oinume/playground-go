package main

import (
	"io"
	"io/ioutil"
)

func before(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

func after(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

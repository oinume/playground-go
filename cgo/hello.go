package main

/*
#cgo LDFLAGS: -L. -lhello
#include <hello.h>
*/
import "C"

func main() {
	C.hello()
}

// $ go run hello.go

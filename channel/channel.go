package main

import "fmt"

func main() {
	//simple()
	buffered()
}

func simple() { //nolint:unused
	messages := make(chan string)
	go func() {
		messages <- "ping"
	}()
	msg := <-messages
	fmt.Println(msg)
}

func deadlock() { //nolint:unused
	messages := make(chan string)
	go func() {
		//	messages <- "ping"
		fmt.Printf("func()\n")
	}()
	msg := <-messages
	fmt.Println(msg)
}

func buffered() {
	messages := make(chan string, 2)

	messages <- "buffered"
	messages <- "channel"

	fmt.Println(<-messages)
	//fmt.Println(<-messages)
}

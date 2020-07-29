package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("f.Args(): %v\n", f.Args())
}

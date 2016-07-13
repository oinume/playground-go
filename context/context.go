package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func main() {
	root := context.Background()
	parent := context.WithValue(root, "parent1", "parent1 value")
	child := context.WithValue(parent, "child", "child value")
	fmt.Println(child.Value("parent1"))
	fmt.Println(child.Value("child"))
}
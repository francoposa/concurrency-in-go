package main

import (
	"fmt"
)

func main() {
	var x int
	go func() { x++ }()
	if x == 0 {
		fmt.Printf("x: %d\n", x)
	}
}

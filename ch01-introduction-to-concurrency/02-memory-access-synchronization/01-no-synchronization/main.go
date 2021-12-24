package main

import (
	"fmt"
)

func main() {
	var x int
	go func() { x++ }()
	go func() { x-- }()
	if x == 0 {
		fmt.Println("x is 0")
	} else {
		fmt.Printf("x is not 0, it is %d\n", x)
	}
}

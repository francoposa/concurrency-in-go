package main

import (
	"fmt"
)

func main() {
	// The following code forks but does not have a join point.
	// The helloWorld goroutines are created and scheduled,
	// but will likely not execute before the program exits

	// calling another func
	go helloWorld()

	// anonymous func
	go func() {
		fmt.Println("Hello, World")
	}()

	// assigning func to variable
	helloWorldInline := func() {
		fmt.Println("Hello, World")
	}
	go helloWorldInline()
}

func helloWorld() {
	fmt.Println("Hello, World")
}

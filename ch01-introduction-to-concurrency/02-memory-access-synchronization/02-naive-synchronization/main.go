package main

import (
	"fmt"
	"sync"
)

func main() {
	var x int
	var memoryAccess sync.Mutex

	go func() {
		memoryAccess.Lock()
		x++
		memoryAccess.Unlock()
	}()
	go func() {
		memoryAccess.Lock()
		x--
		memoryAccess.Unlock()
	}()

	memoryAccess.Lock()
	if x == 0 {
		fmt.Println("x is 0")
	} else {
		fmt.Printf("x is not 0, it is %d\n", x)
	}
	memoryAccess.Unlock()
}

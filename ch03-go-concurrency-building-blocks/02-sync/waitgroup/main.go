package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg0 sync.WaitGroup
	wg0.Add(1)
	go func() {
		defer wg0.Done()
		fmt.Println("goroutine 1 sleeping...")
		time.Sleep(1)
	}()
	wg0.Add(1)
	go func() {
		defer wg0.Done()
		fmt.Println("goroutine 2 sleeping...")
		time.Sleep(2)
	}()

	wg0.Wait()
	fmt.Println("goroutines complete")

	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %d!\n", id)
	}

	const numGreeters = 5
	var wg1 sync.WaitGroup
	// we try to add to the waitgroup as close as possible to the actual goroutines
	wg1.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg1, i+1)
	}
	wg1.Wait()
	fmt.Println("goroutines complete")
}

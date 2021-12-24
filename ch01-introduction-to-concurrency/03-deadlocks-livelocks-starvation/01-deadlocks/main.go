package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	type value struct {
		mu    sync.Mutex
		value int
	}

	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock() // attempt to claim a shared resource
		defer v1.mu.Unlock()

		// simulate doing work with the resource
		// this is not a perfect deadlock, technically this whole process could
		// complete before another process attempts to claim the same resources
		time.Sleep(2 * time.Second)

		v2.mu.Lock() // attempt to claim a second shared resource before releasing the first
		defer v2.mu.Unlock()

		fmt.Printf("sum: %d\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b) // lock a then attempt to lock b, without releasing a first
	go printSum(&b, &a) // lock b then attempt to lock a, without releasing b first
	wg.Wait()
}

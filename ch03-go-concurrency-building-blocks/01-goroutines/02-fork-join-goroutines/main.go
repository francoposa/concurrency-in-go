package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	//helloWorld := func() {
	//	fmt.Println("Hello, World")
	//}
	//go helloWorld()

	// The above code forks but does not have a join point.
	// The helloWorld goroutine is created and scheduled,
	// but will likely not execute before the program exits
	// We must create a join point
	var wg0 sync.WaitGroup
	helloWorld := func() {
		defer wg0.Done()
		fmt.Println("Hello, World")
	}
	wg0.Add(1)
	go helloWorld()
	wg0.Wait() // this is the join point

	// Goroutines are very lightweight
	// Here we measure the average memory consumed by starting a goroutine
	memConsumed := func() uint64 {
		runtime.GC()
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)
		return stats.Sys
	}

	var c <-chan interface{}
	var wg1 sync.WaitGroup
	noop := func() { wg1.Done(); <-c } // this will never exit
	const numGoroutines = 1e4
	wg1.Add(numGoroutines)
	before := memConsumed()
	for i := 0; i < numGoroutines; i++ {
		go noop()
	}
	wg1.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}

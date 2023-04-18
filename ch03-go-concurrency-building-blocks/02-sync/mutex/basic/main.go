package main

import (
	"fmt"
	"sync"
)

// main demonstrates basic mutex usage
// Two functions: one to increment a counter, another to decrement.
// Each function is scheduled to be called in n goroutines.
// Each goroutine competes to obtain an exclusive lock on the counter variable.
// Each function locks the counter before incrementing/decrementing the counter by 1,
// then prints the updated value to stdout before unlocking the counter.
//
// Locking is necessary because each increment/decrement function is four atomic operations:
//     1. Retrieve the value of the counter variable from its memory location
//     2. Increment/Decrement the value of the counter by 1
//     3. Store the new value of the counter back in its memory location
//     4. Retrieve the new value of the counter from its memory location and print it to stdout
//
// With locking, the counter can only ever increment or decrement by one from the previous value.
// Further, the operation printed to stdout will always make sense with the change in value
//
// Without locking, we could observe unexpected counter behavior via stdout.
// Example concurrency scenario without locking:
//     0. Counter var is `0`; main func prints `initial state: 0` to stdout.
//     1. Increment func reads the counter value of `0` from memory
//     2. Decrement func reads the counter value of `0` from memory
//     3. Increment func adds `1` to the `0` value
//     4. Decrement func subtracts `1` from the `0` value
//     5. Increment func writes its new value of `1` to memory
//     6. Decrement func writes its new value of `-1` to memory
//     7. Increment func reads the new value of `-1` from memory and prints `incrementing: -1`
//     8. Decrement func reads the new value of `-1` from memory and prints `decrementing: -1`
// stdout:
//     initial state: 0
//     incrementing: -1
//     decrementing: -1
// This is clearly not correct.
// Incrementing the counter from 0 should get us 1, not -1.
// Decrementing from 0 correctly got us -1, but from stdout it appears that the counter value
// was at -1 before the decrement operation, so it appears that the decrement operation failed.
//
// Comment out all lock usage to observe the naive behavior.
func main() {
	var count int
	fmt.Printf("initial state: %d\n", count)

	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("decrementing: %d\n", count)
	}

	const numOps = 5
	var arithmetic sync.WaitGroup
	for i := 0; i < numOps; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}
	for i := 0; i < numOps; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("complete")
}

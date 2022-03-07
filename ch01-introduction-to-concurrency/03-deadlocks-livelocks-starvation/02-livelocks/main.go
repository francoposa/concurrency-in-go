package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// TODO document this code after I understand it
	// I do not really understand this code at all
	// I get what type of situation we are trying to model
	// but not why this code successfully models the situation

	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			// wake up all goroutines waiting on `cadence`
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDirection := func(directionName string, direction *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, "%s\n", directionName)
		atomic.AddInt32(direction, 1)
		takeStep()
		if atomic.LoadInt32(direction) == 1 {
			fmt.Fprintln(out, "success")
			return true
		}

		takeStep() // not sure why we takeStep again here
		// decrement to indicate we didn't actually succeed in stepping?
		atomic.AddInt32(direction, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool {
		return tryDirection("left", &left, out)
	}
	tryRight := func(out *bytes.Buffer) bool {
		return tryDirection("right", &right, out)
	}

	walk := func(walking *sync.WaitGroup, walkerName string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		fmt.Fprintf(&out, "%s is trying to walk: ", walkerName)
		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}

		fmt.Fprintf(&out, "%s has given up", walkerName)
	}

	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "A")
	go walk(&peopleInHallway, "B")
	peopleInHallway.Wait()
}

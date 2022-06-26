# Chapter 3: Go's Concurrency Building Blocks

## Goroutines
**Goroutine**: a  functions closure, or method running concurrently in Go.

Each program has at least one goroutine, the *main goroutine*.

### Starting Goroutines
Goroutines can be started in a few different ways:

```go
func main() {
	// The following code forks but does not have a join point.
	// The helloWorld goroutines are created and scheduled,
	// but will likely not execute before the program exits

	// calling another func
	go helloWorld()

	// anonymous func
	// must be called immediately in order to use the `go` keyword
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
```
### What Are Goroutines?
Goroutines are not OS threads or green threads (threads managed by a language's runtime).
Goroutines are **coroutines** - concurrent subroutines which are *non-pre-emptive*.
In [**non-preemptive multitasking**](https://en.wikipedia.org/wiki/Cooperative_multitasking),
the operating system never initiates a context switch from a running process to another process.
Instead, processes voluntarily yield control periodically or when idle or logically blocked.
This called **cooperative multitasking** because all programs must cooperate for the scheduling scheme to work.

What makes goroutines unique is their deep integration with the Go runtime.
Goroutines do not define their own suspension or reentry points.
Go's runtime observes the runtime behavior of the goroutines and automatically suspends and resumes them as they become blocked and unblocked.
This makes the goroutines preemptable, but only at points where they have become blocked.
Thus, goroutines can be considered a special class of coroutines.

Coroutines and therefore goroutines are concurrent constructs, but concurrency is not a property of a coroutine.
Something must host multiple coroutines simultaneously and give each an opportunity to execute - this is a cooperative process scheduler.

Go employs an *M:N* scheduler, mapping *M* green threads onto *N* OS threads.
Goroutines are then scheduled onto the green threads.
The scheduler handles the distribution of the goroutines across the available green threads and ensures that when some goroutines becomes blocked, others can be run.

#### The Fork-Join Concurrency Model

Go follows the **fork-join concurrency model**.
*Fork* refers to the fact that at any point in the program, the *parent* branch of execution can split off a *child* branch to be run concurrently with the parent.
*Join* refers to the fact that at some point in the future, these concurrent branches of execution will join back together, called a *join point*.

The fork-join model is a logical model describing how concurrency is performed at a conceptual level.
It is agnostic to any implementation details regarding scheduling of concurrency processes, memory management, etc.

To see the fork-join model at work in Go, we first return to the previous example:

```go
helloWorld := func() {
    fmt.Println("Hello, World")
}
go helloWorld()
```

In this example, `helloWorld` is scheduled to be run in its own goroutine while the rest of the program continues.
The goroutine will execute `helloWorld` at some undetermined point in the future.
There is no join point, so the main goroutine will not wait for the result of the other goroutine.

In fact the main goroutine, and therefore the program, may exit before `sayHello` ever gets a chance to run.
Because in this simple example there is no further code to execute after scheduling `helloWorld`, it is almost certain to exit before `hellowWorld` is run.
We will not see anything printed to stdout.

We can add a sleep condition to the main goroutine, but this is a race condition.
We are hoping to give the other goroutine a chance to complete before the program exits, but we cannot guarantee it.

We must introduce a join point in order to guarantee the correctness of our program and remove the race condition.
There are several ways to do this in Go, but our most basic way is to use the `sync` package's `WaitGroup`.

Without worrying (yet) about exactly how `sync.WaitGroup` operates, here is the correct version of the example:

```go
var wg0 sync.WaitGroup
helloWorld := func() {
    defer wg0.Done()
    fmt.Println("Hello, World")
}
wg0.Add(1)
go helloWorld()
wg0.Wait() // this is the join point
```

### The `sync` Package

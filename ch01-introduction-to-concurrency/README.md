# Chapter 1: An  Introduction to Concurrency

Broad, casual definition of concurrency:
multiple simultaneous processes making progress at approximately the same time.

## Moore's Law and the Mess We're In

**Moore's Law** states that the number of transistors in a compute chip will double roughly every two years.
This held approximately true until around 2012.

Computing companies foresaw the limits of Moore's law and started investing in alternate ways to increase computing power, namely multicore processors.

Multicore - and therefore parallel - computing is bound by Amdahl's Law.
**Ahmdahl's Law** states that the gains of parallel computing are limited by how much of the program must be executed sequentially, and is therefore non-parallelizable.

Problems that are easily divided into parallel tasks are called **embarassingly parallel**.
Embarrassingly parallel problems are perfect candidates for **horizontal scaling**, that is distributing sub-problems across multiple CPU cores or machines.

Cloud computing has enabled horizontal scaling due to the ease and affordability of spinning up many compute resources on demand.

Horizontal scaling and distributed computing have introduced an entirely new class of challenges to go along with the new capabilities.
Provisioning compute resources, communicating between distributed machines, and aggregating and reporting program results are far from solved problems.

Modeling concurrent code is essential to proper understanding of distributed systems.
Due to the ubiquity of distributed systems today, this is something almost every engineer must eventually face.

## Why is Concurrency Hard?

### Race Conditions

**Race Condition**: when two or more operations must execute in the correct order, but the program has not been written so that this order is guaranteed.

This often shows up as a **data race**, where one or more concurrent operations attempt to read a variable while other concurrent operations are writing to the variable.

Data races can be introduced when engineers are incorrectly thinking about the problem sequentially, when in reality there are no guarantees that one operation will run before another.

When modeling concurrent code, it is helpful to meticulously list out the possible concurrency scenarios.

Take this basic example:

```go
func main() {
    var x int // initializes to 0
    go func() { x++ }() // gets scheduled to run concurrently with the remaining code
    if x == 0 {
        fmt.Printf("x: %d\n", x)
    }
}
```

Here we call a goroutine (`go func() { // do something }()`), which schedules it to be run outside of the main goroutine.
We do not apply any guards or inter-goroutine communication to coordinate the execution of the two goroutines.
Any step of either of the two goroutines could be executing, or not executing at any time, in any order.

Right off the bat, there are three basic concurrency scenarios we can see here.

**Concurrency Scenario 1:** Nothing is printed.
1. The second goroutine increments the variable to 1.
2. The main goroutine checks the value of the variable, which is not 0, then exits.

**Concurrency Scenario 2:** `x: 0` is printed.
1. The main goroutine checks the value of the variable, which is 0.
2. The main goroutine prints the value of the variable, which is still 0, then exits.

In this scenario, the second goroutine could increment the variable after the print, but before the main goroutine exits.
Alternatively, the main goroutine could exit before the second goroutine gets a chance to execute.
Either way, the observed behavior is the same.

**Concurrency Scenario 3:** `x: 1` is printed.
1. The main goroutine checks the value of the variable, which is 0.
2. The second goroutine increments the variable to 1.
3. The main goroutine prints the value of the variable, which is now 1, then exits.

#### Danger!

If we run this code ourselves, we will likely see `x: 0` printed every time, or almost every time.
This is because the code we are running in the main goroutine is so simple that it is very likely that it completes checking and printing the variable before the second goroutine increments the value.

This is an easy trap to fall into!

Using observed behavior can lead us to either:

1. completely miss the fact that a race condition exists, or
2. assume that one operation will complete so quickly or slowly that a race condition with another operation is highly unlikely

Number 2 is especially tempting.
However, in the real world any operation can get suspended for an arbitrary amount of time, due to scheduling of CPU resources, long-running or hung network calls, etc.

When modeling concurrency scenarios, it is best to discard any assumptions about how long an operation may be suspended or take to complete.
It helps to imagine arbitrarily large delays (multiple minutes, hours, or days!) for any process.

It may be tempting to try to solve these problems by adding sleeps into our code, like below:

```go
func main() {
    var x int
    go func() { x++ }()
    time.Sleep(1 * time.Second) // this will not solve your problems
    if x == 0 {
        fmt.Printf("x: %d\n", x)
    }
}
```

This is bad!
We have not solved the race condition, we have just made it less likely.
Further, we are almost certainly wasting execution time by sleeping for too long.

With the correct approach and tools, we may wait for a millisecond or less.

### Atomicity

**Atomic operations** are indivisible, or uninterruptible *within the context that they are operating*.

Something may be atomic in one context, but not another.
An operation may be atomic in the context of our process but not our operating system.
Another operation may be atomic in the context of our operating system but not the machine it is executing on.

By "indivisible" or "uninterruptible", we mean that an atomic operation will execute in its entirety without any other operation in that context executing simultaneously.

**Example:** Incrementing an Integer Variable

In Go, this just looks like `i++`.
Hidden under this short statement are three different sub-operations:

1. Retrieve the value of `i` from its memory location
2. Increment the value of `i`
3. Store the new value of `i` back in its memory location

In the context of a single program with no concurrent processes (threads, goroutines, or otherwise), this is an atomic operation.
If we have multiple concurrent processes, but `i` is not exposed outside a single process, then the operation is still atomic.
However, if multiple concurrent processes can access `i`, as in our race condition example, the operation is no longer atomic.

Why do we care?
If an operation is atomic, it is implicity safe within concurrent contexts.

### Memory Access Synchronization

We can expand upon our previous data race example.
We now have goroutine calls to increment *and* decrement, and we always print the value.

```go
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
```

Again, as it is written, the output is nondeterministic.
We do not know:

* what value will be printed.
* if the value in memory when the if-statement is evaluated will match the value printed
  * the value could be 1 when the if-statement is evaluated, then a goroutine could decrement the value before the print executes
    * in this case, the program would print `x is not 0, it is 0`.
* if the value printed will be the same as the value in memory when the program completes
  * a goroutine could change the value after the print executes, but before the program exits


There are multiple sections of this program that need *exclusive access* to a shared resource, in this case a variable in memory.

A section of a program that need exclusive access to a shared resource is called a **critical section**.

We have four critical sections:

1. the goroutine which increments the value
2. the goroutine which decrements the value
3. the if-statement which evaluates the value and branches accordingly
4. whichever print statement executes, which evaluates the value to output

We can guard these critical sections with a **mutex**, short for "mutual exclusion".

This code is not good Go code or good concurrency control!
This is just a basic example to demonstrate memory access synchronization.

```go
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
```

This memory access synchronization has solved two of our three issues:

* the value in memory when the if-statement is evaluated will match the value printed
* the value printed will be the same as the value in memory when the program completes

However, we still do not know *what* the value will be when we enter the critical section to print the value.
This means that not only do we not know what value will be printed, we also do not know which branch of the if-statement will execute.
The order of operations is still nondeterministic.

Further, memory access synchronization via locking has performance ramifications.
While the data is locked, other processes that need access to the data cannot proceed.
If programs are locking very often, locking longer than necessary, or many processes are competing for the lock, performance can deteriorate rapidly.

# Go Concurrency

1. Go concurrency model is the `goroutine`, which is the minimum executable entity in Go. Everything in Go is executed as a goroutine, either implicitly or explicitly.
2. Each goroutine is executed on a single OS thread according to the instructions of the Go scheduler, which is responsible for the execution of goroutines. The OS scheduler does not dictate how many threads the Go runtime is going to create because the Go runtime will spawn enough threads to ensure that `GOMAXPROCS` threads are available to run Go code.
3. However, goroutines cannot directly communicate with each other. Data sharing in Go is implemented using either channels or shared memory.
4. When we combine multiple channels and goroutines, we can create data flows, which in Go terminology are also called `pipelines`. 

`Processes, Threads, and Goroutines`

A goroutine is the minimum Go entity that can be executed concurrently. The use of the word minimum is very important here, `as goroutines are not autonomous entities like UNIX processes—goroutines live in OS threads that live in OS processes`. The good thing is that goroutines are lighter than threads, which, in turn, are lighter than processes—running thousands or hundreds of thousands of goroutines on a single machine is not a problem. In practice, this means that a process can have multiple threads as well as lots of goroutines, whereas a goroutine needs the environment of a process to exist. `So, to create a goroutine, we need to have a process with at least one thread`. The OS takes care of the process and thread scheduling, while Go creates the necessary threads, and the developer creates the set of goroutines needed for the task at hand.

`The Go Scheduler`

**m : n scheduling** - The OS kernel scheduler is responsible for the execution of the threads of a program. Similarly, the Go runtime has its own scheduler, which is responsible for the execution of the goroutines using a technique known as m:n scheduling, where m goroutines are executed using n OS threads using multiplexing.

**Fork-join concurrency** - Go uses the fork-join concurrency model.States that a child branch can be created at any point of a program. Analogously, the join part of the Go concurrency model is where the child branch ends and joins with its parent. 

Keep in mind that both sync.Wait() statements, as well as channels that collect the results of goroutines, are join points, whereas each new goroutine creates a child branch.

**Scheduling strategy** 

1. Fair Scheduling - (Old)
2. Work Stealing Algorithm - (New)
    1. Go now uses work-stealing:
    2. Each logical processor (P) has a local run queue.
    3. If a processor is idle, it steals work from others.
    4. Go steals continuations rather than tasks for efficiency.
    5. This approach improves load balancing and reduces idle time.

Core Entities: G, M, P
The scheduler manages three key resources:

G (goroutine): Represents a lightweight thread of execution.
M (machine): An OS thread that runs Go code or blocks in syscalls.
P (processor): Holds resources like scheduler state and memory allocator.


The Go scheduler works using three main kinds of entities: OS threads (M), which are related to the OS in use; goroutines (G); and logical processors (P). The number of processors that can be used by a Go program is specified by the value of the GOMAXPROCS environment variable—at any given time, there are at most GOMAXPROCS processors. Now, let us return to the m:n scheduling algorithm used in Go. Strictly speaking, at any time, we have m goroutines that are executed and, therefore, scheduled to run on n OS threads using, at most, GOMAXPROCS number of logical processors. 

Go multiplexes lightweight goroutines (G) onto operating system threads (M) managed by logical processors (P). Each P holds a Local Run Queue (LRQ) of goroutines waiting to execute.

Here’s the Go Scheduler Architecture showing:

Global Run Queue → distributes goroutines to local queues.
Local Run Queues → attached to each logical processor (P).
P (Processor) → manages execution resources.
M (Machine) → OS thread executing Go code.
Goroutines (G) are scheduled across Ps and Ms using work-stealing.

The Go scheduler is a masterpiece of simplicity and efficiency. By combining m:n scheduling, work-stealing, and lightweight goroutines, Go delivers high-performance concurrency without the complexity of traditional thread management.

`The GOMAXPROCS Environment Variable`

GOMAXPROCS is a Go runtime setting that defines the maximum number of operating system threads (M) that can execute Go code simultaneously by controlling the number of logical processors (P) in the Go scheduler.

The GMP Model - P (Processor): Logical execution contexts; the number of Ps is equal to GOMAXPROCS

Default Value: Traditionally defaults to runtime.NumCPU() (the total number of logical CPU cores on the host machine). There is also the runtime.GOMAXPROCS() function, which allows us to set and get the value of GOMAXPROCS programmatically.

Container & Kubernetes Awareness

Prior to Go 1.25, the Go runtime ignored Linux cgroup CPU quotas, often matching the host machine's total core count (e.g., 32 cores) rather than a container's tight limit (e.g., 2 cores), which triggered heavy CPU throttling and garbage collection delays. 

From Go 1.25 onward, GOMAXPROCS is container-aware by default and automatically adapts to cgroup CPU limits.

Oversubscription Risks - Setting GOMAXPROCS higher than the number of available physical/logical cores forces excessive context switching and cache thrashing rather than increasing throughput

`Concurrency and Parallelism`

In Go, the distinction between concurrency and parallelism is foundational. As defined by Go co-creator Rob Pike, `"Concurrency is about dealing with lots of things at once. Parallelism is about doing lots of things at once"`.

While they are closely related, they represent completely different aspects of software: concurrency is a design property of your code, whereas parallelism is a runtime execution property of the hardware.

`Goroutines`

1. We can define, create, and execute a new goroutine using the go keyword followed by a function name or an anonymous function.
2. As a Go program(main finction) does not wait for its goroutines to end before exiting, we need to delay it manually, Ex: the time.Sleep() call.

```go
func main() {
    //anonymous function as a goroutine
	go func(x int) {
		fmt.Printf("%d ", x)
	}(10) //The (10) at the end is how we pass a parameter to an anonymous function.
    //function as a goroutine
    go somefun()
}
```

`Creating Multiple Goroutines`
There is nothing prohibiting us from using a for loop to create multiple goroutines, especially when we want to create lots of them.
```go
for i := 0; i < count; i++ {
	go func(x int) {
		fmt.Printf("%d ", x)
	}(i)
}
```

`Waiting for Our Goroutines to Finish`

**The sync.WaitGroup variable** - The synchronization process begins by defining a `sync.WaitGroup` variable and using the `Add(), Done() and Wait()` methods. Each call to `sync.Add()` increases a counter in the state1 field, which is an array with three uint32 elements. Notice that it is really important to call sync.Add() before the go statement in order to prevent any race conditions. When each goroutine finishes its job, the `sync.Done()` function should be executed in order to decrease the same counter by one. Behind the scenes, `sync.Done()` runs an `Add(-1)` call. The `Wait()` method waits until that counter becomes 0 in order to return. The return of Wait() inside the main() function means that main() is going to return and the program ends.
```
var waitGroup sync.WaitGroup
fmt.Printf("%#v\n", waitGroup)
for i := 0; i < count; i++ {
	waitGroup.Add(1)
	go func(x int) {
		defer waitGroup.Done()
		fmt.Printf("%d ", x)
	}(i)
}

fmt.Printf("%#v\n", waitGroup)
waitGroup.Wait()
fmt.Println("\nExiting...")
```

`What If the Number of Add() and Done() Calls Differ?`
1. `panic: sync: negative WaitGroup counter`, which is caused by more calls to Done() than calls to Add().
2. `fatal error: all goroutines are asleep - deadlock!`, which means that the program should wait indefinitely for a goroutine to finish, that is, for a Done() call that is never going to happen.
3. These two simulations demonstrate the importance, in the design of concurrent programs, of ensuring that the number of sync.Add() calls equals the number of sync.Done() calls.

`Channels`
1. A channel is a communication mechanism that, among other things, allows goroutines to exchange data.
2. We should declare a new channel using make() and the chan keyword (make(chan int)), and we can close a channel using the close() function - ` make(chan int, 1)`
3. A `pipeline` is a virtual method for connecting goroutines and channels so that the output of one goroutine becomes the input of another goroutine, using channels to transfer our data. One of the benefits that we get from using pipelines is that there is a constant data flow in our program, as no goroutine or channel has to wait for everything to be completed in order to start their execution. 

**Writing to and reading from a channel**
1. Writing = `ch <- val`. The arrow shows the direction of the value, and both var and ch are of the same data type.
2. Reading from a channel named c by executing `<-c`. We can save that value into a variable using `aVar := <-c`.

**Buffered Channels vs UnBuffered Channels**

1. Buffered channel means that as soon as we fill that buffer, we can close the channel, and the goroutine is going to continue its execution and return. `ch := make(chan int, 1)`
2. A channel that is unbuffered has a different behavior: when we try to send a value to that channel, it blocks forever because it is waiting for someone to fetch that value. `ch := make(chan int)`
3. When you use a `range loop` on an unbuffered channel in Go, the loop will sequentially receive values sent to the channel. Because the channel is unbuffered, each iteration blocks until another goroutine sends a value. The loop will only terminate when the channel is explicitly closed by the sender. If the channel is never closed and no more values are sent, the loop will block forever, causing a deadlock.
```go
func main() {
	// 1. Create an unbuffered channel
	ch := make(chan int)

	// 2. Start a sender goroutine
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i // Blocks until main goroutine receives it
		}
		close(ch) // Crucial: Closes the channel to terminate the range loop
	}()

	// 3. Range over the unbuffered channel in the main goroutine
	for val := range ch {
		fmt.Println(val) // Prints 1, 2, 3 sequentially
	}
	
	fmt.Println("Channel closed, loop finished.")
}
```

**Receiving from a Closed Channel** - Reading from a closed channel returns the zero value of its data type. However, if we try to write to a closed channel, our program is going to crash in a bad way (panic)

`Channels as Function Parameters`
1. When using a channel as a function parameter, we can specify its direction; that is, whether it is going to be used for sending or receiving data.
2.  We will not be able to send data accidentally to a channel from which we should only receive data or receive data from a channel to which we should only be sending data.
3. `func writeToChannel(out <-chan int, in chan<- int) {}`

`Race Conditions`
**A data race condition** - A data race condition is a situation where two or more running elements, such as threads and goroutines, try to take control of or modify a shared resource or shared variable of a program. `The Go race detector` - go run -race

`The select Keyword`
The select statement in Go is used to let a goroutine wait on multiple communication operations (channel sends or receives). It blocks until one of its cases can run, making it the core tool for managing concurrency and handling asynchronous events. Think of it like a switch statement, but instead of evaluating values, it evaluates channel operations.
A select block can have multiple cases and an optional default case, which mimics the switch statement. It is good for select blocks to have a timeout option just in case. Last, a select without any cases (select{}) waits forever.
In practice, this means that select allows a goroutine to wait on multiple communication operations. So, select gives us the power to listen to multiple channels using a single select block. As a consequence, we can have nonblocking operations on channels, provided that we have implemented our select blocks appropriately.


Core Mechanics of select

1. Blocking by Default: If no channels are ready, the select statement blocks and waits until at least one case can proceed.
2. Random Selection: If multiple cases are ready at the same time, Go picks one at uniform pseudo-random selection. This prevents "starvation" of lower cases.
3. Single Execution: A select block executes exactly one case and then exits. To continuously monitor channels, you must wrap it in a for loop.


`Timing Out a Goroutine`

```go
select {
case res := <-c2:
    fmt.Println(res)
case <-time.After(4 * time.Second):
    fmt.Println("timeout c2")
}
```

`Buffered Channels`
1. These channels allow us to put jobs in a queue quickly in order to be able to deal with more requests and process requests later on. Moreover, we can use buffered channels as semaphores in order to limit the throughput of our application.
2. The capacity of the buffer of the channel restricts the number of simultaneous requests that it can keep.
3. `numbers := make(chan int, 5)` 

`Worker Pools`
1. Worker pools in Go are implemented with the help of buffered channels because they allow us to limit the number of goroutines running at the same time.

`Shared Memory and Shared Variables`

1. A `mutex` variable, which is an abbreviation of `mutual exclusion` variable, is mainly used for thread synchronization and for protecting shared data when multiple writes can occur at the same time.
2. A mutex works like a buffered channel with a capacity of one, which allows, at most, one goroutine to access a shared variable at any given time.
3. Go offers the `sync.Mutex` and `sync.RWMutex` data types.
4. All of the interesting work is done by the `sync.Lock()` and `sync.Unlock()` functions, which can lock and unlock a sync.Mutex variable, respectively. Locking a mutex means that nobody else can lock it until it has been released using the sync.Unlock() function.
```go
var m sync.Mutex
var v1 int

func change(i int) {
	m.Lock()
	time.Sleep(time.Second)
	v1 = v1 + 1
	if v1 == 10 {
		v1 = 0
		fmt.Print("* ")
	}
	m.Unlock()
}
```

**Forgetting to unlock a sync.Mutex mutex** - Forgetting to unlock a sync.Mutex mutex creates a panic situation even in the simplest kind of a program. The same applies to the sync.RWMutex mutex - `fatal error: all goroutines are asleep - deadlock!`

`The sync.RWMutex data type`

sync.RWMutex stands for `Read-Write Mutual Exclusion Lock`. It is a specialized version of a standard sync.Mutex designed to optimize performance in scenarios where `data is read frequently but modified rarely`.

Rule:
1. Many Readers: An arbitrary number of readers can hold the lock simultaneously, as long as no writer is waiting or writing.
2. One Writer: Only one writer can hold the lock at a time. While a writer is active, all readers and other writers are completely blocked.

The Four Essential Methods: RLock(), RUnlock(), Lock(), Unlock()

`Do Not Copy the Mutex:` A sync.RWMutex contains internal state counters. If you copy the struct containing the mutex, you copy its internal state, which will corrupt your locking mechanics. `Always pass pointers to structs containing mutexes`.

`Sharing Memory Using Goroutines`

**Share data using a dedicated goroutine** - Although shared memory is the traditional way that threads communicate with each other, Go comes with built-in synchronization features that allow a single goroutine to own a shared piece of data. This means that other goroutines must send messages to this single goroutine that owns the shared data, which prevents the corruption of the data. Such a goroutine is called a `monitor goroutine`.


`Closured Variables and the go Statement`

1. `closured variables`, which are variables inside closures, and the go statement. Notice that closured variables in goroutines are evaluated when the goroutine actually runs and when the go statement is executed in order to create a new goroutine. This means that closured variables are going to be replaced by their values when the Go scheduler decides to execute the relevant code.

**The Loop Problem**

When an anonymous function inside a loop uses a loop variable directly, it captures that variable by reference (sharing its memory address).
Because the loop moves quickly and the functions (like goroutines) execute asynchronously, the loop often finishes before the functions even run. As a result, every single function reads the exact same final value of the loop variable, instead of the unique value from its specific iteration.

**The Fix**
To fix this, pass the loop variable as an argument into a parameter of the anonymous function.
Because Go is strictly pass-by-value, passing the variable as an argument forces Go to immediately copy the current value at that exact millisecond. This gives each function its own isolated, unchanging copy of the data, safely protecting it from the loop's progression.

```go
for i := 0; i <= 20; i++ {
    go func(x int) {
        fmt.Print(x, " ")
    }(i) //pass the loop variable as an argument
}
```


# Introduction to Go


`Intro and Facts`

1. Go will not allow the program to get executed if there are any unused imports or variables. 
2. Go has support for concurrency out of the box through goroutines and channels.
3. Goroutines are really lightweight. We can easily run thousands of goroutines on any modern machine without any performance issues.
4. Unlike C, Go supports functional programming.
5. Go has reserved only 25 keywords, which makes it much easier to remember the language.
6. Go code is also easy to read and offers support for Unicode out of the box, which means that we can make changes to existing Go code easily

`Important features and advantages of Go`

1. **Goroutines and channels** - Go manages OS threads for us and has a powerful runtime that allows us to spawn lightweight units of work (goroutines) that communicate with each other using channels.
2. **Pointers and interfaces** - Although Go is not an object-oriented programming language, Go interfaces are very versatile and allow us to mimic some of the capabilities of object-oriented languages, such as polymorphism, encapsulation, and composition.
3. Go generics and garbage collection - Additionally, the latest Go versions offer support for generics, which simplifies our code when working with multiple data types. Last but not least, Go comes with support for garbage collection, which means that no manual memory management is needed.

`Accessing Go Documentation`

The `go doc` and `godoc` utilities. Because godoc is not installed by default, we might need to install it by running go get golang.org/x/tools/cmd/godoc.

```
go doc fmt.Printf
go doc fmt

godoc -http=:8001
```

`Functions`

1. Everything that begins with a lowercase letter is considered private and is only accessible in the current package. The only exception to this rule is package names, which can begin with either lowercase or uppercase letters. Having said that, we are not aware of a Go package that begins with an uppercase letter!

`Compiling Go code` - Go Build
`Using Go like a scripting language` - Go Run

`Important formatting and coding rules`

1. Go code is delivered in packages, and we are free to use the functionality found in existing packages. However, if we are going to import a package, we should use some of this functionality—there are some exceptions to this rule that mainly have to do with initializing connections, but they are not important for now.
2. We either use a variable or we do not declare it at all. This rule helps us avoid errors such as misspelling an existing variable or function name.
3. There is only one way to format curly braces in Go.
4. Coding blocks in Go are embedded in curly braces, even if they contain just a single statement or no statements at all.
5. Go functions can return multiple values.
6. We cannot automatically convert between different data types, even if they are of the same kind. As an example, we cannot implicitly convert an integer to a floating point.

`Go Variables` 
1. `var` keyword - We can declare a new variable using the var keyword followed by the variable name, followed by the desired data type. If we want, we can follow that declaration with = and an initial value for our variable
2. Short assignment statement - `:=`
3. Although we can declare local variables using either var or :=, only const (when the value of a variable is not going to change) and var work for global variables, which are variables that are defined outside of a function and are not embedded in curly braces.

`Printing Variables`
1. For printing data on the screen, Go uses the functionality of the fmt package.
2. fmt.Println() / fmt.Printf()

`Controlling program flow`

1. `if/else` & `if` has support for `else` and `else if`
2. `switch` - The switch statement has two different forms. In the first form, the switch statement has an expression that is being evaluated, whereas in the second form, the switch statement has no expression to evaluate. In that case, expressions are evaluated in each case statement, which increases the flexibility of switch. The main benefit we get from switch is that when used properly, it simplifies complex and hard-to-read if-else blocks

`Loops`

*Iterating with `for` loops and the `range` keyword*
1. Go supports for loops as well as the range keyword for iterating over all the elements of arrays, slices, and maps. An example of Go’s simplicity is the fact that Go only provides support for the for keyword instead of including direct support for while loops. However, depending on how we write a for loop, it can function as a while loop or an infinite loop. Moreover, for loops can implement the functionality of JavaScript’s forEach function when combined with the range keyword.
2. for i := 0; i < 10; i++ {} - Traditional for loop
3. for ok := true; ok; ok = (i != 10) {} - Idiomatic Go - For loop used as do-while loop.
4. for {} - For loop used as while loop - The `break` keyword exits the loop early and acts as the loop exit condition.

`Getting User Input`

1. `Reading from standard input` - The fmt.Scanln() function can help us read user input while the program is already running and store it to a string variable, which is passed as a pointer to fmt.Scanln().
2. `command-line arguments` - By default, command-line arguments in Go are stored in the os.Args slice. Go also offers the flag package for parsing command-line arguments, but there are better and more powerful alternatives.

`Error Variables` - Using error variables to differentiate between input types - Technique that uses error variables to differentiate between various kinds of user input.

`Concurrency`

1.  The Go concurrency model is implemented using goroutines and channels. A goroutine is the smallest executable Go entity. In order to create a new goroutine, we have to use the go keyword followed by a predefined function or an anonymous function—both methods are equivalent as far as Go is concerned.
2. A channel in Go is a mechanism that, among other things, allows goroutines to communicate and exchange data.
3. As main() runs as a goroutine as well, we do not want main() to finish before the other goroutines of the program because when main() exits, the entire program, along with any goroutines that have not finished yet will terminate.
4. The good thing is that there are various techniques for the main() function to wait for goroutines to exchange data through channels or, less frequently in Go, using shared memory. Ex: Synchronize goroutines using time.Sleep() calls, wait groups and more.

`logging`
1. The UNIX logging service has support for two properties named logging level and logging facility.
2. The log package of the standard Go library does not support working with logging levels. 
3. The log package sends log messages to standard error. Part of the log package is the log/syslog package, which allows us to send log messages to the syslog server of our machine. Although by default, log writes to standard error, and the use of log.SetOutput() modifies that behavior.
4. The list of functions for sending logging data includes log.Printf(), log.Print(), log.Println(), log.Fatalf(), log.Fatalln(), log.Panic(), log.Panicln() and log.Panicf().
5. Keep in mind that log.Panic() is equivalent to a call to log.Print() followed by a call to panic(). panic() is a built-in function that stops the execution of the current function and begins panicking. After that, it returns to the caller function. On the other hand, log.Fatal() calls log.Print() and then os.Exit(1), which is an immediate way of terminating the current program.
6. Writing to a custom log file -  `log.New(file, "log name ", log.LstdFlags)`
7. Printing line numbers in log entries - he desired functionality is implemented with the use of log.Lshortfile in the parameters of log.New() or SetFlags(). The log.Lshortfile flag adds the file name as well as the line number of the Go statement that printed the log entry in the log entry itself. If we use log.Llongfile instead of log.Lshortfile, then we get the full path of the Go source file.

`Go Generics`

1. The main idea behind generics in Go, as well as any other programming language that supports generics, is not having to write special code for supporting multiple data types when performing the same task. Currently, Go supports multiple data types in functions, such as fmt.Println() using the empty interface and reflection.
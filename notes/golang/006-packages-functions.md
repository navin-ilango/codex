# Packages and Functions

Everything in Go is delivered in the form of packages. A Go package is a Go source file that begins with the package keyword, followed by the name of the package.

`Functions`

All function definitions begin with the `func` keyword, followed by the function’s signature and its implementation and that functions accept zero or more arguments and return zero or more values back.

1. The single most popular Go function is main(), which is used in every executable Go program—the main() function accepts no parameters and returns nothing, but it is the starting point of every Go program. Additionally, when the main() function ends, the entire program ends as well.

`Anonymous Functions`
1. Anonymous functions can be defined inline without the need for a name, and they are usually used for implementing things that require a small amount of code.
2. In Go, a function can return an anonymous function or take an anonymous function as one of its arguments. Additionally, anonymous functions can be attached to Go variables.
3. An anonymous function in Go has access to the variables defined in the context where it is itself defined. This property of anonymous functions is called `closure`.
4. It is considered a good practice for anonymous functions to have a small implementation and a local focus. If an anonymous function does not have a local focus, then we might need to consider making it a regular function.

`Functions that return multiple values`

Functions can return multiple distinct values, which saves us from having to create a dedicated structure for returning and receiving multiple values from a function. However, if we have a function that returns more than three values, we should reconsider that decision and maybe redesign it to use a single structure or slice for grouping and returning the desired values as a single entity—this makes handling the returned values simpler and easier.

`Naming Return Values in Go Functions`

Go allows us to name the return values of a Go function. Additionally, when such a function has a return statement without any arguments, the function automatically returns the current value of each named return value in the order in which they were declared in the function signature.

`Functions That Accept Other Functions as Parameters` - Functions can accept other functions as parameters. The best example of a function that accepts another function as an argument can be found in the sort package. We can provide the sort.Slice() function with another function as an argument that specifies the way sorting is implemented - `func Slice(slice interface{}, less func(i, j int) bool)`

`Functions That Return Other Functions` - Apart from accepting functions as arguments, functions can also return anonymous functions, which can be handy when the returned function is not always the same but depends on the function's input or other external parameters.

`Variadic Functions` - Variadic functions are functions that can accept a variable number of parameters — we already know about fmt.Println() and append(), which are both variadic functions that are widely used. In fact, most functions found in the fmt package are variadic - `func variadic(input ...string{})`

**Passing os.Args as ...interface{}** - If we try to pass os.Args, which is a slice of strings ([]string), as ...interface{} to a variadic function, our code will not compile and will generate an error message similar to cannot use `os.Args (type []string) as type []interface {} in argument to <function_name>`. This happens because the two data types ([]string and []interface{}) do not have the same representations in memory—this applies to all data types. In practice, this means that we can’t use os.Args... to pass each individual value of the os.Args slice to a variadic function - The solution to this problem is converting the slice of strings—or any other slice—into a slice of interface{}

`The defer Keyword`

The defer keyword postpones the execution of a function until the surrounding function returns. It is very important to remember that deferred functions are executed in last in, first out (LIFO) order after the surrounding function has been returned.
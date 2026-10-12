# Basic Go Data Types

`Error Data Type`

1. Go provides a special data type for representing error conditions and error messages named `error` — in practice, this means that Go treats errors as values
2. `The next convention for error messages` - if the value of an error variable is nil, then there was no error.
3. `Customized error messages` - Should we wish to return a custom error, we can use `errors.New()` from the errors package. If we want to format our error messages in the way fmt.Printf() works, we can use the fmt.Errorf() function, which simplifies the creation of custom error messages—the `fmt.Errorf()` function returns an error value just like errors.New().

`Numeric Data Types`

int, unit, float32, flaot64, complex64, complex128.

`Non-Numeric Data Types`

**String**
1. Go has support for strings, characters, runes, dates, and times. However, Go does not have a dedicated char data type.
2. Go supports the `string `data type for representing strings. A Go string is just a collection of bytes and can be accessed as a whole or as an array.
3. A rune is an int32 value that is used for representing a single Unicode code point, which is an integer value that is used for representing single Unicode characters.
4. Although a rune is an int32 value, we can’t compare a rune with an int32 value. Go considers these two data types as totally different. While both types have the same underlying int32 representation, the rune type is specifically used to store Unicode characters. This means that a rune value holds a 32-bit signed integer that represents a single Unicode character, distinguishing it from a generic int32 value.
5. We can create a new byte slice from a given string by using a `[]byte("A String")` statement. Given a byte slice variable b, we can convert it into a string using the `string(b)` statement
6. When working with byte slices that contain Unicode characters, the number of bytes in a byte slice is not always connected to the number of characters in the byte slice because most Unicode characters require more than one byte for their representation. As a result, when we try to print each single byte of a byte slice using fmt.Println() or fmt.Print(), the output is not text presented as characters but integer values. If we want to print the contents of a byte slice as text, we should either print it using `string(byteSliceVar)` or using `fmt.Printf() with %s` to tell fmt.Printf() that we want to print a string.

**Rune**
1. We can define a rune using single quotes: `r := '€'`, and we can print the integer value of the bytes that compose it as `fmt.Println(r)` — in this case, the integer value is 8364. Printing it as a single Unicode character requires the use of the `%c control string in fmt.Printf()`. The `%c` control string in fmt.Printf() prints a rune as a character and 
2. The `%x` control string in fmt.Printf() prints the string as a rune.
```go
aString := "Hello World! €"
// Print an existing string as runes
for _, v := range aString {
    fmt.Printf("%x ", v)
}
fmt.Println()
```

`Integer Conversions`

1. We can convert an integer value into a string in two main ways: using string() and using a function from the strconv package. However, the two methods are fundamentally different.
2. The `string()` function converts an integer value into a `Unicode code point, which is a single character`. Means that when you use the string() function on an integer, it doesn’t turn the number into a string of digits but instead treats the integer as a Unicode code point and converts it into the corresponding single character. For example, string(65) produces “A” because 65 is the Unicode code for the letter A, so the result is one character, not the string “65”. 
3. whereas functions such as `strconv.FormatInt()` and `strconv.Itoa()` converts an integer value into a string value with the same representation and the same number of characters.

`The unicode package`

The unicode standard Go package contains various handy functions for working with Unicode code points. One of them is called unicode.IsPrint(), and it can help us identify the parts of a string that are printable using runes.

```go
package main

import (
	"fmt"
	"unicode"
)

func main() {
	const sL = "\x99\x00ab\x50\x00\x23\x50\x29\x9c"

	for i := 0; i < len(sL); i++ {
		if unicode.IsPrint(rune(sL[i])) {
			fmt.Printf("%c\n", sL[i])
		} else {
			fmt.Println("Not printable!")
		}
	}
}
```

`The String Package`

1. `String Comparisons` - The strings.EqualFold() function compares two strings `without considering their case` and returns true when they are the same and false otherwise. Ex: `s.EqualFold("Mihalis", "MIHAlis")`
2. `String Searches`
    1. The `strings.Index()` function checks whether the string of the second parameter can be found in the string that is given as the first parameter and returns the index where it was found for the first time. On an unsuccessful search, it returns -1. Ex: s.Index("Mihalis", "ha")
    2. The `strings.HasPrefix()` function checks whether the given string, which is the first parameter, begins with the string that is given as the second parameter. Similarly, the `strings.HasSuffix()` function checks whether the given string ends with the second string.
3. `Tokenization` - Tokenization is the process of breaking a string or text into smaller pieces called tokens, which can be words, phrases, or symbols. The handy `strings.Fields()` function splits the given string around one or more white space characters as defined by the unicode.IsSpace() function and returns a slice of substrings found in the input string. If the input string contains white characters only, it returns an empty slice.
4. `Splits and Replacement`
    1. The s`trings.Split()` function allows us to split the given string according to the desired separator string — the strings.Split() function returns a string slice. Using `""` as the second parameter of strings.Split() allows us to process a string character by character.
    2. The `strings.Replace()` function takes four parameters. The first parameter is the string that we want to process. The second parameter contains the string that, if found, will be replaced by the third parameter of strings.Replace(). The last parameter is the maximum number of replacements that are allowed to happen.
    3. The `strings.SplitAfter()` function splits its first parameter string into substrings based on the separator string that is given as the second parameter to the function. The separator string is included in the returned slice.

`Time and Date`

1. Working with times and dates in Go is the `time.Time` data type, which represents an instant in time with `nanosecond precision`. Each time.Time value is associated with a `location (time zone)`.
2. The `time.Now().Unix()` function returns the popular UNIX epoch time, which is the number of seconds that have elapsed since 00:00:00 UTC, January 1, 1970.
3. Parsing a time string - The single most important topic about Go and dates and times is the way Go parses a string in order to convert it into a date and a time. The reason that this is important is usually that such input is given as a string and not as a valid date variable. The function used for parsing is `time.Parse()`, and its full signature is Parse(layout, value string) (Time, error), where layout is the parse string and value is the input that is being parsed. The time.Time value that is returned is a moment in time with nanosecond precision and contains both date and time information

`Go Constants`

1. Go supports constants, which are variables that cannot change their values. Constants in Go are defined with the help of the const keyword. Generally speaking, constants can be either global or local variables.
2. Strictly speaking, the value of a constant variable is defined at compile-time, not at runtime—this means that it is included in the binary executable.
3. `The constant generator iota` - The constant generator iota is used for declaring a sequence of related values that use incrementing numbers without the need to explicitly type each one of them.

`Grouping Similar Data: Arrays and Slices`

1. When to use arrays and slices - The quick answer is that we can use slices instead of arrays almost anywhere in Go, but we are also demonstrating arrays because they can still be useful and because slices are implemented by Go using arrays
2. Arrays:
    1. When defining an array variable, we must define its size. Otherwise, we should put [...] in the array declaration and let the Go compiler find out the length for us. So, we can create an array with four string elements, either as `[4]string{"Zero", "One", "Two", "Three"}` or as `[...]string{"Zero", "One", "Two", "Three"}`. If we put nothing in the square brackets, then a slice is going to be created instead. 
    2. We can’t change the size of an array after we have created it.
    3. When we pass an array to a function, what happens is that Go creates a copy of that array and passes that copy to that function—therefore, any changes we make to an array inside a function are lost when the function returns.
3. Slices:
    1. Slices in Go are more powerful than arrays mainly because they are dynamic, which means that they can grow or shrink after creation if needed. Additionally, any changes we make to a slice inside a function also affect the original slice because they are pointing toward the same memory.
    2. In reality, a slice value is a header that contains the following - A pointer to an underlying array where the elements are actually stored, The length of the array, The array’s capacity.
    3. Note that the slice value does not include its elements, just a pointer to the underlying array. So, when we pass a slice to a function, Go makes a copy of that header and passes it to the function. This copy of the slice header includes the pointer to the underlying array. A side effect of passing the slice header is that it is `faster to pass a slice to a function` because Go does not need to make a copy of the slice and its elements, just the slice header.
    4. Slice creation - We can create a slice using make() or like an array without specifying its size or using [...]. If we do not want to initialize a slice, then using make() is better and faster.
3. `Slice Operations` - We can find the length of an array or a slice using len(). We can add new elements to a full slice using the append() function. The append() function automatically allocates the required memory space.
4. `Slice Len and Cap` - Both arrays and slices support the len() function for finding out their length. However, slices also have an additional property called capacity that can be found using the cap() function. Although the capacity of a slice is handled by Go after slice creation, a developer can define the capacity of a slice at creation time using the `make()` function—after that, the capacity of the slice doubles each time the length of the slice is about to become bigger than its current capacity. `a := make([]int, 4)` -  Only length is defined. Capacity = length. `a := make([]int, 4)` == `a := make([]int, 4,4)`

`Selecting a Part of a Slice`

In Go, we select a part of a slice by defining two indexes: the first one is the beginning of the selection, whereas the second one is the end of the selection, without including the element at that index, separated by `:`. For Example, If we want to process all the command-line arguments of a utility apart from the first one, which is its name, we can assign it to a new variable (arguments := os.Args) for ease of use and use the `arguments[1:]` notation to skip the first command-line argument.
However, there is a variation where we can add a third parameter that controls the capacity of the resulting slice. So, using aSlice[0:2:4] selects the first two elements of a slice (at indexes 0 and 1) and creates a new slice with a maximum capacity of 4. The resulting capacity is defined as the result of the 4-0 subtraction where 4 is the maximum capacity, and 0 is the first index—if the first index is omitted, it is automatically set to 0. In this case, the capacity of the result slice will be 4 because 4-0 equals 4

`Deleting an Element from a Slice`

Deleting an element from a slice can be tricky, we present two techniques for doing so. 
1. The first technique virtually divides the original slice into two slices, split at the index of the element that needs to be deleted. Neither of the two slices includes the element that is going to be deleted. After that, we concatenate these two slices and create a new one. 
2. The second technique copies the last element in the place of the element that is going to be deleted and creates a new slice by excluding the last element from the original slice.

`Sorting Slices`

The sort package can sort slices of built-in data types without the need to write any extra code. Additionally, Go provides the `sort.Reverse()` function for sorting in reverse rather than the default order. However, what is really interesting is that sort allows us to write our own sorting functions for custom data types by implementing the sort.Interface interface.

So, we can sort a slice of integers saved as x by typing `sort.Ints(x)`. When sorting a slice of integers in reverse order using sort.Reverse(), we need to pass the desired slice to sort.Reverse() using sort.IntSlice(x) because the IntSlice type implements the sort.Interface internally, which allows us to sort in a different way than usual. The same applies to the other standard Go data types.

`Pointers`

1. A pointer is the memory address of a variable. We need to dereference a pointer in order to get its value— dereferencing is performed using the * character in front of the pointer variable.
2. Additionally, we can get the memory address of a normal variable using an & in front of it.
3. If a pointer variable points to an existing regular variable, then any changes we make to the stored value using the pointer variable will modify the regular variable.
4.  `zero value` for pointers in `nil`
5. `Reasons for using pointers`
    1. The main benefit we get from pointers is that passing a variable to a function as a pointer does not discard any changes we make to the value of that variable inside that function when the function returns. Remember that slices are passed to functions without the need to use a pointer—it is Go that passes the pointer to the underlying array of a slice, and there is no way to change that behavior.
    2. Pointers allow us to share data between functions. However, when sharing data between functions and goroutines, we should be extra careful with race condition issues.
    3. Pointers are also very handy when we want to tell the difference between the zero value of a variable and a value that is not set (nil). This is particularly useful with structures because pointers to structures can have the nil value, which means that we can compare a pointer to a structure with the nil value, which is not allowed for normal structure variables.
    4. Having support for pointers and, more specifically, pointers to structures allows Go to support data structures such as linked lists and binary trees, which are widely used in computer science.
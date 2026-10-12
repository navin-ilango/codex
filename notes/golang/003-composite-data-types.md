# Composite data types

`Maps`

1. Both arrays and slices limit us to using positive integers as indexes. Maps are powerful data structures because they allow us to use indexes of various data types as keys to look up our data as long as these keys are comparable. 
2. A practical rule is that we should use a map when we are going to need indexes that are not positive integer numbers or when the integer indexes have big gaps.
3. We can create a new map variable using either make() or a map literal. Creating a new map with string keys and int values using make() is as simple as writing make(map[string]int) and assigning its return value to a variable - `m := make(map[string]int)`
4. On the other hand, if we decide to create a map using a map literal, we need to write something like the following: `m := map[string]int {"key1": -1, "key2": 123}`
5. We should make no assumptions about the order of the elements inside a map. Go randomizes keys when iterating over a map—this is done on purpose and is an intentional part of the language design.
6. `Map operations` - len(), delete().
7. Additionally, we can tell whether a key k exists on a map named aMap by the second return value of the `v, ok := aMap[k]` statement. If ok is set to true, then k exists, and its value is v. If it does not exist, v will be set to the zero value of its data type, which depends on the definition of the map
8. If we try to get the value of a key that does not exist in a map, Go will not complain about it and will return the zero value of the data type of the value

`Storing to a Nil Map`

Now, let’s discuss a special case where a map variable has the nil value. We are allowed to assign a map variable to nil. In that case, we will not be able to use that variable until we assign it to a new map variable. Put simply, if we try to store data on a nil map, our program will crash.

Testing whether a map points to nil before using it is a good practice.

`Iterating over Maps`

When for is combined with the range keyword, it implements the functionality of foreach loops found in other programming languages and allows us to iterate over all the elements of a map without knowing its size or its keys. When range is applied on a map, it returns key and value pairs in that order - `for k, v := range sampleMap {}`

`Go Structures`

1. `Defining a Structure` - When we define a new structure, we group a set of values into a single data type, which allows us to pass and receive this set of values as a single entity. A structure has fields, and each field has its own data type, which can even be another structure or slice of structures. Additionally, because a structure is a new data type, it is defined using the `type` keyword followed by the name of the structure and ending with the `struct` keyword, which signifies that we are defining a new structure.
2. `Using the new keyword` - Additionally, we can create new structure instances using the new() keyword: `pS := new(Entry)`. The new() keyword has the following properties: It allocates the proper memory space, which depends on the data type, and then it zeroes it, It always returns a pointer to the allocated memory, It works for all data types except channels and maps.
3. `Type identity` -  The order in which we put the fields in the definition of a structure type is significant for the type identity of the defined structure. Put simply, two structures with the same fields will not be considered identical in Go if their fields are not in the same order.

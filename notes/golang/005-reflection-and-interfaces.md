# Reflection and Interfaces

One might wonder how we can find out the names of the fields of a structure at execution time. In such cases, we need to use reflection.

`The reflect package`

1. Reflection allows us to dynamically learn the type of an arbitrary object along with information about its structure.
2. Reflection allows us to handle and work with data types that do not exist at the time at which we write our code but might exist in the future, which happens when we set up our code to use an external package with user-defined data types.
3. reflect.Value, reflect.Type, reflect.TypeOf(), reflect.ValueOf(), reflect.NumField(), Field(), reflect.Kind, Kind(), the Int() and String() methods return the integer and string value of a reflect.Value, respectively.

`Type Methods`

A type method is a function that is attached to a specific data type. Although type methods (or methods on types) are, in reality, functions, they are defined and used in a slightly different way.

```
func (t type) FunctionName(parameters) <return values> {}
```

Keep in mind that under the hood, the Go compiler does turn methods into regular function calls with the self value as the first parameter. However, interfaces require the use of type methods to work.

`Interfaces`

An interface is a Go mechanism for defining behavior that is implemented using a set of methods. Interfaces play a key role in Go and can simplify the code of our programs when they have to deal with multiple data types that perform the same task

Interfaces work with methods on types (or type methods), which are like functions attached to given data types, which in Go are usually structures (although we can use any data type we want).

Once we implement the required type methods of an interface, that interface is satisfied implicitly.

The empty interface is defined as just interface{}. As the empty interface has no methods, it means that it is already implemented by all data types.

The biggest advantage we get from interfaces is that, if needed, we can pass a variable of a data type that implements a particular interface to any function that expects a parameter of that specific interface, which saves us from having to write separate functions for each supported data type. However, Go offers an alternative to this with the recent addition of generics.

Ex: sort.Interface - len(), less(), swap() are the methods to implement.

`The Empty Interface`

The empty interface is defined as just interface{} and is already implemented by all data types. Thus, variables of any data type can be put in the place of a parameter of the empty interface data type. Therefore, a function with an interface{} parameter can accept variables of any data type in this place.

`Type Assertions` - A type assertion is used to extract the underlying concrete value from an interface value. It answers the question: "I have an interface, but what specific type is actually hiding inside it?" `value, ok := x.(T)` & `	switch t := x.(type)`

`The error Data Type` 

the error data type, which is an interface defined as follows:
```go
type error interface {
    Error() string
}
```

In order to satisfy the error interface, we just need to implement the Error() string type method.

`Object-oriented concepts in Go`


**Encapsulation** - Go rejects access modifiers like `public` or `private`. Instead, visibility is entirely determined by **letter casing** at the **package level**.

*   ** Uppercase First Letter** -> **Public** (Exported outside the package).
*   ** Lowercase First Letter** -> **Private** (Unexported, hidden inside the package).



**Composition (Struct Embedding)** - Go does not have a structural `extends` keyword. It favors **composition over inheritance** using **Struct Embedding**. You drop an anonymous struct inside another, and its capabilities are instantly promoted.

*   **"Has-A" instead of "Is-A"**: A car *has an* engine, it is not an engine.
*   **Method Promotion**: The outer struct gains direct access to all fields and methods of the inner struct.


**Polymorphism (Implicit Interfaces)** - Polymorphism in Go is achieved through **Interfaces**. Go interfaces are **implicitly satisfied**, meaning you never use an `implements` keyword. 

*   **The Duck Test**: If a struct implements all the methods defined in an interface, Go automatically accepts it as that interface type.
*   **Decoupled Design**: Creators of structs don't need to know about the interfaces that will eventually consume them.

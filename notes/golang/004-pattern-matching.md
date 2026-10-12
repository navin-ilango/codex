# Pattern Matching

Pattern matching is a technique for searching a string for some set of characters based on a specific search pattern that is based on regular expressions and grammar

`The regexp package`
1. The Go package responsible for defining regular expressions and performing pattern matching is called `regexp`.
2. We use the `regexp.MustCompile()` function to create the regular expression and the `Match()` function to see whether the given string is a match or not. 
3. The regular expression `^[A-Z][a-z]*$` - Matches Names or SurNames
4. Matching Integers - `^[-+]?\d+$` - In plain English, what we say here is that we want to match something that begins with – or +, which is optional (?) and ends with any number of digits (\d+)—it is required that we have at least one digit before the end of the string that is examined ($).
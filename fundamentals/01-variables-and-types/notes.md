# 01 - Basics, Variables, and Data Types

## Key Takeaways

### 1. Variables and Scope
* **`:=` (Short Declaration):** Can only be used **inside** functions. At the package level, you must use `var`.
* **Unused Variables:** Go's compiler will throw an error if you declare a variable and never use it. This keeps code clean and prevents dead code.
* **Global vs. Local:** Package-scoped variables are accessible anywhere in that package, whereas function-scoped variables die when the function exits.

### 2. Strict Typing & Conversion
* Go is a **statically typed** language. 
* Implicit type conversion does not exist. If you have an `int` and want to add it to a `float64`, you must explicitly convert the integer: `float64(myInt) + myFloat`.

### 3. Zero Values
Unlike languages that initialize variables to `null` or leave them undefined (causing bugs), Go initializes every variable to its **zero value** automatically:
* `int` / `float`: `0` or `0.0`
* `string`: `""` (empty string)
* `bool`: `false`
* Pointers, slices, maps, channels, interfaces: `nil`

### 4. Common Formatting Verbs (`fmt.Printf`)
* `%v`: The value in a default format
* `%T`: The type of the value
* `%d`: Base 10 integers
* `%f`: Floating-point numbers
* `%s`: Plain string
* `%t`: Booleans (`true` or `false`)
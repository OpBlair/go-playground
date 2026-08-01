# 03 - Functions

## Key Takeaways

### 0. Function Naming Convention & Visibility

* **Exported (Public):** Functions starting with an **uppercase letter** (conventionally written in `PascalCase`) are visible and accessible outside the package.
* **Unexported (Private):** Functions starting with a **lowercase letter** (conventionally written in `camelCase`) are restricted to the package where they are defined.

### 1. Function Declarations & Signatures

* Syntax: Declared using the func keyword, followed by the name, parameters with explicit types, and the return type(s).

* Parameter Shorthand: If consecutive parameters share the same type, you only need to specify the type once (e.g., func add(a, b int) int).

* Pass by Value: Go is strictly pass-by-value. When arguments are passed into a function, copies are created; modifying them inside the function does not affect the original variables.

### 2. Multiple & Named Return Values

* Multiple Returns: Go natively supports returning multiple values, most commonly used to return a result alongside an error interface.

* Named Returns: You can name return values in the function signature. They act as pre-declared local variables.

* Naked Returns: A return statement without arguments will return the current values of the named return variables. Use sparingly to maintain readability.

### 3. Variadic Functions

* Variable Arguments: Using the ellipsis (...) operator allows a function to accept an arbitrary number of arguments of a specific type.

* Slice Behavior: Inside the variadic function, the parameters are handled as a slice of that type, making it easy to iterate over them using range.

### 4. Anonymous Functions & Closures

* First-Class Citizens: Functions in Go can be assigned to variables, passed as arguments, or returned from other functions.

* Closures: Anonymous functions can bind to and reference variables from outside their immediate body, capturing and maintaining state across executions.

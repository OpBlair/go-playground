# 05 - Pointers

## Key Takeaways

### 1. What is a Pointer?
* **Memory Address:** A pointer is a variable that stores the **memory address** of another value, rather than the value itself.
* **Zero Value:** The zero value of an uninitialized pointer is `nil`.

### 2. The Two Crucial Operators
* **`&` (Address-of Operator):** Prefixed to a variable, it returns the memory address where that variable lives (e.g., `&x`).
* **`*` (Dereference Operator):** Prefixed to a pointer variable, it accesses or modifies the *value* stored at that memory address (e.g., `*ptr`).

### 3. Pass-by-Value vs. Pass-by-Reference
* **Go is Strictly Pass-by-Value:** When you pass a variable into a function, Go copies it.
* **Simulating Pass-by-Reference:** By passing a *pointer* (`*int`) into a function, you pass a copy of the *memory address*. This allows the function to mutate the original variable outside its scope.
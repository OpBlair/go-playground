# 06 - Structs & Methods

## Key Takeaways

### 1. Structs (Custom Data Types)
* **Aggregation:** Structs allow you to group variables (fields) together under a single custom type to represent real-world entities.
* **Initialization:** Can be initialized using positional fields or named field syntax (which is much safer and preferred).

### 2. Methods & Receiver Functions
* **Binding Logic:** In Go, methods are just functions with a special receiver argument placed between the `func` keyword and the method name.
* **Value Receivers (`func (t Type)`):** Operates on a *copy* of the struct. Use this when the method doesn't need to mutate state.
* **Pointer Receivers (`func (t *Type)`):** Operates on a *pointer* to the struct. Use this when the method needs to modify the struct's fields or to avoid copying large structs.
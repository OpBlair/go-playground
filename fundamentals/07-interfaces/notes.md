# 07 - Interfaces

## Key Takeaways

### 1. Implicit Implementation
* **No `implements` Keyword:** A type implements an interface automatically simply by implementing its method signatures.
* **Decoupling:** Write functions that accept interfaces rather than concrete structs to make your code modular, extensible, and testable.

### 2. The Empty Interface (`any` / `interface{}`)
* **Accepts Anything:** An empty interface has zero methods, meaning any type satisfies it.
* **Type Assertions:** To extract the underlying concrete type safely, use the idiom: `val, ok := myVar.(TargetType)`.


# Interfaces in Go

## What is an interface?

An interface defines a set of method signatures.
A type implements an interface implicitly by having all the required methods.

## Why interfaces matter

Interfaces let you write flexible and reusable code.
They help separate what something does from what it is.

## Example

```go
type Speaker interface {
    Speak() string
}

type Dog struct{}

func (d Dog) Speak() string {
    return "woof"
}
```

A `Dog` implements `Speaker` because it has a `Speak()` method.

## Key ideas

- Interfaces are satisfied implicitly.
- There is no `implements` keyword in Go.
- A variable of interface type can hold any value that implements that interface.
- Interfaces are useful for polymorphism and testing.

## Empty interface

The empty interface is `interface{}`.
It can hold any value because every type implements zero methods.

Example:

```go
var x interface{}
x = 42
x = "hello"
x = true
```

## Type assertion

A type assertion lets you extract the concrete value from an interface.

```go
value, ok := x.(string)
```

If `ok` is `true`, the assertion succeeded.

## Type switch

A type switch checks the concrete type stored in an interface.

```go
switch v := x.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
default:
    fmt.Println("unknown type")
}
```

## Best practices

- Keep interfaces small.
- Define interfaces where they are used, not where they are implemented.
- Prefer specific interfaces with one or a few methods.
- Use `interface{}` only when you truly need to accept any type.

## Common mistake

A large interface is harder to reuse and test.
Small interfaces are easier to understand and compose.
# 08 - Error Handling

## Key Takeaways

### 1. Errors are Values
* **No Exceptions:** Go avoids `try/catch` blocks entirely. Errors are ordinary return values implementing the built-in `error` interface.
* **The `if err != nil` Idiom:** Explicit error checking right after function execution is standard practice in Go.

### 2. Advanced Error Patterns
* **Custom Errors:** Use `errors.New("msg")` or `fmt.Errorf("...")` to create static or dynamic errors.
* **Error Wrapping (`%w`):** Wrap lower-level errors with context while keeping the original error intact for unwrapping or inspection.

# Error Handling in Go

## What is an error?

An error is a value that represents a failure or unexpected condition.
In Go, errors are handled explicitly using the built-in `error` interface.

## The error interface

The `error` interface has one method:

```go
type error interface {
    Error() string
}
```

Any type that implements `Error() string` can be used as an error.

## Returning errors

Go functions usually return a value and an error.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}
```

Always check the error before using the result.

## Checking errors

Typical pattern:

```go
result, err := divide(10, 0)
if err != nil {
    fmt.Println("error:", err)
    return
}
fmt.Println("result:", result)
```

## Creating errors

### Using `errors.New`

Use `errors.New("message")` for simple errors.

### Using `fmt.Errorf`

Use `fmt.Errorf` when you want formatted error messages.

```go
err := fmt.Errorf("invalid input: %d", value)
```

## Custom errors

You can define your own error type.

```go
type MyError struct {
    Message string
}

func (e MyError) Error() string {
    return e.Message
}
```

Custom errors are useful when you want to attach extra data.

## Wrapping errors

Wrapping adds context while keeping the original error.

```go
return fmt.Errorf("reading config: %w", err)
```

The `%w` verb wraps the error so it can be inspected later.

## Unwrapping and inspecting

Wrapped errors can be checked with:

- `errors.Is` for matching a specific sentinel error.
- `errors.As` for matching a specific error type.

## Panic and recover

`panic` stops normal execution immediately.
`recover` can be used inside deferred functions to regain control after a panic.

Use them sparingly.
They are meant for truly exceptional situations, not normal error handling.

## Best practices

- Return errors instead of panicking in normal code.
- Add context when wrapping errors.
- Check errors immediately.
- Keep error messages clear and specific.
- Use custom errors only when they add real value.

## Common mistake

Ignoring errors can cause bugs that are hard to trace.
Always handle them explicitly.

## Differences
- **Sentinel Errors (errors.New):** Good for simple, static checks
- **Formatted Errors (fmt.Errorf):** Good for adding wrapper context along the call stack.
- **Custom Error Types (struct + Error() string):** Good for bundling extra structured data(like error codes or custom fields) inside the error itself.
- **errors.As (Type Assertion for Errors):** While errors.Is checks if an error matches a specific sentinel value, errors.As checks if an error matches a specific type. This is heavily used with custom error structs when you need to extract underlying fields or data from an error.  
- **Custom Error Types:** Because the built-in error is just an interface requiring an **Error() string** method, you can create your own structs that implement this interface. This allows your errors to carry rich metadata (like error codes, timestamps, or field names).  
- **errors.Join:** Introduced in Go 1.20, this allows you to combine multiple errors into a single error, which is useful when validating multiple fields at once or gathering multiple failures (e.g., during configuration loading).  

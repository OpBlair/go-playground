# 02 - Control Flow

## Key Ideas

### 1. `if`, `else if`, `else`

* `if` requires a **boolean** condition; there is no truthy/falsy concept in Go.
* Syntax:
  ```go
  if condition {
      // runs when condition is true
  } else if otherCondition {
      // runs when first is false, second is true
  } else {
      // runs when all above are false
  }
  ```
* You can use a **short statement** before the condition, scoped to the `if`:
  ```go
  if x := compute(); x > 10 {
      // x exists only inside this if/else chain
  }
  ```

### 2. Logical operators and precedence

* Operators:
  * `&&` – logical AND
  * `||` – logical OR
  * `!` – logical NOT
* `&&` binds tighter than `||`, so:
  ```go
  a || b && c   // same as: a || (b && c)
  ```
* Use parentheses for clarity in complex conditions:
  ```go
  if (age >= 18 && hasID) || isParent {
      // ...
  }
  ```

### 3. `switch` statements

* `switch` compares a single expression against multiple `case` values:
  ```go
  switch value {
  case 1:
      // ...
  case 2, 3:
      // multiple values in one case
  default:
      // runs if no case matches
  }
  ```
* Key differences from some other languages:
  * **No implicit fall-through**; once a case runs, the switch ends.
  * Use `fallthrough` explicitly if you want the next case to run:
    ```go
    switch n {
    case 1:
        fmt.Println("one")
        fallthrough
    case 2:
        fmt.Println("one or two")
    }
    ```
* You can write a **switch without an expression**, which acts like a cleaner `if/else if` chain:
  ```go
  switch {
  case score >= 90:
      fmt.Println("A")
  case score >= 80:
      fmt.Println("B")
  default:
      fmt.Println("Below B")
  }
  ```

### 4. `for` loops (Go’s only loop)

* Go has only `for`; there is no separate `while` or `do...while`.
* Classic form:
  ```go
  for i := 0; i < 10; i++ {
      // ...
  }
  ```
* While-style loop:
  ```go
  for condition {
      // runs while condition is true
  }
  ```
* Infinite loop:
  ```go
  for {
      // must break or return inside, or it runs forever
  }
  ```
* Control statements:
  * `break` – exit the nearest loop or `switch`.
  * `continue` – skip to the next iteration of the loop.

### 5. Looping over collections

* `for range` is used to iterate over arrays, slices, maps, strings, etc.:
  ```go
  nums := []int{1, 2, 3}
  for i, v := range nums {
      // i = index, v = value
  }
  ```
* If you only need the value:
  ```go
  for _, v := range nums {
      // ignore index with _
  }
  ```
* For maps, the order of iteration is **not guaranteed**:
  ```go
  m := map[string]int{"a": 1, "b": 2}
  for k, v := range m {
      // k = key, v = value
  }
  ```

### 6. Nested control flow and readability

* Control flow constructs (`if`, `switch`, `for`) can be nested, but deeply nested logic quickly becomes hard to read.
* Idiomatic Go favors:
  * Early returns to avoid deep nesting.
  * Clear, simple conditions.
  * Small functions that each handle one responsibility.

### 7. Common patterns to recognize

* **Guard clauses**:
  ```go
  if err != nil {
      return err
  }
  // main logic continues here
  ```
* **State-based `switch`** (e.g., on status or type):
  ```go
  switch status {
  case "pending":
      // ...
  case "done":
      // ...
  default:
      // ...
  }
  ```
* **Loop with exit condition inside**:
  ```go
  for {
      if done {
          break
      }
      // work
  }
  ```
# 04 - Collections (Arrays, Slices, and Maps)

## Key Takeaways

### 1. Arrays (Fixed Size)
* **Fixed Length:** An array’s size is part of its type. For example, `[3]int` and `[5]int` are completely different types.
* **Value Semantics:** Arrays are values in Go. Assigning one array to another, or passing it into a function, copies the entire array.
* **Rarely Used Directly:** Because of their rigidity, arrays are less common in everyday Go code than slices.


### 2. Slices (Dynamic Views)
* **The Workhorse:** Slices are the most commonly used collection type in Go. They are built on top of arrays and provide a flexible view into the underlying data.
* **Length vs. Capacity:**
  * `len(s)` = the number of elements currently in the slice.
  * `cap(s)` = the total number of elements the slice can hold before a new underlying array is allocated.
* **Creation:** Slices can be created with literals or `make`.
  * `s := []int{1, 2, 3}`
  * `s := make([]int, length, capacity)`
* **Appending:** Use `append` to grow a slice. If capacity is exceeded, Go allocates a new underlying array automatically.
* **Copying Slices:** Use `copy(dst, src)` when you want to duplicate elements into another slice.


### 3. Maps (Key-Value Pairs)
* **Unordered Lookups:** Maps are hash tables used for fast key-value access.
* **Comparable Keys:** The key type must be comparable, such as `string`, `int`, or certain structs.
* **Initialization Required:** Always initialize a map before writing to it using `make` or a map literal.
  * `m := make(map[string]int)`
  * `m := map[string]int{"one": 1, "two": 2}`
* **Comma-Ok Idiom:** Use the comma-ok pattern to safely check whether a key exists.
  ```go
  val, exists := myMap["key"]
  if exists {
      // key is present
  }
  ```
* **Deleting Keys:** Use `delete(m, key)` to remove an entry from a map.


### 4. Common Collection Patterns
* Use **arrays** when the size is known and fixed.
* Use **slices** for most ordered collections.
* Use **maps** for lookups by key.
* Combine them often, like `[]Person` or `map[string]int`.
* Loop over slices and maps using `range`.


### 5. Gotchas
* Appending to a slice may create a new underlying array.
* A nil slice can still be appended to.
* A nil map cannot be written to until it is initialized.
* Copying an array copies all of its elements.
* Map iteration order is not guaranteed.
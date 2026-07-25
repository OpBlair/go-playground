package main

import (
	"fmt"
	"reflect"
)

var globalAppVersion = "v1.0.0"

const (
	StatusOK    = 200
	StatusError = 500
	MaxRetries  = 3
)

func main() {
	var username string = "gopher_dev"
	var birthYear = 2009
	isAwesome := true

	var firstName, lastName string = "John", "Doe"
	x, y, z := 1, 2, 3

	fmt.Println("--- 1. Variables ---")
	fmt.Printf("User: %s, Born: %d, Awesome: %t\n", username, birthYear, isAwesome)
	fmt.Printf("Name: %s %s | Coords: %d, %d, %d\n", firstName, lastName, x, y, z)
	fmt.Printf("Global App Version: %s\n\n", globalAppVersion)

	var defaultInt int
	var defaultFloat float64
	var defaultString string
	var defaultBool bool

	fmt.Println("--- 2. Zero Values ---")
	fmt.Printf("int: %d | float: %f | string: %q | bool: %t\n\n", defaultInt, defaultFloat, defaultString, defaultBool)

	var integerNum int = 42
	var floatNum float64 = 3.14

	result := float64(integerNum) + floatNum

	fmt.Println("--- 3. Type Conversion ---")
	fmt.Printf("Converted sum: %f\n\n", result)

	fmt.Println("--- 4. Constants ---")
	fmt.Printf("Status OK: %d, Max Retries: %d\n\n", StatusOK, MaxRetries)

	sampleVar := 99.99

	fmt.Println("--- 5. Formatting & Types ---")
	fmt.Printf("Value: %v\n", sampleVar)
	fmt.Printf("Type via %%T: %T\n", sampleVar)
	fmt.Printf("Type via reflect: %v\n", reflect.TypeOf(sampleVar))

	var name string = "Alice"
	age := 30
	const planet = "Earth"

	isLearning := true
	pi := 3.14159

	var otherDefaultInt int
	var otherDefaultString string

	fmt.Printf("Name: %v (Type: %T)\n", name, name)
	fmt.Printf("Age: %v (Type: %T)\n", age, age)
	fmt.Printf("Planet: %v\n", planet)
	fmt.Printf("Is Learning? %v\n", isLearning)
	fmt.Printf("Pi: %f\n", pi)

	fmt.Printf("Default Int (Zero value): %v\n", otherDefaultInt)
	fmt.Printf("Default String (Zero value): %q\n", otherDefaultString)
}

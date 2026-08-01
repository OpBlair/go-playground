package main

import "fmt"

// Define the struct
type Person struct {
	Name string
	Age  int
}

type Counter struct {
	Value int
}

// 1. Value Receiver (just reads the value)
func (c Counter) PrintValue() {
	fmt.Printf("Current Value: %d\n", c.Value)
}

// 2. Pointer Receiver (modifies the actual struct instance)
func (c *Counter) Increment() {
	c.Value++ // Go automatically dereferences c behind the scenes!
}

// 3. Struct Embedding (Go's Version of Inheritance)
type Engine struct {
	Horsepower int
	Type       string
}

type Car struct {
	Brand  string
	Model  string
	Engine // Embedded struct (fields and methods are promoted!)
}

func main() {
	// 1. Struct Basics & Initialization
	// Preferred way: Named field initialization
	p1 := Person{Name: "Alice", Age: 28}

	// Accessing fields using dot notation (works for values and pointers alike!)
	fmt.Println("Name:", p1.Name)
	fmt.Println("Age:", p1.Age)

	// 2. Value Receivers vs. Pointer Receivers
	myCounter := Counter{Value: 0}

	myCounter.PrintValue() // Current Value: 0

	myCounter.Increment() // Modifies it via pointer receiver
	myCounter.Increment() // Modifies it again

	myCounter.PrintValue()

	// 3. Struct Embedding (Go's Version of Inheritance)
	myCar := Car{
		Brand: "Toyota",
		Model: "Supra",
		Engine: Engine{
			Horsepower: 382,
			Type:       "Inline-6",
		},
	}

	// Field promotion: can access embedded fields directly!
	fmt.Printf("%s %s has a %d HP %s engine\n", myCar.Brand, myCar.Model, myCar.Horsepower, myCar.Type)

	// Exercise
	// 1: Value vs. Pointer Receivers
	fmt.Printf("\n ==== Struct Exercise ====\n")
	rect := Rectangle{Width: 5.0, Height: 3.0}

	fmt.Printf("Initial Area: %.2f\n", rect.Area())

	rect.Scale(2.0)

	fmt.Printf("Scaled Area: %.2f\n", rect.Area())

	// 2: Struct Embedding & Method Promotion
	admin := Admin{
		User: User{
			Username: "Tonny",
			Email:    "tonny@tonny.org",
		},
		Level: 1,
	}

	fmt.Println(admin.Details())

	// 3: Encapsulation & Validation Pattern
	wallet, err := NewWallet(100.0)
	if err != nil {
		fmt.Println("Error creating wallet:", err)
		return
	}

	err = wallet.Withdraw(150.0)
	if err != nil {
		fmt.Println("Withdrawal error:", err)
	}

	err = wallet.Withdraw(40.0)
	if err == nil {
		fmt.Printf("Remaining Balance: %.2f\n", wallet.Balance)
	}
}

// Exercise
// 1: Value vs. Pointer Receivers
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r *Rectangle) Scale(factor float64) {
	r.Height *= factor
	r.Width *= factor
}

// 2: Struct Embedding & Method Promotion
type User struct {
	Username string
	Email    string
}

func (u User) Details() string {
	return fmt.Sprintf("Username: [%s], Email: [%s]", u.Username, u.Email)
}

type Admin struct {
	User
	Level int
}

// 3: Encapsulation & Validation Pattern
type Wallet struct {
	Balance float64
}

// / Note: Go errors usually lowercase, no ending punctuation
func NewWallet(initial float64) (*Wallet, error) {
	if initial < 0 {
		return nil, fmt.Errorf("initial balance cannot be negative")
	}
	return &Wallet{Balance: initial}, nil
}

func (w *Wallet) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be greater than zero")
	}
	if amount > w.Balance {
		return fmt.Errorf("insufficient funds")
	}

	w.Balance -= amount
	return nil
}

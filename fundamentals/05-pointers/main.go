package main

import (
	"fmt"
)

func main() {
	fmt.Println("====== Pointers all the way we go yeah.... ======")

	x := 42

	// & gives memory address
	var ptr *int = &x

	fmt.Println("Value of x:", x)
	fmt.Println("Memory Address of x:", &x)
	fmt.Println("Value of ptr:", ptr)

	// * derefferences the pointer to get or change the value at that address
	fmt.Println("Value at pointer (*ptr):", *ptr)

	*ptr = 100
	fmt.Println("New value of x:", x)

	num := 10
	fmt.Println("Before:", num)

	// We pass the memory address using &
	doubleValue(&num)

	fmt.Println("After:", num)

	// Struct Pointers & The Arrow Syntax (-> / . shorthand)
	user := &User{Name: "Alice", Email: "alice@old.com"}

	fmt.Println("Before:", user.Email)

	updateEmail(user, "alice@new.com")

	fmt.Println("After:", user.Email) // alice@new.com

	// The new() Built-in Function
	// new(int) allocates memory for an int, sets it to 0, and returns *int
	p := new(int)

	fmt.Println("Address:", p) // Some memory address
	fmt.Println("Value:", *p)  // 0 (zero value of int)

	*p = 42
	fmt.Println("Updated Value:", *p) // 42

	fmt.Printf("\n === GO Pointer Exercise ==== \n")
	x, y := 10, 20
	fmt.Printf("Before: x = %d, y = %d\n", x, y)

	swap(&x, &y)

	fmt.Printf("After:  x = %d, y = %d\n", x, y)

	// Banking Scenario
	account := &BankAccount{Owner: "Alex", Balance: 100.0}

	Deposit(account, 50.0)
	fmt.Printf("Balance after deposit: %.2f\n", account.Balance)

	err := Withdraw(account, 200.0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	Withdraw(account, 30.0)
	fmt.Printf("Final Balance: %.2f\n", account.Balance)

	// Double slice values
	nums := []int{1, 2, 3, 4}

	DoubleSliceValues(&nums)

	fmt.Println(nums)
}

// Pointers in Functions
func doubleValue(n *int) {
	*n = *n * 2
}

// Struct pointer & Arror syntax
type User struct {
	Name  string
	Email string
}

func updateEmail(u *User, newEmail string) {
	u.Email = newEmail
}

// === GO Pointer Exercise ====
// In-Place Swap
func swap(a, b *int) {
	var temp int = *a
	*a = *b
	*b = temp
}

// State Mutator Function
type BankAccount struct {
	Owner   string
	Balance float64
}

func Deposit(acc *BankAccount, amount float64) {
	acc.Balance = acc.Balance + amount
}

func Withdraw(acc *BankAccount, amount float64) error {
	if amount > acc.Balance {
		return fmt.Errorf("Insufficient funds")
	}
	acc.Balance = acc.Balance - amount
	return nil
}

// Array Double via Pointer
func DoubleSliceValues(s *[]int) {
	for i := range *s {
		(*s)[i] = (*s)[i] * 2
	}
	*s = append(*s, 99)
}

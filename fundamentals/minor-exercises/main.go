package main

import "fmt"

func main() {
	fmt.Println("===== Lets start to practice some Go =====")
	user1 := Account{}
	user1.Balance = 500

	if err := user1.Deposit(-10); err != nil {
		fmt.Println("Deposit failed:", err)
	}

	if err := user1.Withdraw(-20); err != nil {
		fmt.Println("Withdrawal failed:", err)
	}

	book1 := Book{}
	book1.Title = "Atomic Habits"
	book1.Author = "James Clear"
	book1.isAvailable = true

	err := book1.Borrow()
	if err != nil {
		fmt.Println("Acquisition Failed:", err)
	}

	err = book1.Return()
	if err != nil {
		fmt.Println()
		fmt.Println("Return Failed: ", err)
	}

	book2 := Book{}
	book2.Title = "The Art of War"
	book2.Author = "Sun Tzu"
	book2.isAvailable = true

	if err := book2.Borrow(); err != nil {
		fmt.Println("Acquisition Failed:", err)
	}

	if err := book2.Return(); err != nil {
		fmt.Println()
		fmt.Println("Return Failed: ", err)
	}

	fmt.Printf("\n ==== Shopping Cart =====\n")
	cart := Cart{}

	phone := Product{Name: "Phone", Price: 300, Quantity: 2}
	laptop := Product{Name: "Laptop", Price: 500, Quantity: 5}
	keyboard := Product{Name: "Keyboard", Price: -100, Quantity: 1}

	if err := cart.AddProduct(phone); err != nil {
		fmt.Println("Error:", err)
	}
	if err := cart.AddProduct(laptop); err != nil {
		fmt.Println("Error:", err)
	}
	if err := cart.AddProduct(keyboard); err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\n--- Cart Before Removal ---")
	cart.ListProducts()

	fmt.Println("\n--- Removing Phone ---")
	if err := cart.RemoveProduct(phone.Name); err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\n--- Cart After Removal ---")
	cart.ListProducts()

	fmt.Printf("\n ===== Zoo ====== \n")
	dog := Dog{}
	cat := Cat{}
	lion := Lion{}

	MakeThemSpeak(dog)
	MakeThemSpeak(&cat)
	MakeThemSpeak(&lion)

	animals := []Animal{Dog{}, &Cat{}, &Lion{}}

	for _, animal := range animals {
		MakeThemSpeak(animal)
	}

	fmt.Printf("====== Payment System ======\n")
	paymentOptions := []PaymentGateway{Stripe{}, PayPal{}, MoMo{}}
	for _, paymentOption := range paymentOptions {
		err := Checkout(paymentOption)
		if err != nil {
			break
		}
	}
}

// Bank account
type Account struct {
	Balance float64
}

func (a *Account) Deposit(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("deposit amount must be greater than 0")
	}
	a.Balance += amount
	fmt.Printf("\nDeposit successful. New Balance: %.2f", a.Balance)
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be greater than 0")
	}
	if amount > a.Balance {
		return fmt.Errorf("insufficient funds")
	}
	a.Balance -= amount
	fmt.Printf("\nWithdraw successful, new balance: %.2f", a.Balance)
	return nil
}

// Library Managment
type Book struct {
	Title       string
	Author      string
	isAvailable bool
}

func (b *Book) Borrow() error {
	if b.isAvailable == false {
		return fmt.Errorf("the book is currently unavailable")
	}
	if b.isAvailable == true {
		fmt.Printf("You have borrowed the book %s", b.Title)
		b.isAvailable = false
	}
	return nil
}

func (b *Book) Return() error {
	if b.isAvailable == true {
		return fmt.Errorf("you haven't borrowed this book yet.")
	}
	b.isAvailable = true
	fmt.Printf("\n%s has been returned\n", b.Title)
	return nil
}

// Shopping Cart
type Product struct {
	Name     string
	Price    float64
	Quantity int
}

type Cart struct {
	Items []Product
}

func (c *Cart) AddProduct(product Product) error {
	if product.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}

	if product.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	c.Items = append(c.Items, product)
	return nil
}

func (c *Cart) RemoveProduct(name string) error {
	for item := 0; item < len(c.Items); item++ {
		if name == c.Items[item].Name {
			c.Items = append(c.Items[:item], c.Items[item+1:]...)
			return nil
		}
	}
	return fmt.Errorf("product %s not found in cart", name)
}

func (c *Cart) CalculateTotal() float64 {
	var total float64
	for _, p := range c.Items {
		total += p.Price * float64(p.Quantity)
	}
	return total
}

func (c *Cart) ListProducts() {
	if len(c.Items) == 0 {
		fmt.Println("The Cart is empty.")
		return
	}

	fmt.Printf("\n| %-10s | %-8s | %-8s | %-8s |\n", "Name", "Quantity", "Price", "Total")
	fmt.Println("-----------------------------------------------")
	for _, product := range c.Items {
		fmt.Printf("| %-10s | %-8d | %-8.2f | %-8.2f |\n", product.Name, product.Quantity, product.Price, product.Price*float64(product.Quantity))
	}
	fmt.Printf("| %-32s | %-8.2f |\n", "Total Item Price", c.CalculateTotal())
}

// Zoo
type Animal interface {
	Speak() string
}

type Dog struct{}
type Cat struct{}
type Lion struct{}

func (d Dog) Speak() string {
	return "Dog says \"woof woof\""
}

func (c *Cat) Speak() string {
	return "Cat says \"meow meow\""
}

func (l *Lion) Speak() string {
	return "Lion says \"roar roar\""
}

func MakeThemSpeak(a Animal) error {
	fmt.Println(a.Speak())
	return nil
}

// Payment System
type PaymentGateway interface {
	Pay(amount float64) error
}

type Stripe struct{}
type PayPal struct{}
type MoMo struct{}

func (s Stripe) Pay(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Amount must be greater than zero")
	}
	if amount >= 1000 {
		return fmt.Errorf("The amount $%.2f exceeds daily transaction limit", amount)
	}
	fmt.Printf("Paid %.2f using Stripe\n", amount)
	return nil
}

func (p PayPal) Pay(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Amount must be greater than zero")
	}
	if amount >= 1200 {
		return fmt.Errorf("The amount $%.2f exceeds daily transaction limit", amount)
	}
	fmt.Printf("Paid %.2f using Paypal\n", amount)
	return nil
}

func (m MoMo) Pay(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Amount must be greater than zero")
	}
	if amount >= 3000 {
		return fmt.Errorf("The amount $%.2f exceeds daily transaction limit", amount)
	}
	fmt.Printf("Paid %.2f using Mobile Money\n", amount)
	return nil
}

func Checkout(g PaymentGateway) error {
	err := g.Pay(100)

	if err != nil {
		fmt.Println()
		fmt.Println("Error:", err)
	}

	return err
}

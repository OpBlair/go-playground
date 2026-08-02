package main

import (
	"fmt"
)

// PaymentGateway defines the interface for processing payments.
// As a student, think of this as a rulebook or contract: any type that wants
// to be a PaymentGateway must have a Pay method with this exact signature.
type PaymentGateway interface {
	Pay(amount float64) error
}

// Stripe represents the Stripe payment processor.
type Stripe struct{}

// Pay implements the PaymentGateway interface for Stripe.
func (s Stripe) Pay(amount float64) error {
	fmt.Printf("Paid $%.2f using Stripe\n", amount)
	return nil
}

// PayPal represents the PayPal payment processor.
type PayPal struct{}

// Pay implements the PaymentGateway interface for PayPal.
// Because it has a Pay method matching the interface, it automatically satisfies
// PaymentGateway without needing an explicit "implements" keyword.
func (p PayPal) Pay(amount float64) error {
	fmt.Printf("Paid $%.2f using PayPal\n", amount)
	return nil
}

// Checkout accepts any payment gateway, demonstrating polymorphism.
// Instead of writing separate checkout functions for Stripe and PayPal,
// this single function works with any type that implements PaymentGateway.
func Checkout(gateway PaymentGateway, amount float64) {
	// We call Pay without needing to know if it's Stripe or PayPal under the hood
	err := gateway.Pay(amount)
	if err != nil {
		fmt.Println("Payment failed:", err)
	}
}

// Speaker defines an interface for things that can speak.
type Speaker interface {
	Speak() string
}

// Dog represents a dog entity.
type Dog struct{}

// Speak satisfies the Speaker interface for Dog.
func (d Dog) Speak() string {
	return "Dog says \"woof woof\""
}

// Robot represents a robot entity.
type Robot struct{}

// Speak satisfies the Speaker interface for Robot.
func (r Robot) Speak() string {
	return "Robot says \"beep boop\""
}

// MakeThemSpeak takes any Speaker and prints what it says.
func MakeThemSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	fmt.Println("====== Interfaces & Error Handling Tutorial ======")

	// 1. Using the payment gateways polymorphically
	var stripeGateway PaymentGateway = Stripe{}
	var paypalGateway PaymentGateway = PayPal{}

	Checkout(stripeGateway, 49.99)
	Checkout(paypalGateway, 25.50)

	fmt.Println()

	// 2. Using the speaker example
	dog := Dog{}
	robot := Robot{}

	MakeThemSpeak(dog)
	MakeThemSpeak(robot)
}

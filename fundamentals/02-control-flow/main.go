package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Control Flow")

	if score := 85; score >= 80 {
		fmt.Println("Greate job!")
	} else {
		fmt.Println("Try harder next time")
	}

	day := "Tuesday"

	switch day {
	case "Monday":
		fmt.Println("Start of the work week.")
	case "Tuesday":
		fmt.Println("Tackle control flow!")
	default:
		fmt.Println("Another day, another bug fixed.")
	}

	fmt.Println("\nCounting up:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("Count: %d\n", i)
	}

	FizzBuzz()
	LeapYearValidator()
	Calculator()
	PasswordAttempt()
}

func FizzBuzz() {
	fmt.Println("======= FizzBuzz here man ======")
	for i := 0; i < 31; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Printf("Number %2d : FizzBuzz\n", i)
		} else if i%3 == 0 {
			fmt.Printf("Number %2d : Fizz\n", i)
		} else if i%5 == 0 {
			fmt.Printf("Number %2d : Buzz\n", i)
		} else {
			fmt.Printf("Number %2d : %d \n", i, i)
		}
	}
}

func LeapYearValidator() {
	year := 2020

	fmt.Printf("\n ==== Leap Year Validator ==== \n")
	if year%400 == 0 || year%4 == 0 && year%100 != 0 { // && takes precedence over ||
		fmt.Printf("%d is a leap year \n", year)
	} else {
		fmt.Printf("%d is not a leap year \n", year)
	}
}

func Calculator() {
	var choice int
	var num1, num2 float64

	for {
		fmt.Println("=== Calculator ===")
		fmt.Println("1. +")
		fmt.Println("2. -")
		fmt.Println("3. *")
		fmt.Println("4. /")
		fmt.Println("0. Exit")

		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		if choice == 0 {
			fmt.Println("Exiting... bye")
			return
		} else {
			switch choice {
			case 1:
				fmt.Print("Enter first number: ")
				fmt.Scanln(&num1)

				fmt.Print("Enter second number: ")
				fmt.Scanln(&num2)

				fmt.Printf("%.2f + %.2f = %.2f \n", num1, num2, num1+num2)
			case 2:
				fmt.Print("Enter first number: ")
				fmt.Scanln(&num1)

				fmt.Print("Enter second number: ")
				fmt.Scanln(&num2)

				fmt.Printf("%.2f - %.2f = %.2f \n", num1, num2, num1-num2)
			case 3:
				fmt.Print("Enter first number: ")
				fmt.Scanln(&num1)

				fmt.Print("Enter second number: ")
				fmt.Scanln(&num2)

				fmt.Printf("%.2f x %.2f = %.2f \n", num1, num2, num1*num2)
			case 4:
				fmt.Print("Enter first number: ")
				fmt.Scanln(&num1)

				fmt.Print("Enter second number: ")
				fmt.Scanln(&num2)

				if num2 == 0 {
					fmt.Println("Division by zero is not allowed")
				} else {
					fmt.Printf("%.2f / %.2f = %.2f \n", num1, num2, num1/num2)
				}
			default:
				fmt.Println("Invalide input, Try again")
			}
		}
	}
}

func PasswordAttempt() {
	password := "password123"
	var inputPwd string
	var failedAttempt int
	maxAttempts := 3

	for {
		fmt.Print("Enter Password: ")
		fmt.Scanln(&inputPwd)

		if inputPwd == password {
			fmt.Println("Logging in...")
			return
		}

		fmt.Println("Incorrect password. Try again")
		failedAttempt++

		if failedAttempt >= maxAttempts {
			fmt.Println("Too many attempts. Try again in ...")

			for i := 10; i > 0; i-- {
				fmt.Printf("Try again in %d \r", i)
				time.Sleep(1 * time.Second)
			}

			fmt.Println("You can now enter password.")
			failedAttempt = 0
		}
	}
}

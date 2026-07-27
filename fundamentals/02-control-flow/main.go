package main

import "fmt"

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
}

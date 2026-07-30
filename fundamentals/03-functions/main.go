package main

import (
	"fmt"
)

func main() {
	fmt.Println("==== Functions in Go ====")
	greet()

	// Function with return value
	sum := add(10, 20)
	fmt.Printf("Add Result: %d\n", sum)

	// Multiple returns handling
	result, err := divide(10.0, 2.0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Divide Result: %.2f\n", result)
	}

	// Named returns handling
	s, p := calculate(5, 4)
	fmt.Printf("Calculate -> Sum: %d, Product: %d\n", s, p)

	// Variadic function handling
	total := sumAll(1, 2, 3, 4, 5)
	fmt.Printf("Sum All Result: %d\n", total)

	// functions as values
	subFunc := func(minuend, subtrahend int) int {
		return minuend - subtrahend
	}

	// Function as value handling
	difference := subFunc(10, 7)
	fmt.Printf("Subtraction Result: %d\n", difference)

	/* ==== Exercise ===== */
	fmt.Printf("\n ===== Go Exercises ===== \n")
	data := []int{4, 2, 9, 1, 5, 6}
	min, max, avg := AnalyzeStats(data)
	fmt.Printf("Min: %d, Max: %d, Avg: %.2f\n", min, max, avg)

	LogInfo("WARNING:", "Disk space low", "CPU temperature high", "Service restarting")

	// Anonymous Function
	fmt.Printf("\nDealing with Anonymous Functions\n")
	counter := MakeCounter()
	fmt.Println(counter())
	fmt.Println(counter())
	fmt.Println(counter())
}

func greet() {
	fmt.Println("Hello, there, am a function and you just have summoned me.")
}

func add(num1 int, num2 int) int {
	return num1 + num2
}

// Multiple return values
func divide(dividend float64, divisor float64) (float64, error) {
	if divisor == 0 {
		return 0, fmt.Errorf("Cannot divide by Zero")
	}
	return dividend / divisor, nil
}

// Named return values
func calculate(x, y int) (sum int, product int) {
	sum = x + y
	product = x * y
	return
}

// Variadic Functions
func sumAll(numbers ...int) int {
	total := 0

	for _, num := range numbers {
		total += num
	}
	return total
}

/*========= Exercise ========*/
func AnalyzeStats(nums []int) (int, int, float64) {
	min := nums[0]
	max := nums[0]
	var sum int
	for i := 0; i <= len(nums)-1; i++ {
		if nums[i] < min {
			min = nums[i]
		}
		if nums[i] > max {
			max = nums[i]
		}
		sum += nums[i]
	}
	avg := float64(sum) / float64(len(nums))
	return min, max, avg
}

func LogInfo(prefix string, messages ...string) {
	fmt.Printf("\nPrinting Logs\n")
	for _, msg := range messages {
		fmt.Printf("%s : %s \n", prefix, msg)
	}
}

func MakeCounter() func() int {
	count := 0

	return func() int {
		count++
		return count
	}
}

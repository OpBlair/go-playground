package main

import (
	"fmt"
)

func main() {
	fmt.Println("==== Collections Tutorial ====")

	// Arrays
	var fruits [5]string = [5]string{"mango", "apple", "strawberry", "orange", "guava"}

	for i := range fruits {
		fmt.Printf("%s is a fruit\n", fruits[i])
	}

	// Slice
	nums := []int{10, 20, 30, 40}
	nums = append(nums, 40)
	fmt.Printf("Numbers = %d\n", nums)

	// creating slice with 'make' keyword
	scoreSheet := make([]int, 3, 5)
	scoreSheet = append(scoreSheet, 10)
	fmt.Println("Scoresheet ", scoreSheet)

	letters := []string{"a", "b", "c"}
	fmt.Println("Length of slice 'letters' {length is current size}:", len(letters))
	fmt.Println("Capacity of slice 'letters', {capacity is max current size}:", cap(letters))

	letters = append(letters, "d")
	fmt.Println("New length after appending", len(letters))
	fmt.Println("Capacity doubles in size when new item is appended:", cap(letters))

	// Maps (key-value pairs)
	ages := map[string]int{
		"Alice": 24,
		"Bob":   22,
	}

	ages["Charlie"] = 26

	//Commma-Ok Idiom for checking a value
	age, exists := ages["David"]
	if exists {
		fmt.Printf("David's age: %d\n", age)
	} else {
		fmt.Printf("David not found in map")
	}

	// Exerice Solutions
	fmt.Printf("\n ===== Go Collections Exercise Solution =====")
	words := []string{"apple", "banana", "apple", "orange", "banana", "apple"}
	counts := CountWords(words)

	fmt.Printf("\n %v", counts)

	// Exercise No. 2 slice & creating subslices
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	subSlice := numbers[3:7]
	fmt.Printf("\nSubSlice: %v, Len: %d, Cap: %d\n", subSlice, len(subSlice), cap(subSlice))
	subSlice[0] = 99
	fmt.Println("Original numbers:", numbers)

	fmt.Printf("\nGrouping By Value\n")
	class := []Student{
		{"Alice", "A"},
		{"Bob", "B"},
		{"Charlie", "A"},
		{"Diana", "C"},
		{"Eve", "B"},
	}

	grouped := GroupByGrade(class)

	for grade, names := range grouped {
		fmt.Printf("%s: %v\n", grade, names)
	}
}

// Exercise For Collections
func CountWords(words []string) map[string]int {
	count := make(map[string]int)

	for _, words := range words {
		count[words]++
	}
	return count
}

type Student struct {
	Name  string
	Grade string
}

func GroupByGrade(students []Student) map[string][]string {
	grouped := make(map[string][]string)

	for _, student := range students {
		// Append the student's name to the slice located at the key 'student.Grade'
		grouped[student.Grade] = append(grouped[student.Grade], student.Name)
	}

	return grouped
}

package main

import (
	"fmt"
)

// Student Record Manager
type Student struct {
	Id     int
	Name   string
	Course string
	Marks  []int
}

var studentRecord []Student

func studentRecordManager() {
	for {
		var choice int
		fmt.Println("\n--- MENU ---")
		fmt.Println("1. Add a student")
		fmt.Println("2. Remove student")
		fmt.Println("3. Search student")
		fmt.Println("4. List all students")
		fmt.Println("5. Calculate average marks")
		fmt.Println("6. Exit")

		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var id int
			var name, course string

			fmt.Print("Enter ID: ")
			fmt.Scanln(&id)
			fmt.Print("Enter Name (No spaces): ")
			fmt.Scanln(&name)
			fmt.Print("Enter Course: ")
			fmt.Scanln(&course)

			startingMarks := []int{80, 85, 90}

			addStudent(id, name, course, startingMarks)

		case 2:
			var id int
			fmt.Print("Enter ID to remove: ")
			fmt.Scanln(&id)
			removeStudent(id)

		case 3:
			var id int
			fmt.Print("Enter ID to search: ")
			fmt.Scanln(&id)
			searchStudent(id)

		case 4:
			listStudents()

		case 5:
			var id int
			fmt.Print("Enter ID to calculate average: ")
			fmt.Scanln(&id)
			calculateAverage(id)

		case 6:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Please enter a valid choice.")
		}
	}
}

// Add student to the slice
func addStudent(id int, name string, course string, marks []int) {
	newStudent := Student{Id: id, Name: name, Course: course, Marks: marks}
	studentRecord = append(studentRecord, newStudent)
	fmt.Printf("Successfully added %s\n", name)
}

// List students
func listStudents() {
	fmt.Println("\n--- Current Student Lists ---")
	if len(studentRecord) == 0 {
		fmt.Println("No records found.")
		return
	}
	for _, student := range studentRecord {
		fmt.Printf("ID: %d | Name: %s | Course: %s\n", student.Id, student.Name, student.Course)
	}
}

// Search Students
func searchStudent(id int) {
	for _, student := range studentRecord {
		if student.Id == id {
			fmt.Printf("\n[Found] ID: %d, Name: %s, Course: %s\n", student.Id, student.Name, student.Course)
			return
		}
	}
	fmt.Printf("\n[Error] Student with ID %d not found.\n", id)
}

// CALCULATE AVERAGE MARKS
func calculateAverage(id int) {
	for _, student := range studentRecord {
		if student.Id == id {
			if len(student.Marks) == 0 {
				fmt.Printf("%s has no marks recorded.\n", student.Name)
				return
			}

			total := 0
			for _, mark := range student.Marks {
				total = total + mark
			}

			average := float64(total) / float64(len(student.Marks))
			fmt.Printf("\nAverage marks for %s: %.2f\n", student.Name, average)
			return
		}
	}
	fmt.Printf("\n[Error] Cannot calculate average. Student %d not found.\n", id)
}

// REMOVE STUDENT
func removeStudent(id int) {
	for index, student := range studentRecord {
		if student.Id == id {
			studentRecord = append(studentRecord[:index], studentRecord[index+1:]...)
			fmt.Printf("\nSuccessfully removed student ID %d.\n", id)
			return
		}
	}
	fmt.Printf("\n[Error] Cannot remove. Student %d not found.\n", id)
}

func main() {
	addStudent(1, "John_Doe", "BSCS", []int{70, 80, 85, 89})

	fmt.Println("=== Student Record Manager ===")
	studentRecordManager()
}

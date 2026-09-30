package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var tasks []Task

type Task struct {
	Title       string
	isCompleted bool
	Priority    string
}

const dataFile = "tasks.json"

func main() {

	LoadTasks()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("===== TO-DO App ======")
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Mark Tasks Complete")
		fmt.Println("4. Remove Tasks")
		fmt.Println("5. Edit Task")
		fmt.Println("0. Exit")

		fmt.Print("Enter your choice: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
			return
		}

		input := scanner.Text()

		choice, err := strconv.Atoi(input)

		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("Adding tasks")
			AddTask(scanner)
			SaveTasks()
		case 2:
			ListTasks()
			SaveTasks()
		case 3:
			fmt.Println("marking tasks complete")
			MarkComplete(scanner)
			SaveTasks()
		case 4:
			fmt.Println("removing tasks")
			RemoveTask(scanner)
			SaveTasks()
		case 5:
			fmt.Println("Editing task")
			EditTask(scanner)
			SaveTasks()
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter choice again")
		}
	}
}

func SaveTasks() {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		fmt.Println("Error Saving tasks:", err)
		return
	}

	err = os.WriteFile(dataFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func LoadTasks() {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println("Error parsing Tasks json:", err)
	}
}

func AddTask(scanner *bufio.Scanner) {
	fmt.Print("Enter name of task: ")

	scanner.Scan()
	name := scanner.Text()

	fmt.Println()
	fmt.Println("Set Task Priority: ")
	fmt.Println("1. High")
	fmt.Println("2. Medium")
	fmt.Println("3. Low")

	fmt.Print("Enter Priority:")
	scanner.Scan()
	priority, err := strconv.Atoi(scanner.Text())

	if err != nil {
		fmt.Println("Invalid priority. Please enter a number.")
		return
	}

	if priority < 1 || priority > 3 {
		fmt.Println("Invalid Priority Number. Try Again")
		return
	}

	var priorityName string
	switch priority {
	case 1:
		priorityName = "High"
	case 2:
		priorityName = "Medium"
	case 3:
		priorityName = "Low"
	default:
		priorityName = "Low"
	}

	task := Task{Title: name, isCompleted: false, Priority: priorityName}

	tasks = append(tasks, task)
}

func ListTasks() {
	if len(tasks) == 0 {
		fmt.Println("Task List is Empty.")
		return
	}

	fmt.Printf("| %-3s | %-20s | %-10s | %-9s |\n", "No.", "Task Name", "Priority", "Complete")
	for i, task := range tasks {
		status := ""

		if !task.isCompleted {
			status = "\u2717"
		}
		if task.isCompleted {
			status = "\u2713" // unicode for check mark character
		}
		fmt.Printf("| %-3d | %-20s | %-10s | %-9s |\n", i+1, task.Title, task.Priority, status)
	}
}

func MarkComplete(scanner *bufio.Scanner) {
	if len(tasks) == 0 {
		fmt.Println("Task List is Empty. Cannot toggle from an empty List.")
		return
	}

	fmt.Print("Enter the task number:")
	scanner.Scan()

	number, err := strconv.Atoi(scanner.Text())

	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return
	}

	if number < 1 || number > len(tasks) {
		fmt.Println("Invalid Task number. Check Task number and try again")
		return
	}

	tasks[number-1].isCompleted = true
}

func RemoveTask(scanner *bufio.Scanner) {
	if len(tasks) == 0 {
		fmt.Println("Task List is Empty. Cannot remove from an empty List.")
		return
	}

	fmt.Print("Enter the task number:")
	scanner.Scan()

	number, err := strconv.Atoi(scanner.Text())

	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return
	}

	if number < 1 || number > len(tasks) {
		fmt.Println("Invalid Task number. Check Task number and try again")
		return
	}

	fmt.Printf("You are removing the task \"%s\" from task list.\n", tasks[number-1].Title)

	fmt.Printf("Enter: \n1)Y to continue \n2)N to cancel \n")

	fmt.Println("Your choice: ")
	scanner.Scan()

	choice := scanner.Text()

	if choice == "y" || choice == "Y" {
		taskName := tasks[number-1].Title
		tasks = append(tasks[:number-1], tasks[number:]...)
		fmt.Printf("Task \"%s\" has been removed from task list successfully.\n", taskName)
	}

	if choice == "n" || choice == "N" {
		fmt.Println("Removal Cancelled")
		return
	}
}

func EditTask(scanner *bufio.Scanner) {
	if len(tasks) == 0 {
		fmt.Println("Task List is Empty. Cannot edit from an empty List.")
		return
	}

	fmt.Print("Enter the task number: ")
	scanner.Scan()

	number, err := strconv.Atoi(scanner.Text())

	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return
	}

	if number < 1 || number > len(tasks) {
		fmt.Println("Invalid Task number. Check Task number and try again")
		return
	}

	// Identify the task being edited
	taskIndex := number - 1
	currentTask := tasks[taskIndex]

	fmt.Printf("Editing Task #%d: %s\n", number, currentTask.Title)

	// Edit Title
	fmt.Print("Enter new task name (leave blank to keep current): ")
	scanner.Scan()
	newTitle := scanner.Text()

	if newTitle != "" {
		tasks[taskIndex].Title = newTitle
	}

	// Edit Priority
	fmt.Printf("Current Priority: %s\n", currentTask.Priority)
	fmt.Println("Set New Task Priority (leave blank to keep current): ")
	fmt.Println("1. High")
	fmt.Println("2. Medium")
	fmt.Println("3. Low")
	fmt.Print("Enter New Priority: ")

	scanner.Scan()
	priorityInput := scanner.Text()

	if priorityInput != "" {
		priority, err := strconv.Atoi(priorityInput)
		if err != nil || priority < 1 || priority > 3 {
			fmt.Println("Invalid priority input. Keeping current priority.")
		} else {
			switch priority {
			case 1:
				tasks[taskIndex].Priority = "High"
			case 2:
				tasks[taskIndex].Priority = "Medium"
			case 3:
				tasks[taskIndex].Priority = "Low"
			}
		}
	}

	fmt.Println("Task updated successfully!")
}

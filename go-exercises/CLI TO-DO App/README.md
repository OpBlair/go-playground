# Go CLI To-Do App

A simple command-line TODO application built with Go.

This project is a learning project focused on practicing Go fundamentals by building a practical CLI application.

## Features

* **Add Tasks:** Create new tasks with custom names and set priority levels (**High**, **Medium**, **Low**).
* **List Tasks:** View all tasks in a clean, organized table with their current status and priority.
* **Mark Complete:** Mark tasks as completed once they are finished.
* **Remove Tasks:** Delete tasks with a confirmation prompt.
* **Edit Tasks:** Edit an existing task.
* **Validate User Input:** Handle invalid input with error checking.
* **Interactive CLI:** Navigate the application through a simple command-line menu.

## Example

```text
===== TO-DO App ======
1. Add Task
2. List Tasks
3. Mark Tasks Complete
4. Remove Tasks
5. Edit Task
0. Exit

Enter your choice: 1

Enter name of task: Learn Go

Set Task Priority:
1. High
2. Medium
3. Low

Enter Priority: 1
```

Tasks are displayed with their priority and completion status:

| Number | Task Name | Priority | Complete |
| ------ | --------- | -------- | -------- |
| 1      | Learn Go  | High     | ✓        |

## ⚙️ Requirements

* Go 1.20 or newer
* [Download and Install Go](https://go.dev/dl/)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/go-playground
cd go-playground/go-exercises/CLI\TO-DO\App
```

### 2. Initialize the Go module

If a Go module has not already been initialized:

```bash
go mod init todo-app
```

### 3. Run the application

```bash
go run .
```

## Usage

When you launch the application, you will see the following menu:

```text
===== TO-DO App ======
1. Add Task
2. List Tasks
3. Mark Tasks Complete
4. Remove Tasks
0. Exit
Enter your choice:
```

* **`1` — Add Task:** Enter a task title and select a priority level (`1. High`, `2. Medium`, `3. Low`).
* **`2` — List Tasks:** Displays all tasks with their task number, title, priority, and completion status.
* **`3` — Mark Tasks Complete:** Enter a task number to mark it as completed.
* **`4` — Remove Tasks:** Enter a task number and confirm with `Y` to delete it or `N` to cancel.
* **`0` — Exit:** Exits the application.

## Go Concepts Practiced

This project is being developed incrementally while learning Go.

Some of the core concepts practiced in the codebase include:

* Variables and data types
* Functions and passing values
* Structs and struct fields
* Slices
* `append`
* Slice indexing and manipulation
* `for` loops and control flow
* `switch` statements
* User input with `bufio.Scanner`
* Error handling
* String-to-integer conversion with `strconv.Atoi`
* String formatting with `fmt.Printf`
* Unicode characters

## Project Structure

```text
.
├── main.go
├── go.mod
└── README.md
```

### Core Code Breakdown

* **`main()`** — Handles the main application loop, input routing, and menu display.
* **`Task`** — Struct that defines the task data model with `Title`, `isCompleted`, and `Priority`.
* **`AddTask()`** — Handles user input for creating and adding new tasks.
* **`ListTasks()`** — Displays tasks in a formatted table.
* **`MarkComplete()`** — Updates a task's status to completed.
* **`RemoveTask()`** — Removes a task from the slice after user confirmation.

## Future Improvements

* Toggle tasks between complete and incomplete
* Improve CLI formatting and styling
* Edit existing tasks
* Persist tasks to a file (JSON/CSV)
* Load tasks automatically when the application starts
* Add automated unit tests
* Improve input validation
* Refactor the application into multiple packages

## Purpose

The goal of this project is to learn Go by building a practical application and gradually applying new concepts as they are learned.

# Student Record Manager (Go CLI + Companion UI)

A lightweight student record management system built in **Go (Golang)** featuring a command-line interface (CLI) for in-memory data operations and an optional modern web UI companion.

---

## Current Features (Go CLI In-Memory Engine)

The core backend runs as a Go console program using slices and structs (`[]studentRecord`) to manage records entirely in memory:

1. **Add Student (`addStudent`):** Input a unique Student ID, Name, Course, and a default slice of marks.
2. **Remove Student:** Delete a student record from memory by their ID.
3. **Search Student:** Quickly look up a student's profile and recorded marks by ID.
4. **List All Students:** Display all active student records currently stored in memory.
5. **Calculate Average:** Compute and display the exact grade average from a student's marks slice.
6. **Exit Console:** Cleanly terminate the CLI session.

---

## Project Structure

```text
├── main.go              # Go backend logic & CLI menu loop
├── script.js            # Frontend interactivity for the UI mockup
├── index.html           # Simple Tailwind CSS companion mockup UI
└── README.md            # Project documentation & future roadmap

```

---

## Getting Started (Running the CLI)

1. Ensure you have [Go installed](https://golang.org/) on your system.
2. Clone this repository and:
```bash
cd student-manager
```
3. Run the CLI program:
```bash
go run main.go

```

4. Follow the interactive console prompts to test choices `1` through `6`.

---

## Future Roadmap: Scaling to a Full School Management System (SMS)

We are actively evolving this project from a single-user CLI tool into a full-scale, multi-role web application.

### 1. Architectural Evolution

* **Shift from In-Memory to Persistent DB:** Migrating from Go slices to **SQLite** (for local development) and PostgreSQL/MySQL for production using a robust relational schema.
* **Multi-Page Web Architecture:** Moving from a single mockup dashboard to a dedicated multi-page system with server-side or frontend routing.

### 2. Multi-Role Authentication & Access Control

Using a clean domain-based mock login (`admin@westridge.ac`, `lecturer@westridge.ac`, `student@westridge.ac`) to route users into role-specific portals:

* **Admin Portal:** User provisioning (creating lecturer/student accounts), managing course offerings, departments, and system-wide reporting.
* **Lecturer Portal:** Viewing assigned class rosters, entering and modifying assignment/exam grades, and creating online quizzes.
* **Student Portal:** Personal dashboard to view enrolled courses, take active quizzes set by lecturers, view instant auto-graded feedback, and check academic transcripts.

---

## 📄 License

This project is open-source and available under the [MIT License](https://opensource.org/license/mit).

# Go Fundamentals

This directory contains my hands-on code experiments, notes, and exercises covering core Go programming concepts.

## Topics Covered

| Topic | Folder | Status | Key Concepts |
| :--- | :--- | :--- | :--- |
| **01. Basics & Types** | [`01-basics/`](./01-basics/) | Completed | Variables, Constants, Zero Values, Type Conversion, Formatting Verbs |
| **02. Control Flow** | `02-control-flow/` | Upcoming | `if/else`, `switch`, `for` loops |
| **03. Functions** | `03-functions/` | Upcoming | Multiple returns, variadic functions, closures |
| **04. Collections** | `04-collections/` | Upcoming | Arrays, Slices, Maps |
| **05. Pointers** | `05-pointers/` | Upcoming | Memory addresses, dereferencing, pass-by-reference |
| **06. Structs & Methods** | `06-structs/` | Upcoming | Custom types, receiver functions |
| **07. Interfaces** | `07-interfaces/` | Upcoming | Behavior implementation, polymorphism |
| **08. Error Handling** | `08-errors/` | Upcoming | The `error` interface, custom errors, panic/recover |

---

## How to Run Code

To run any subfolder independently, navigate into it and use `go run`:

```bash
cd fundamentals/01-basics
go run main.go
```
## Folder Structure
```
go-playground/
└── fundamentals/
    ├── README.md
    ├── 01-variables-and-types/
    │   ├── main.go
    │   └── notes.md
    ├── 02-control-flow/
    │   ├── main.go
    │   └── notes.md
    ├── 03-functions/
    │   ├── main.go
    │   └── calc.go
    ├── 04-arrays-slices-maps/
    │   └── main.go
    ├── 05-structs-and-interfaces/
    │   └── main.go
    └── 06-pointers-and-errors/
        └── main.go
```

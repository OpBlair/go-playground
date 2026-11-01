# Student Record Manager (Go CLI Evolving to Web SMS)

A lightweight student record management system built in **Go (Golang)**, originally started as a command-line interface (CLI) for in-memory data operations, and currently evolving into a secure, role-based web application backed by PostgreSQL and JWT authentication.

---

## Evolution & Progress

### Phase 1: Go CLI In-Memory Engine (Initial Foundation)
The core backend started as a Go console program using slices and structs (`[]studentRecord`) to manage records entirely in memory:
1. **Add Student (`addStudent`):** Input a unique Student ID, Name, Course, and a default slice of marks.
2. **Remove Student:** Delete a student record from memory by their ID.
3. **Search Student:** Quickly look up a student's profile and recorded marks by ID.
4. **List All Students:** Display all active student records currently stored in memory.
5. **Calculate Average:** Compute and display the exact grade average from a student's marks slice.
6. **Exit Console:** Cleanly terminate the CLI session.

### Phase 2: PostgreSQL Persistence & Database Constraints (Current Stage)
Migrated from temporary memory slices to a robust relational database schema hosted on an Alpine Linux VM:
- **Relational Integrity:** Implemented `UUID` identifiers (`uuid-ossp`), strict `user_role` ENUM types (`admin`, `lecturer`, `student`), and composite foreign key constraints to strictly guarantee that courses and quizzes can only be assigned to valid accounts holding the `'lecturer'` role.
- **Broadcast Messaging:** Normalized notification structures separating core message content (`notifications`) from individual user read statuses (`notification_recipients`) with partial indexing for rapid unread lookups.

### Phase 3: JWT Authentication & Role-Based Access Control (RBAC)
- Shifted from static session mocks to stateless, signed JSON Web Tokens (JWT).
- Implemented a custom Go `authMiddleware` that verifies token signatures, inspects roles, and secures sensitive endpoints (such as admin-managed user creation).

---

## Project Structure

```text
├── dbschema.sql         # Master PostgreSQL schema, ENUMs, and constraints
├── main.go              # Go backend logic, HTTP router, JWT middleware, and DB handlers
├── script.js            # Frontend interactivity, JWT storage, and fetch API logic
├── login.html           # Authentication portal
├── admin.html           # Admin control center
├── lecturer.html        # Lecturer workspace
├── student.html         # Student dashboard
├── .env                 # Environment variables (DB connection, JWT secret)
└── README.md            # Project documentation & evolution roadmap
```
---

## 🛠️ Project Tech Stack

- **Backend / API:** Go (Golang) using standard `net/http` and `github.com/golang-jwt/jwt/v5` for authentication (running locally on Fedora host).
- **Database:** PostgreSQL (with `uuid-ossp` extension, ENUM types, and strict relational foreign key constraints) running inside an **Alpine Linux VM**..
- **Security & Auth:** `bcrypt` for secure password hashing and stateless JWT-based Role-Based Access Control (RBAC).
- **Frontend:** HTML5, Tailwind CSS (via CDN), and Vanilla JavaScript (`fetch` API interacting with Go JSON endpoints).
- **Infrastructure:** Database running in a guest VM and backend on host on fedora.

---

## ⚙️ PostgreSQL Installation & Database Setup

If you are setting up the database from scratch on a fresh environment, follow these steps:

### 1. Install PostgreSQL
* **On Alpine Linux:**
  ```bash
  apk update
  apk add postgresql postgresql-contrib

  ```

* **On Ubuntu / Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib

```

### 2. Initialize & Start the Database Service

* **On Alpine Linux:**
```bash
/etc/init.d/postgresql setup
rc-service postgresql start

```

* **On Ubuntu / Debian:**
```bash
sudo systemctl start postgresql
sudo systemctl enable postgresql

```

### 3. Create the Database and User

Switch to the postgres user and create your database environment:

```bash
sudo -i -u postgres
psql

```

Inside the PostgreSQL interactive terminal, run:

```sql
CREATE DATABASE <database_name>;
CREATE USER <database_username> WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE <database_name> TO <database_username>;
\q
exit

```

### 4. Apply the Master Schema

Run your schema file to set up tables, extensions, ENUMs, and foreign key constraints:

```bash
psql -U <database_username> -d <database_name> -f dbschema.sql

```
---

## Getting Started (Running the CLI)

1. Ensure you have [Go installed](https://golang.org/) on your system.
2. Clone this repository and:
```bash
cd student-manager/v1
```
3. Run the CLI program:
```bash
go run main.go

```

4. Follow the interactive console prompts to test choices `1` through `6`.

---

## Getting Started (Running the Web Version)

1. Ensure you have [Go installed](https://golang.org/) on your system.
2. Clone this repository and:
```bash
cd student-manager
```
3. Run the CLI program:
```bash
go run main.go
```
4. Navigate to the directed url on terminal
---

## Future Roadmap: Feature Implementations

We are actively building out full-stack API integration to replace static UI mocks across our multi-role portals:

### 1. Admin Control Center (`admin.html` & Go Handlers)
- **User Provisioning API (`POST /api/admin/create-user`):** Securing user creation behind the `authMiddleware` (admin role check) and hashing passwords via `bcrypt`.
- **Program & Course Management APIs:** Wiring up database CRUD endpoints for creating programs and academic courses with strict lecturer validation constraints.
- **Dashboard Metrics API (`GET /api/admin/stats`):** Dynamically fetching total users, active courses, programs, and student counts for summary cards.

### 2. Lecturer & Student Portals
- **Lecturer Workspace:** Publishing weekly course modules and setting target-week quizzes.
- **Student Dashboard:** Enrolling in courses, viewing grades, and checking unread inbox notifications.

---

## 📄 License

This project is open-source and available under the [MIT License](https://opensource.org/license/mit).

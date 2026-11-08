package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// ==========================================
// TYPES & STRUCTS
// ==========================================

type Config struct {
	Port   string
	DBConn string
}

// Custom JWT Claims struct
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Login Request Payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create Program Payload
type CreateProgramRequest struct {
	ProgramCode     string `json:"program_code"`
	ProgramName     string `json:"program_name"`
	ProgramDuration int    `json:"duration"`
}

// CREATE USER REQUEST PAYLOAD
type CreateUserRequest struct {
	Password  string `json:"user_password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"user_email"`
	Role      string `json:"user_role"`
	ProgramID *int   `json:"program_id"`
}

// Struct for program response items
type ProgramResponse struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

// CREATE COURSE REQUEST PAYLOAD
type CreateCourseRequest struct {
	Code         string  `json:"code"`
	Title        string  `json:"title"`
	ProgramID    int     `json:"program_id"`
	IsElective   bool    `json:"is_elective"`
	AcademicYear int     `json:"academic_year"`
	Semester     int     `json:"semester"`
	LecturerID   *string `json:"lecturer_id"`
}

// STRUCT FOR LECTURER RESPONSE ITEMS
type LecturerResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// STRUCT FOR DASHBOARD RESPONSE ITEMS
type DashBoardStats struct {
	TotalUsers     int `json:"total_users"`
	TotalAdmins    int `json:"total_admins"`
	TotalLecturers int `json:"total_lecturers"`
	TotalStudents  int `json:"total_students"`
	TotalPrograms  int `json:"total_programs"`
	ActiveCourses  int `json:"active_courses"`
}

// STRUCT FOR USER RECORDS
type UserResponse struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	AssignedProgram string `json:"assigned_program,omitempty"`
}

// STRUCT FOR COURSE RECORDS
type CourseResponse struct {
	ID           int    `json:"id"`
	Code         string `json:"code"`
	Title        string `json:"title"`
	IsElective   bool   `json:"is_elective"`
	AcademicYear int    `json:"academic_year"`
	Semester     int    `json:"semester"`
	Capacity     int    `json:"capacity"`
}

// Struct for creating a notification request
type CreateNotificationRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// Struct for notification response items
type NotificationResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// Struct for password change response items
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// Context keys for injecting user details safely
type contextKey string

const userContextKey contextKey = "userClaims"

// ==========================================
// CONFIGURATION & DATABASE
// ==========================================

func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found.")
	}

	port := os.Getenv("DB_PORT")
	dbConn := os.Getenv("DB_CONN")

	return Config{
		Port:   port,
		DBConn: dbConn,
	}
}

func connectDB(ctx context.Context) *pgxpool.Pool {
	dbConn := os.Getenv("DB_CONN")
	pool, err := pgxpool.New(ctx, dbConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database pool: %v\n", err)
		os.Exit(1)
	}
	return pool
}

// ==========================================
// MIDDLEWARE & SECURITY
// ==========================================

// Auth Middleware to protect restricted routes
func authMiddleware(next http.HandlerFunc, requiredRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims := &Claims{}

		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Check role restriction if one was specified
		if requiredRole != "" && claims.Role != requiredRole {
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		}

		// Inject claims into request context so the handler can read the UserID
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Helper to generate JWT Token
func generateJWT(userID, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "fallback_secret"
	}

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ==========================================
// API HANDLERS
// ==========================================

// LOGIN API
func handleLogin(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Query database to get ID, role, and hash
		var userID string
		var storedHash string
		var role string

		query := "SELECT id, role, password_hash FROM users WHERE email = $1"
		err := conn.QueryRow(r.Context(), query, req.Email).Scan(&userID, &role, &storedHash)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate JWT Token
		tokenString, err := generateJWT(userID, role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Success! Send back the JWT token and role
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"success": true,
			"token":   tokenString,
			"role":    role,
		}
		json.NewEncoder(w).Encode(response)
	}
}

// CREATE PROGRAM API
func handleCreateProgram(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check method type
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode incoming json into struct
		var req CreateProgramRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Extract user claims from request context (injected by authMiddleware)
		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check the role just to be extra sure
		if claims.Role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}

		// Insert Program into database
		query := "INSERT INTO programs(code, name, duration) VALUES($1, $2, $3)"

		// Execute the insert query
		_, err := conn.Exec(r.Context(), query, req.ProgramCode, req.ProgramName, req.ProgramDuration)
		if err != nil {
			// Check if it's a duplicate code error (unique constraint violation)
			if strings.Contains(err.Error(), "duplicate key") {
				http.Error(w, "Program code already exists", http.StatusConflict)
				return
			}
			http.Error(w, "Failed to create program", http.StatusInternalServerError)
			return
		}

		// Send back success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Program created successfully!",
		})
	}
}

// HANDLE GET PROGRAMS (Public or Protected - protected for admins/users)
func handleGetPrograms(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := "SELECT id, code, name, duration FROM programs ORDER BY name ASC"
		rows, err := conn.Query(r.Context(), query)
		if err != nil {
			http.Error(w, "Failed to fetch programs", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var programs []ProgramResponse
		for rows.Next() {
			var p ProgramResponse
			if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.Duration); err != nil {
				http.Error(w, "Error parsing programs", http.StatusInternalServerError)
				return
			}
			programs = append(programs, p)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(programs)
	}
}

// CREATE USER API
func handleCreateUser(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check method type
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode incoming json into struct
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Extract user claims from request context (injected by authMiddleware)
		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check the role just to be extra sure (Admins only)
		if claims.Role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}

		// Hash the password securely with bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}

		// Insert User into database
		query := "INSERT INTO users(password_hash, first_name, last_name, email, role, program_id) VALUES($1, $2, $3, $4, $5, $6)"

		// Execute the insert query (using hashedPassword instead of raw req.Password)
		_, err = conn.Exec(r.Context(), query, string(hashedPassword), req.FirstName, req.LastName, req.Email, req.Role, req.ProgramID)
		if err != nil {
			// Check if it's a duplicate email error (unique constraint violation)
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "users_email_key") {
				http.Error(w, "Email already exists", http.StatusConflict)
				return
			}
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Send back success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "User created successfully!",
		})
	}
}

// CREATE COURSE API
func handleCreateCourse(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateCourseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Extract user claims from request context (injected by authMiddleware)
		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Restrict to admins only
		if claims.Role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}

		// Insert course into database
		// Note: lecturer_role is hardcoded to 'lecturer' satisfy your composite foreign key constraint
		query := `
            INSERT INTO courses (code, title, program_id, is_elective, academic_year, semester, lecturer_id, lecturer_role) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, CASE WHEN $7::uuid IS NOT NULL THEN 'lecturer'::user_role ELSE NULL END)
        `

		_, err := conn.Exec(r.Context(),
			query,
			req.Code,
			req.Title,
			req.ProgramID,
			req.IsElective,
			req.AcademicYear,
			req.Semester,
			req.LecturerID,
		)

		if err != nil {
			log.Printf("Database error creating course: %v\n", err)
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "courses_code_key") {
				http.Error(w, "Course code already exists", http.StatusConflict)
				return
			}
			http.Error(w, "Failed to create course", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Course created successfully!",
		})
	}
}

// GET LECTURERS API
func handleGetLecturers(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := "SELECT id, first_name, last_name FROM users WHERE role = 'lecturer' ORDER BY last_name ASC"
		rows, err := conn.Query(r.Context(), query)
		if err != nil {
			http.Error(w, "Failed to fetch lecturers", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var lecturers []LecturerResponse
		for rows.Next() {
			var l LecturerResponse
			if err := rows.Scan(&l.ID, &l.FirstName, &l.LastName); err != nil {
				http.Error(w, "Error parsing lecturers", http.StatusInternalServerError)
				return
			}
			lecturers = append(lecturers, l)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lecturers)
	}
}

// GET STATS API
func handleGetStats(conn *pgxpool.Pool) http.HandlerFunc { // Outer Function -> The Factory
	return func(w http.ResponseWriter, r *http.Request) { // Inner Function -> The Delivery (argument 1 is used to send data back the user, argument 2 contains all data sent by user's browser or client)
		// Check incoming Request Method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()       // (The Timeline) grabs incoming http requst. It's a request's lifeline.
		var stats DashBoardStats // (The bucket); holds data pulled from the database.

		// Total Users in the system
		err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers) // conn is the active db connection,
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Total Active Admins in the system
		err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&stats.TotalAdmins)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Total Active Lecturers in the system
		err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'lecturer'").Scan(&stats.TotalLecturers)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Total Active Students in the system
		err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'student'").Scan(&stats.TotalStudents)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Total Active Programs in the system
		err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM programs").Scan(&stats.TotalPrograms)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Total Active courses in the system
		err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM courses").Scan(&stats.ActiveCourses)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Return JSON Response
		w.Header().Set("Content-Type", "application/json") // Modifies http response headers (Tell receiver browser the kind of data coming)
		w.WriteHeader(http.StatusOK)                       // Stamp of Approval; tells client browser request was processed successfully
		json.NewEncoder(w).Encode(stats)                   // The Delivery System; w opens a direct pipeline for client browser and its streamed directly to the client by Encode after serialization
	}
}

// HANDLE GET USERS
func handleGetUsers(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Explicitly cast u.role::text to prevent ENUM scanning errors in pgx
		query := `
            SELECT u.id::text, u.first_name, u.last_name, u.email, u.role::text, p.name 
            FROM users u 
            LEFT JOIN programs p ON u.program_id = p.id 
            ORDER BY u.last_name ASC
        `
		rows, err := conn.Query(r.Context(), query)
		if err != nil {
			log.Printf("DB Query Error (Users): %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []UserResponse
		for rows.Next() {
			var u UserResponse
			var progName *string
			if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &progName); err != nil {
				log.Printf("Scan Error (Users): %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Smart role-based fallback when program is not applicable
			if progName != nil {
				u.AssignedProgram = *progName
			} else {
				switch u.Role {
				case "admin":
					u.AssignedProgram = "System-wide"
				case "lecturer":
					u.AssignedProgram = "Faculty / General"
				default:
					u.AssignedProgram = "Unassigned"
				}
			}
			users = append(users, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// HANDLE GET COURSES
func handleGetCourses(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := `
            SELECT id, code, title, is_elective, academic_year, semester 
            FROM courses 
            ORDER BY code ASC
        `
		rows, err := conn.Query(r.Context(), query)
		if err != nil {
			log.Printf("DB Query Error (Courses): %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var courses []CourseResponse
		for rows.Next() {
			var c CourseResponse
			if err := rows.Scan(&c.ID, &c.Code, &c.Title, &c.IsElective, &c.AcademicYear, &c.Semester); err != nil {
				log.Printf("Scan Error (Courses): %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			courses = append(courses, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(courses)
	}
}

// ADMIN BROADCAST NOTIFICATION
func handleCreateNotification(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok || claims.Role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}

		var req CreateNotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		// Begin transaction because we insert into notifications and fan out to recipients
		tx, err := conn.Begin(ctx)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(ctx)

		// Insert notification and get its new ID
		var notificationID int
		err = tx.QueryRow(ctx, "INSERT INTO notifications (title, message) VALUES ($1, $2) RETURNING id", req.Title, req.Message).Scan(&notificationID)
		if err != nil {
			http.Error(w, "Failed to create notification", http.StatusInternalServerError)
			return
		}

		// Fan out to all users in the system via the junction table
		_, err = tx.Exec(ctx, `
            INSERT INTO notification_recipients (notification_id, user_id)
            SELECT $1, id FROM users
        `, notificationID)
		if err != nil {
			http.Error(w, "Failed to dispatch notification recipients", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			http.Error(w, "Failed to commit notification transaction", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Notification broadcasted successfully!",
		})
	}
}

// GET USER NOTIFICATIONS
func handleGetNotifications(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		query := `
            SELECT n.id, n.title, n.message, nr.is_read, n.created_at
            FROM notifications n
            JOIN notification_recipients nr ON n.id = nr.notification_id
            WHERE nr.user_id = $1
            ORDER BY n.created_at DESC
        `
		rows, err := conn.Query(r.Context(), query, claims.UserID)
		if err != nil {
			http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notifications []NotificationResponse
		for rows.Next() {
			var n NotificationResponse
			if err := rows.Scan(&n.ID, &n.Title, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
				http.Error(w, "Error parsing notifications", http.StatusInternalServerError)
				return
			}
			notifications = append(notifications, n)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifications)
	}
}

// MARK NOTIFICATION AS READ
func handleMarkNotificationRead(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			NotificationID int `json:"notification_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		query := `
            UPDATE notification_recipients 
            SET is_read = TRUE, read_at = CURRENT_TIMESTAMP 
            WHERE notification_id = $1 AND user_id = $2
        `
		_, err := conn.Exec(r.Context(), query, req.NotificationID, claims.UserID)
		if err != nil {
			http.Error(w, "Failed to update notification status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}
}

// Change password API
func handleChangePassword(conn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req ChangePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		// Fetch current stored password hash for the user
		var storedHash string
		err := conn.QueryRow(ctx, "SELECT password_hash FROM users WHERE id = $1", claims.UserID).Scan(&storedHash)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Verify old password
		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.OldPassword))
		if err != nil {
			http.Error(w, "Incorrect old password", http.StatusUnauthorized)
			return
		}

		// Hash the new password securely
		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to process new password", http.StatusInternalServerError)
			return
		}

		// Update database
		_, err = conn.Exec(ctx, "UPDATE users SET password_hash = $1 WHERE id = $2", string(newHashedPassword), claims.UserID)
		if err != nil {
			http.Error(w, "Failed to update password in database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Password changed successfully!",
		})
	}
}

// ==========================================
// SERVER & ROUTING
// ==========================================

func startServer(conn *pgxpool.Pool) {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy", "message": "Westridge API is running!"}`))
	})

	// Public Route
	http.HandleFunc("/api/login", handleLogin(conn))

	// Protected Admin Route for Creating Programs
	http.HandleFunc("/api/programs", authMiddleware(handleCreateProgram(conn), "admin"))

	// User Creation Route
	http.HandleFunc("/api/users", authMiddleware(handleCreateUser(conn), "admin"))

	// Protected route to fetch programs dynamically
	http.HandleFunc("/api/programs-list", authMiddleware(handleGetPrograms(conn), ""))

	// Fetch Lecturers Route
	http.HandleFunc("/api/lecturers-list", authMiddleware(handleGetLecturers(conn), ""))

	// Create Course Route
	http.HandleFunc("/api/courses", authMiddleware(handleCreateCourse(conn), "admin"))

	// Stats Route
	http.HandleFunc("/api/stats", authMiddleware(handleGetStats(conn), "admin"))

	// users list Route
	http.HandleFunc("/api/users-list", authMiddleware(handleGetUsers(conn), ""))

	// courses list Route
	http.HandleFunc("/api/courses-list", authMiddleware(handleGetCourses(conn), ""))

	// Notification Routes
	http.HandleFunc("/api/notifications", authMiddleware(handleGetNotifications(conn), ""))
	http.HandleFunc("/api/notifications/broadcast", authMiddleware(handleCreateNotification(conn), "admin"))
	http.HandleFunc("/api/notifications/read", authMiddleware(handleMarkNotificationRead(conn), ""))

	// Password Change Route
	http.HandleFunc("/api/auth/change-password", authMiddleware(handleChangePassword(conn), ""))

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	serverAddr := ":" + serverPort
	log.Printf("Starting Westridge web server on http://localhost%s\n", serverAddr)

	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

// ==========================================
// MAIN ENTRY POINT
// ==========================================

func main() {
	loadConfig()
	log.Println("Server initializing configuration...")

	ctx := context.Background()
	dbPool := connectDB(ctx)
	defer dbPool.Close()

	fmt.Println("Successfully connected to Westridge DB!")
	startServer(dbPool)
}

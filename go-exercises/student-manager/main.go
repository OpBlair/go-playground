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
	"github.com/jackc/pgx/v5"
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

func connectDB(ctx context.Context) *pgx.Conn {
	dbConn := os.Getenv("DB_CONN")
	conn, err := pgx.Connect(ctx, dbConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
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
func handleLogin(conn *pgx.Conn) http.HandlerFunc {
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
func handleCreateProgram(conn *pgx.Conn) http.HandlerFunc {
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
func handleGetPrograms(conn *pgx.Conn) http.HandlerFunc {
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
func handleCreateUser(conn *pgx.Conn) http.HandlerFunc {
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
func handleCreateCourse(conn *pgx.Conn) http.HandlerFunc {
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
func handleGetLecturers(conn *pgx.Conn) http.HandlerFunc {
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

// ==========================================
// SERVER & ROUTING
// ==========================================

func startServer(conn *pgx.Conn) {
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
	conn := connectDB(ctx)
	defer conn.Close(ctx)

	fmt.Println("Successfully connected to Westridge DB!")
	startServer(conn)
}

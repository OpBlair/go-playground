package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Port   string
	DBConn string
}

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

// connectDB handles the connection logic and returns the connection or an error.
func connectDB(ctx context.Context) *pgx.Conn {
	dbConn := os.Getenv("DB_CONN")

	conn, err := pgx.Connect(ctx, dbConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func startServer(conn *pgx.Conn) { // <-- Accept conn here
	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	// Sample API health check endpoint
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy", "message": "Westridge API is running!"}`))
	})

	// Register login route with the database connection!
	http.HandleFunc("/api/login", handleLogin(conn))

	// Use PORT from .env or default to 8080
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

// Logic to handle Logging in
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// the handler is wrapped so it has access to the pgx connection
func handleLogin(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check Request Type
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode the incoming JSON from script.js
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Query the database using pgx
		var storedHash string
		var role string

		query := "SELECT role, password_hash FROM users WHERE email = $1"
		err := conn.QueryRow(r.Context(), query, req.Email).Scan(&role, &storedHash)
		if err != nil {
			// If user isn't found or query fails, return unauthorized
			log.Printf("LOGIN ERROR: User not found for email '%s': %v\n", req.Email, err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Verify password using bcrypt against the seeded hash
		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
		if err != nil {
			log.Printf("LOGIN ERROR: Password mismatch for email '%s'\n", req.Email)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Success! Send back the role to frontend script.js
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "role": "` + role + `"}`))
	}
}

func main() {
	loadConfig()

	log.Println("Server initializing configuration...")

	ctx := context.Background()

	conn := connectDB(ctx)
	defer conn.Close(ctx)

	fmt.Println("Successfully connected to Westridge DB!")

	// Start the web server
	startServer(conn)
}

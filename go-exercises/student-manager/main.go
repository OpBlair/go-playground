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

func startServer(conn *pgx.Conn) {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy", "message": "Westridge API is running!"}`))
	})

	// Public Route
	http.HandleFunc("/api/login", handleLogin(conn))

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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

// Context keys for injecting user details safely
type contextKey string

const userContextKey contextKey = "userClaims"

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

func main() {
	loadConfig()
	log.Println("Server initializing configuration...")

	ctx := context.Background()
	conn := connectDB(ctx)
	defer conn.Close(ctx)

	fmt.Println("Successfully connected to Westridge DB!")
	startServer(conn)
}

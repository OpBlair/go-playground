package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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

// setup and start the server
func startServer() {
	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	// API health check endpoint
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy", "message": "Westridge API is running!"}`))
	})

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	serverAddr := ":" + serverPort
	log.Printf("Starting Westridge web server on http://localhost%s\n", serverAddr)

	// Start the HTTP server
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
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
	startServer()
}

package main

import (
	"context"
	"fmt"
	"log"
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

func main() {
	loadConfig()

	log.Println("Server initializing configuration...")

	ctx := context.Background()

	conn := connectDB(ctx)
	defer conn.Close(ctx)

	fmt.Println("Successfully connected to Westridge DB!")
}

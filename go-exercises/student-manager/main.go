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

func main() {
	cfg := loadConfig()

	log.Println("Server initializing configuration...")

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DBConn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	fmt.Println("Successfully connected to Westridge DB!")
}

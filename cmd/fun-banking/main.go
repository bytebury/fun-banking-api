package main

import (
	"funbanking/internal/infrastructure/persistence"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("🟡 Unable to load .env configuration")
	}
	persistence.SetUpConnection()
}

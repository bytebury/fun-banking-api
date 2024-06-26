package main

import (
	"funbanking/internal/domain/announcements"
	"funbanking/internal/domain/banking"
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/persistence"
	"funbanking/internal/interfaces/api"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("🟡 Unable to load .env configuration")
	}

	persistence.SetUpConnection()
	users.RunMigrations()
	banking.RunMigrations()
	announcements.RunMigrations()

	api.Run()
}

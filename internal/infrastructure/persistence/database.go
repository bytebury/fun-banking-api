package persistence

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type databaseConfig struct {
	DSN string
}

func newConnection(config *databaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DSN))

	if err != nil {
		return nil, err
	}

	return db, err
}

func SetUpConnection() {
	db, err := newConnection(&databaseConfig{DSN: os.Getenv("DATABASE_URL")})

	if err != nil {
		log.Fatalln("🔴 Unable to connect to the database", err)
	}

	DB = db
	log.Println("🟢 Successfully connected to the database")
}

func RunMigrations() {
	// TODO
}

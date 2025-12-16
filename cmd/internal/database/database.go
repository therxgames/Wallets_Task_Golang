package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"wallet/cmd/internal/models"
)

var DB *gorm.DB

func Init() {
	connectDB()
	migrateTables()
	fmt.Println("Connected to wallet DB")
}

func InitTestDB() {
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	createTestDatabase(dbName, host, user, password, port)
	connectToTestDB(dbName, host, user, password, port)
	resetTables()
	migrateTables()

	fmt.Println("Connected to wallet_test DB")
}

func connectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
}

func createTestDatabase(dbName, host, user, password, port string) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		host, user, password, port,
	)

	sqlDB, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))

	if err != nil && !contains(err.Error(), "already exists") {
		log.Fatalf("failed to create test db: %v", err)
	}
}

func connectToTestDB(dbName, host, user, password, port string) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}
}

func resetTables() {
	DB.Exec("DROP TABLE IF EXISTS wallet_operations CASCADE")
	DB.Exec("DROP TABLE IF EXISTS wallets CASCADE")
}

func migrateTables() {
	if err := DB.AutoMigrate(&models.Wallet{}, &models.WalletOperation{}); err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && contains(s[1:], substr)))
}

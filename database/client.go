package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnDB() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// get data from .env file
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")

	// connect database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_pass, db_name)
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return DB
}

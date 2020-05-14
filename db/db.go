package db

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/lib/pq"
)

var (
	host = "localhost"
	port = "5432"
	user = "friend"
	password = "p123123"
	dbname = "friend"
)

var db *sql.DB

func GetInstance() *sql.DB {
	return db
}

func Init() {
	var err error
	db, err = sql.Open("postgres", connString())
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting database...")
}

func connString() string {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl != "" {
		return dbUrl
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
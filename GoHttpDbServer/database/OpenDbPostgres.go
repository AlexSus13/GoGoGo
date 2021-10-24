package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Config struct {
	User     string
	Host     string
	Password string
	Port     string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s host=%s password=%s port=%s dbname=%s sslmode=%s",
		cfg.User, cfg.Host, cfg.Password, cfg.Port, cfg.DBName, cfg.SSLMode))
	if err != nil {
		log.Fatal(err) //error handling
	}
	//defer db.Close()

	errr := db.Ping()
	if errr != nil {
		log.Fatal(errr) //error handling
	}

	return db, nil
}

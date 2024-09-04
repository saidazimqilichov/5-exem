package postgres

import (
	"database/sql"
	"fmt"
	"os"
	_"github.com/lib/pq"
	_"github.com/joho/godotenv/autoload"

)

func ConnectPostgres() (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

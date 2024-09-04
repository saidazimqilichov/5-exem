package postgress

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "user_postgres"
	user     = "postgres" 
	password = "7777"
	port     = 5432
	dbname   = "hotel"      
)

func ConnectPostgres() (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", host, user, password, port, dbname)
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

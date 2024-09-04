package connpostgresql

import (
	"database/sql"
	"fmt"
	"hotel-service/config"

	_ "github.com/lib/pq"
)

func ConnectToDB(cfg config.Config) (*sql.DB, error) {
	psqlString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)
	connDb, err := sql.Open("postgres", psqlString)
	if err != nil {
		return nil, err
	}

	err = connDb.Ping()
	if err!= nil {
        return nil, err
    }

	return connDb, nil
}
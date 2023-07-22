package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func CreateDbConnection() (*sqlx.DB, error) {
	dbHost := "localhost"
	dbPort := 5432
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := "restapi-books"

	dbUrl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}

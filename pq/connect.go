package pq

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host="localhost"
	user="postgres"
	password="1234"
	dbname="chatbox"
	port=5432
)

func ConnectDB() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbname, port)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
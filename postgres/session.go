package postgres

import "database/sql"

import _ "github.com/lib/pq"

func CreateDBConnection() *sql.DB {
	connStr := "postgresql://postgres:postgres@localhost/test?sslmode=disable" // TODO: Env

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

var DB = CreateDBConnection()

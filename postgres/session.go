package postgres

import "database/sql"

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

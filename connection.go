package main

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

var Connector *sql.DB

func Connect() error {
	var e error
	Connector, e = sql.Open(
		"postgres",
		`host=10.14.206.27
		 port=5432
		 user=student
		 password=1234
		 dbname=my_database
		 sslmode=disable`,
	)
	if e != nil {
		return errors.New("не удалось подключиться к базе данных")
	}

	return nil
}

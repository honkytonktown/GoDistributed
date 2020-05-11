package datamanager

import (
	"database/sql"

	//blank import
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:####@localhost:5432/GoDistributed?sslmode=disable")

	if err != nil {
		panic(err.Error())
	}
}

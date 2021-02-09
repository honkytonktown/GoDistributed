package datamanager

import (
	"database/sql"
	"fmt"

	//blank import
	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	hostname     = "localhost"
	hostport     = 5432
	username     = "postgres"
	password     = ""
	databasename = "GoDistributed"
)

func init() {
	var err error
	pg_con_string := fmt.Sprintf("port=%d host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostport, hostname, username, password, databasename)

	db, err = sql.Open("postgres", pg_con_string)
	if err != nil {
		panic(err)
	}
	//defer db.Close()W

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("You are Successfully connected!")
}

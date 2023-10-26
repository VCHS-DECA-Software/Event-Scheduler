package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	schemaLocation := flag.String("sql", "./schema.sql", "specify sql to execute on creation")
	dbLocation := flag.String("db", "./local.db", "specify the location of the database file")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schemaStatement, err := os.ReadFile(*schemaLocation)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(schemaStatement))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("executed given sql file.")
}

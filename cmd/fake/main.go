package main

import (
	"database/sql"
	"event-scheduler/lib/input"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := sql.Open("sqlite3", "local.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	f, err := os.Open("./conference.csv")
	if err != nil {
		log.Fatal(err)
	}

	df := dataframe.ReadCSV(f)
	err = input.LoadInputs(db, df)
	if err != nil {
		log.Fatal(err)
	}
}

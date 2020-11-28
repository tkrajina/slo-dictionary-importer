package main

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tkrajina/slo-dictionary-importer/importer"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	dbFile := "dict.sqlite3"
	_ = os.Remove(dbFile)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		tn := "thesaurus"
		panicIfErr(createTable(db, tn))
		panicIfErr(importer.ImportThesaurus(db, tn))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tn := "collocations"
		panicIfErr(createTable(db, tn))
		panicIfErr(importer.ImportCollocations(db, tn))
	}()

	wg.Wait()
}

func createTable(db *sql.DB, tableName string) error {
	if _, err := db.Exec("CREATE TABLE " + tableName + " (id INTEGER primary key not null, word TEXT, info_json TEXT, search_str TEXT)"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE UNIQUE INDEX " + tableName + "_word_unique_idx ON " + tableName + "(word)"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE INDEX " + tableName + "_search_str_idx ON " + tableName + "(search_str)"); err != nil {
		return err
	}
	return nil
}

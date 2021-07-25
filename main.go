package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tkrajina/slo-dictionary-importer/importer"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

var commands = map[string]func(){
	"app-db":      BuildAppDb,
	"kindle-dict": BuildKindleDict,
}

func main() {
	if len(os.Args) < 2 {
		var cmds []string
		for cmd := range commands {
			cmds = append(cmds, cmd)
		}
		sort.Strings(cmds)
		fmt.Printf("Command not given, select one of %v\n", cmds)
		return
	}
	command := os.Args[1]
	if f, found := commands[command]; found {
		f()
	} else {
		fmt.Println("Command", command, "not found")
	}
}

func BuildAppDb() {
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

func BuildKindleDict() {
	var wg sync.WaitGroup

	// Load synonyms
	// Load collocations
	// Load slolex

	var (
		thesaurus    []importer.ThesaurusEntry
		collocations []importer.CollocationXMLEntry
		slolex       []importer.SlolexLexicalEntry
		err          error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		thesaurus, err = importer.LoadThesaurus()
		panicIfErr(err)
	}()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	collocations, err = importer.LoadCollocations()
	// 	panicIfErr(err)
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	slolex, err = importer.SlolexLoader()
	// 	panicIfErr(err)
	// }()

	wg.Wait()

	_ = thesaurus
	_ = collocations
	_ = slolex

	var dict importer.KindleDict
	for n, thesaurusEntry := range thesaurus {
		if n > 1000 {
			break
		}
		synonymsCore := []string{}
		for _, g := range thesaurusEntry.GroupsCore.Group {
			for _, g2 := range g.Candidate {
				synonymsCore = append(synonymsCore, g2.S.Text)
			}
		}
		synonymsNear := []string{}
		for _, g := range thesaurusEntry.GroupsNear.Group {
			for _, g2 := range g.Candidate {
				synonymsNear = append(synonymsNear, g2.S.Text)
			}
		}
		dict.Entries = append(dict.Entries, importer.KindleDictEntry{
			Word:        thesaurusEntry.Headword.Text,
			Description: "<p>" + strings.Join(synonymsCore, "; ") + "</p>" + "<p>" + strings.Join(synonymsNear, "; ") + "</p>",
		})
	}
	importer.ExportOPF(dict)
	fmt.Println("Now open kindledict/slo.opf in Kindle previewer and export the dictionary")
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

package main

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"

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
		thesaurus            []importer.ThesaurusEntry
		slolexByLema         = map[string][][]string{}
		collocationsByLema   = map[string][][]string{}
		frequencyOrderByLema = map[string]importer.WordFrequency{}
		err                  error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		thesaurus, err = importer.LoadThesaurus()
		panicIfErr(err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := importer.LoadCollocationsChan()
		panicIfErr(err)

		for col := range ch {
			panicIfErr(col.Err)
			collocationsByLema[col.Entry.Word()] = col.Entry.GetFrequentCollocations()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := importer.SlolexLoaderChan()

		for e := range ch {
			panicIfErr(e.Err)
			lema := e.Entry.Lema.FindLema()
			slolexByLema[lema] = e.Entry.FindFormRepresentations()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := importer.LoadFrequencyChan()

		for e := range ch {
			panicIfErr(e.Err)
			frequencyOrderByLema[e.Entry.Lemma] = e.Entry
		}

	}()

	wg.Wait()

	var dict importer.KindleDict
	for _, thesaurusEntry := range thesaurus {
		word := thesaurusEntry.Headword.Text

		for _, r := range word {
			if !unicode.IsLetter(r) {
				fmt.Println("Contains non-letters:", word)
				continue
			}
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
		entry := importer.KindleDictEntry{
			ID:          1 + len(dict.Entries),
			Word:        word,
			Description: "<br/><i>syn.</i> " + html.EscapeString(strings.Join(synonymsCore, "; ")) + "<br/>",
		}
		if len(synonymsNear) > 0 {
			entry.Description += "<i>syn.</i> " + html.EscapeString(strings.Join(synonymsNear, "; ")) + "<br/>"
		}
		inflections := map[string]bool{word: true}
		if slolexEntry, found := slolexByLema[entry.Word]; found {
			//fmt.Println("inflections:", slolexEntry)
			for _, g1 := range slolexEntry {
				for _, inflection := range g1 {
					if _, found := inflections[inflection]; !found {
						entry.Inflections = append(entry.Inflections, inflection)
					}
					inflections[inflection] = true
				}
			}
		} else {
			fmt.Println("No slolex for", entry.Word, "skipping")
			continue
		}
		if col, found := collocationsByLema[entry.Word]; found {
			for n, g := range col {
				if n >= 2 {
					continue
				}
				for m, collocation := range g {
					if m >= 2 {
						continue
					}
					entry.Description += "<i>e.g.</i> " + html.EscapeString(collocation) + "<br/>"
				}
			}
		}
		if freq, found := frequencyOrderByLema[word]; found {
			entry.Description += "<i>freq.</i> " + fmt.Sprintf("%.3f", freq.Frequency) + "%<br/>"
		}
		dict.Entries = append(dict.Entries, entry)
		if len(dict.Entries)%100 == 0 {
			fmt.Println("Building kindle dictionary entry #", len(dict.Entries))
		}
	}

	panicIfErr(importer.ExportOPF(dict))
	fmt.Println("Exported", len(dict.Entries), "entries")
	fmt.Println("Now open kindledict/slo-thesuarus.opf in Kindle previewer and export the dictionary")
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

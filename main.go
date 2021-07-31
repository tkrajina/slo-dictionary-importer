package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
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
	"app-db":              BuildAppDb,
	"kindle-dict":         BuildKindleDict,
	"rebuild-kindle-dict": ExportOPF,
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
thesurus_loop:
	for _, thesaurusEntry := range thesaurus {
		word := thesaurusEntry.Headword.Text

		for _, r := range word {
			if !unicode.IsLetter(r) || unicode.IsSpace(r) {
				fmt.Println("Contains non-letters:", word)
				continue thesurus_loop
			}
		}

		synonymsCore := []importer.KindleDictSynonym{}
		for _, g := range thesaurusEntry.GroupsCore.Group {
			for _, g2 := range g.Candidate {
				score, err := strconv.ParseFloat(g2.Score, 64)
				panicIfErr(err)
				synonymsCore = append(synonymsCore, importer.KindleDictSynonym{Word: g2.S.Text, Score: score})
			}
		}
		synonymsNear := []importer.KindleDictSynonym{}
		for _, g := range thesaurusEntry.GroupsNear.Group {
			for _, g2 := range g.Candidate {
				score, err := strconv.ParseFloat(g2.Score, 64)
				panicIfErr(err)
				synonymsNear = append(synonymsNear, importer.KindleDictSynonym{Word: g2.S.Text, Score: score})
			}
		}
		entry := importer.KindleDictEntry{
			ID:       1 + len(dict.Entries),
			Word:     word,
			Synonyms: [][]importer.KindleDictSynonym{synonymsCore, synonymsNear},
		}
		if len(synonymsCore) == 0 && len(synonymsNear) == 0 {
			continue thesurus_loop
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
			//continue
		}
		if col, found := collocationsByLema[entry.Word]; found {
			entry.Collocations = col
		}
		if freq, found := frequencyOrderByLema[word]; found {
			entry.Frequency = fmt.Sprintf("%.3f", freq.Frequency) + "%"
		}
		dict.Entries = append(dict.Entries, entry)
		if len(dict.Entries)%100 == 0 {
			fmt.Println("Building kindle dictionary entry #", len(dict.Entries))
		}
	}

	byts, err := json.Marshal(dict)
	panicIfErr(err)

	err = ioutil.WriteFile(kindleDictJson, byts, 0700)
	panicIfErr(err)

	ExportOPF()
}

const kindleDictJson = "kindledict/slo.json"

func ExportOPF() {
	max := math.MaxInt64
	if len(os.Args) >= 3 {
		m, err := strconv.ParseInt(os.Args[2], 10, 64)
		panicIfErr(err)
		max = int(m)
	}
	fmt.Println("max=", max)

	byts, err := ioutil.ReadFile(kindleDictJson)
	panicIfErr(err)

	var dict importer.KindleDict
	panicIfErr(json.Unmarshal(byts, &dict))

	if max < len(dict.Entries) {
		dict.Entries = dict.Entries[0:max]
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

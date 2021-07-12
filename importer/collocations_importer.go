package importer

import (
	"database/sql"
	"fmt"
	"sort"
)

func ImportCollocations(db *sql.DB, tableName string) error {
	collocations, err := LoadCollocations()
	if err != nil {
		return err
	}
	sumFrequencies := 0.0
	numFrequencies := 0
	for n, col := range collocations {
		var info [][][2]interface{}
		for n, gramrel := range col.Body.Sense.Gramrels.Gramrel {
			if n > 5 {
				continue
			}
			info = append(info, [][2]interface{}{})
			collocations := gramrel.Collocations.Collocation
			sort.Sort(collocations)
			if len(collocations) > 10 {
				collocations = collocations[0:10]
			}
			for _, collocation := range collocations {
				//fmt.Println("-", collocation.Frequency, collocation.Form)
				info[len(info)-1] = append(info[len(info)-1], [2]interface{}{collocation.Frequency, collocation.Form})
				sumFrequencies += collocation.Frequency
				numFrequencies++
			}
			//fmt.Println()
		}
		if err := insert(db, tableName, row{id: n, word: col.Header.LexicalUnit.Text, info: info}); err != nil {
			return err
		}
	}
	fmt.Println("Avg collocation frequency:", sumFrequencies/float64(numFrequencies))
	return nil
}

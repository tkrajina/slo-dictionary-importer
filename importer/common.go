package importer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type row struct {
	id   int
	word string
	info interface{}
}

func insert(db *sql.DB, tableName string, r row) error {
	ij, err := json.Marshal(r.info)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("insert into %s (id, word, info_json, search_str) values (?, ?, ?, ?)", tableName)
	vals := []interface{}{r.id, r.word, string(ij), strings.ToLower(r.word)}
	if r.id%100 == 0 {
		fmt.Printf("%d (%s): %#v\n", r.id, tableName, vals)
	}
	_, err = db.Exec(sql, vals...)
	return err
}

package importer

import (
	"database/sql"
)

func ImportThesaurus(db *sql.DB, tableName string) error {
	// Open a zip archive for reading.
	entries, err := LoadThesaurus()
	if err != nil {
		return err
	}

	for n, entry := range entries {
		var info [][][]Word
		group, err := importGroup("Core", entry.GroupsCore.Group)
		if err != nil {
			return err
		}
		info = append(info, group)
		group, err = importGroup("Near", entry.GroupsNear.Group)
		if err != nil {
			return err
		}
		info = append(info, group)

		if err := insert(db, tableName, row{id: n, word: entry.Headword.Text, info: info}); err != nil {
			return err
		}
	}

	return nil
}

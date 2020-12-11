package importer

import (
	"archive/zip"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type CollocationXMLEntry struct {
	XMLName xml.Name `xml:"entry"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text        string `xml:",chardata"`
		LexicalUnit struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"lexical_unit"`
		Category string `xml:"category"`
		Measure  struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Source string `xml:"source,attr"`
		} `xml:"measure"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Sense struct {
			Text     string `xml:",chardata"`
			Gramrels struct {
				Text    string `xml:",chardata"`
				Gramrel []struct {
					Text         string `xml:",chardata"`
					Name         string `xml:"name,attr"`
					Collocations struct {
						Text        string          `xml:",chardata"`
						Collocation CollocationsXML `xml:"collocation"`
					} `xml:"collocations"`
				} `xml:"gramrel"`
			} `xml:"gramrels"`
		} `xml:"sense"`
	} `xml:"body"`
}

type CollocationsXML []CollocationXML

var _ sort.Interface = CollocationsXML(nil)

func (c CollocationsXML) Len() int           { return len(c) }
func (c CollocationsXML) Less(i, j int) bool { return c[i].Frequency > c[j].Frequency }
func (c CollocationsXML) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type CollocationXML struct {
	Text      string  `xml:",chardata"`
	Form      string  `xml:"form,attr"`
	Cluster   string  `xml:"cluster,attr"`
	Frequency float64 `xml:"frequency,attr"`
	LogDice   string  `xml:"logDice,attr"`
	Comp      []struct {
		Text     string `xml:",chardata"`
		Position string `xml:"position,attr"`
		Type     string `xml:"type,attr"`
		Sloleks  string `xml:"sloleks,attr"`
	} `xml:"comp"`
}

func ImportCollocations(db *sql.DB, tableName string) error {
	// Open a zip archive for reading.
	r, err := zip.OpenReader("data/KSSS.zip")
	if err != nil {
		return err
	}
	defer r.Close()

	sumFrequencies := 0.0
	numFrequencies := 0
	id := -1
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".xml") {
			id++
			fc, err := f.Open()
			if err != nil {
				return err
			}
			byts, err := ioutil.ReadAll(fc)
			if err != nil {
				return err
			}
			var entry CollocationXMLEntry
			xml.Unmarshal(byts, &entry)
			//fmt.Println(entry.Header.Category, entry.Header.LexicalUnit.Text)
			var info [][][2]interface{}
			for n, gramrel := range entry.Body.Sense.Gramrels.Gramrel {
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
			insert(db, tableName, row{id: id, word: entry.Header.LexicalUnit.Text, info: info})
		}
	}
	fmt.Println("Avg collocation frequency:", sumFrequencies/float64(numFrequencies))
	return nil
}

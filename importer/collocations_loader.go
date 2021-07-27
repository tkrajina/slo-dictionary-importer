package importer

import (
	"archive/zip"
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

func (ce CollocationXMLEntry) Word() string {
	return ce.Header.LexicalUnit.Text
}

func (ce CollocationXMLEntry) GetFrequentCollocations() [][]string {
	var frequent [][]string
	for n, gramrel := range ce.Body.Sense.Gramrels.Gramrel {
		if n > 5 {
			continue
		}
		frequent = append(frequent, []string{})
		collocations := gramrel.Collocations.Collocation
		sort.Sort(collocations)
		if len(collocations) > 10 {
			collocations = collocations[0:10]
		}
		for _, collocation := range collocations {
			frequent[len(frequent)-1] = append(frequent[len(frequent)-1], collocation.Form)
		}
	}
	return frequent
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

func LoadCollocations() ([]CollocationXMLEntry, error) {
	ch := LoadCollocationsChan()
	var res []CollocationXMLEntry
	for e := range ch {
		if e.Err != nil {
			return nil, e.Err
		}
		res = append(res, e.Entry)
	}
	return res, nil
}

type CollocationEntryWithErr struct {
	Entry CollocationXMLEntry
	Err   error
}

func LoadCollocationsChan() <-chan CollocationEntryWithErr {
	// Open a zip archive for reading.
	res := make(chan CollocationEntryWithErr)
	r, err := zip.OpenReader("data/KSSS.zip")
	if err != nil {
		res <- CollocationEntryWithErr{Err: err}
		return res
	}

	go func() {
		for n, f := range r.File {
			if strings.HasSuffix(f.Name, ".xml") {
				if n%1000 == 0 {
					fmt.Println("Importing", n, "collocations")
				}
				fc, err := f.Open()
				if err != nil {
					res <- CollocationEntryWithErr{Err: err}
					return
				}
				byts, err := ioutil.ReadAll(fc)
				if err != nil {
					res <- CollocationEntryWithErr{Err: err}
					return
				}
				var entry CollocationXMLEntry
				if err := xml.Unmarshal(byts, &entry); err != nil {
					res <- CollocationEntryWithErr{Err: err}
					return
				}
				//fmt.Println(entry.Header.Category, entry.Header.LexicalUnit.Text)
				res <- CollocationEntryWithErr{Entry: entry}
			}
		}
		r.Close()
		close(res)
	}()

	return res
}

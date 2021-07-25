package importer

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
)

type ThesaurusEntry struct {
	Text     string `xml:",chardata"`
	Headword struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"headword"`
	GroupsCore struct {
		Text  string     `xml:",chardata"`
		Group []GroupXML `xml:"group"`
	} `xml:"groups_core"`
	GroupsNear struct {
		Text  string     `xml:",chardata"`
		Group []GroupXML `xml:"group"`
	} `xml:"groups_near"`
}

type ThesaurusXML struct {
	XMLName        xml.Name         `xml:"base"`
	Text           string           `xml:",chardata"`
	Xsi            string           `xml:"xsi,attr"`
	SchemaLocation string           `xml:"schemaLocation,attr"`
	Entries        []ThesaurusEntry `xml:"entry"`
}

type GroupXML struct {
	Text      string         `xml:",chardata"`
	Candidate []CandidateXML `xml:"candidate"`
}

type CandidateXML struct {
	Text  string `xml:",chardata"`
	Score string `xml:"score,attr"`
	S     struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"s"`
	Labels struct {
		Text string   `xml:",chardata"`
		La   []string `xml:"la"`
	} `xml:"labels"`
}

type Word [2]interface{}

func LoadThesaurus() ([]ThesaurusEntry, error) {
	// Open a zip archive for reading.
	r, err := zip.OpenReader("data/CJVT_Thesaurus-v1.0.zip")
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var byts []byte
	for _, f := range r.File {
		if f.Name == "CJVT_Thesaurus-v1.0.xml" {
			fc, err := f.Open()
			if err != nil {
				return nil, err
			}
			byts, err = ioutil.ReadAll(fc)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(byts) == 0 {
		return nil, fmt.Errorf("File not found")
	}

	var data ThesaurusXML
	fmt.Println("Unmarshalling thesaurus")
	err = xml.Unmarshal(byts, &data)
	if err != nil {
		return nil, err
	}
	fmt.Println("Unmarshalled thesaurus")

	for n, entry := range data.Entries {
		if n%1000 == 0 {
			fmt.Println("Importing", n, "thesaurus")
		}
		var info [][][]Word
		group, err := importGroup("Core", entry.GroupsCore.Group)
		if err != nil {
			return nil, err
		}
		info = append(info, group)
		group, err = importGroup("Near", entry.GroupsNear.Group)
		if err != nil {
			return nil, err
		}
		info = append(info, group)
	}

	return data.Entries, nil
}

func importGroup(name string, grps []GroupXML) ([][]Word, error) {
	//fmt.Println("	", name)
	res := [][]Word{}
	for n, grp := range grps {
		_ = n
		//fmt.Println("		Group #", n)
		candidates := []Word{}
		for _, candidate := range grp.Candidate {
			//fmt.Printf("			* [%s] %s\n", candidate.Score, candidate.S.Text)
			score, err := strconv.ParseFloat(candidate.Score, 64)
			if err != nil {
				return nil, err
			}
			//fmt.Println(CandidateJSON{Score: score, Word: candidate.S.Text})
			candidates = append(candidates, Word([2]interface{}{score, candidate.S.Text}))
		}
		res = append(res, candidates)
	}
	return res, nil
}

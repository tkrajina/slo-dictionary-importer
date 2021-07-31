package importer

import (
	"fmt"
	"html"
	"math"
	"os"
	"strings"
)

type KindleDictSynonym struct {
	Word  string  `json:"word"`
	Score float64 `json:"score"`
}

func (kds KindleDictSynonym) ScoreNormalized() float64 {
	return math.Max(0.0, math.Min(0.99, kds.Score/0.4))
}

type KindleDictEntry struct {
	ID           int                   `json:"id"`
	Word         string                `json:"word"`
	Inflections  []string              `json:"inflections"`
	Synonyms     [][]KindleDictSynonym `json:"synonyms"`
	Scores       [][]float64           `json:"synScores"`
	Collocations [][]string            `json:"collocations"`
	Frequency    string                `json:"frequency"`
}

func (de KindleDictEntry) toXML() string {
	res := `<idx:entry name="english" scriptable="yes" spell="yes"><idx:short><a id="` + fmt.Sprint(de.ID) + `"></a><idx:orth value="` + html.EscapeString(de.Word) + `"><h5>` + html.EscapeString(de.Word) + `</h5><idx:infl>`
	for _, infl := range de.Inflections {
		res += `<idx:iform value="` + html.EscapeString(infl) + `" />`
	}
	res += `</idx:infl></idx:orth>`

	var synonyms []string
	var collocations []string

	for _, syns := range de.Synonyms {
		for _, syn := range syns {
			synonyms = append(synonyms, "<i>syn.</i> "+html.EscapeString(syn.Word)) // TODO
		}
	}
col_loop:
	for n, colls := range de.Collocations {
		if n >= 2 {
			continue col_loop
		}
		for m, col := range colls {
			if m >= 2 {
				continue col_loop
			}
			collocations = append(collocations, "<i>e.g.</i> "+html.EscapeString(col))
		}
	}

	if len(synonyms) == 0 {
		return ""
	}
	res += "<p>" + strings.Join(synonyms, "<br>") + "</p>"

	if len(collocations) > 0 {
		res += "<p>" + strings.Join(collocations, "<br>") + "</p>"
	}

	if de.Frequency != "" {
		res += "<p><i>freq.</i> " + html.EscapeString(de.Frequency) + "</p>"
	}

	res += `</idx:short></idx:entry>
`
	return res
}

type KindleDict struct {
	Entries []KindleDictEntry `json:"entries"`
}

func ExportOPF(dict KindleDict) error {
	f, err := os.Create("kindledict/content.html")
	if err != nil {
		return err
	}

	_, err = f.WriteString(`<html xmlns:math="http://exslt.org/math" xmlns:svg="http://www.w3.org/2000/svg"
	xmlns:tl="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:saxon="http://saxon.sf.net/" xmlns:xs="http://www.w3.org/2001/XMLSchema"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xmlns:cx="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:mbp="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:mmc="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:idx="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf">

<head>
	<style>
	h5 {
		font-size: 1.1em;
		margin: 0.2em 0 0 0;
		padding: 0;
	}
	p {
		margin: 0 0 0.6em 0.5em;
		padding: 0;
	}
	</style>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>

<body>
	<mbp:frameset>`)
	if err != nil {
		return err
	}

	for _, w := range dict.Entries {
		_, err = f.WriteString(w.toXML())
		if err != nil {
			return err
		}
	}
	_, err = f.WriteString(`</mbp:frameset>
</body>
</html>`)
	return err
}

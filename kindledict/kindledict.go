package main

import (
	"fmt"
	"strings"
)

const dictHtml = `<html xmlns:math="http://exslt.org/math" xmlns:svg="http://www.w3.org/2000/svg"
	xmlns:tl="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:saxon="http://saxon.sf.net/" xmlns:xs="http://www.w3.org/2001/XMLSchema"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xmlns:cx="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:mbp="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:mmc="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:idx="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf">

<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>

<body>
	<mbp:frameset>
	</mbp:frameset>
</body>
</html>`

type dictEntry struct {
	word        string
	inflections []string
	description string
}

func (de dictEntry) toXML() string {
	res := `<idx:entry name="english" scriptable="yes" spell="yes">
	<idx:short>
		<a id="1"></a>
		<idx:orth value="` + de.word + `"><b>` + de.word + `</b>
			<idx:infl>`
	for _, infl := range de.inflections {
		res += `<idx:iform value="` + infl + `" />`
	}
	res += `			</idx:infl>
		</idx:orth>`
	res += `<p>first` + de.description + `</p>`
	res += `<p>second` + de.description + `</p>`
	res += ` </idx:short>
</idx:entry>`
	return res
}

type dict struct {
	entries []dictEntry
}

func main() {
	var dict dict
	for i := 0; i < 10; i++ {
		dict.entries = append(dict.entries, dictEntry{
			word:        strings.Repeat(fmt.Sprint(i), 5),
			inflections: []string{"aaa", "bbb", "ccc"},
			description: "<p>jkljkljdaklfjdskl <strong>kdflsj</strong> fdskl <i>dfsjklf djskl</i></p>",
		})
	}

	fmt.Println(`<html xmlns:math="http://exslt.org/math" xmlns:svg="http://www.w3.org/2000/svg"
	xmlns:tl="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:saxon="http://saxon.sf.net/" xmlns:xs="http://www.w3.org/2001/XMLSchema"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xmlns:cx="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:mbp="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:mmc="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf"
	xmlns:idx="https://kindlegen.s3.amazonaws.com/AmazonKindlePublishingGuidelines.pdf">

<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>

<body>
	<mbp:frameset>`)
	for _, w := range dict.entries {
		fmt.Println(w.toXML())
	}
	fmt.Println(`</mbp:frameset>
</body>
</html>`)
}

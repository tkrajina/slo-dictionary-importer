package importer

import (
	"os"
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

type KindleDictEntry struct {
	Word        string
	Inflections []string
	Description string
}

func (de KindleDictEntry) toXML() string {
	res := `<idx:entry name="english" scriptable="yes" spell="yes">
	<idx:short>
		<a id="1"></a>
		<idx:orth value="` + de.Word + `"><strong>` + de.Word + `</strong>
			<idx:infl>`
	for _, infl := range de.Inflections {
		res += `<idx:iform value="` + infl + `" />`
	}
	res += `			</idx:infl>
		</idx:orth>`
	res += de.Description
	res += ` </idx:short>
</idx:entry>`
	return res
}

type KindleDict struct {
	Entries []KindleDictEntry
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

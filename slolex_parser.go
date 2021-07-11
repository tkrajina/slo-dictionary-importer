package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//fn := "tml.xml"
	fn := "./data/Sloleks2.0.LMF/sloleks_clarin_2.0.xml"
	// byts, err := ioutil.ReadFile(fn)
	// panicIfErr(err)

	// var slolex Slolex
	// err = xml.Unmarshal(byts, &slolex)
	// panicIfErr(err)

	// fmt.Println(slolex)

	f, err := os.Open(fn)
	panicIfErr(err)
	d := xml.NewDecoder(f)
	count := 0
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			log.Fatalf("Error decoding token: %s", err)
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			if ty.Name.Local == "LexicalEntry" {
				count++
				var loc SlolexLexicalEntry
				panicIfErr(d.DecodeElement(&loc, &ty))
				if count%1000 == 0 {
					fmt.Printf("#%d. %#v\n", count, loc.Lema)
					//fmt.Printf("  forms: %#v\n", loc.Forms)
					fmt.Printf("  lema: %s\n", loc.Lema.FindLema())
					fmt.Printf("  representations: %v\n", loc.FindFormRepresentations())
				}
			}
		default:
		}
	}

	fmt.Println("count =", count)
}

/*
<LexicalEntry id="LE_a0f08b4f97c2606889f612e77f9d93f5">
  <feat att="ključ" val="S_Pierre"/>
  <feat att="besedna_vrsta" val="samostalnik"/>
  <feat att="SP2001" val="samostalnik"/>
  <feat att="vrsta" val="lastno_ime"/>
  <feat att="spol" val="moški"/>
  <feat att="SPSP" val="C1a1g"/>
  <Lemma>
    <feat att="zapis_oblike" val="Pierre"/>
    <feat att="naglašena_beseda_1" val="Pierrè"/>
  </Lemma>
  <WordForm>
    <feat att="msd" val="Slmei"/>
    <feat att="število" val="ednina"/>
    <feat att="sklon" val="imenovalnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierre"/>
      <feat att="pogostnost" val="6865"/>
      <feat att="naglašena_beseda_1" val="Pierrè"/>
      <feat att="SAMPA_1" val="piE&quot;rrE"/>
      <feat att="IPA_1" val="piɛˈrrɛ"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmer"/>
    <feat att="število" val="ednina"/>
    <feat att="sklon" val="rodilnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierra"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="1304"/>
      <feat att="naglašena_beseda_1" val="Piêrra"/>
      <feat att="SAMPA_1" val="pi&quot;E:rra"/>
      <feat att="IPA_1" val="piˈɛːrra"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrea"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="17"/>
      <feat att="naglašena_beseda_1" val="Pierrêa"/>
      <feat att="SAMPA_1" val="piErr&quot;E:a"/>
      <feat att="IPA_1" val="piɛrrˈɛːa"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmed"/>
    <feat att="število" val="ednina"/>
    <feat att="sklon" val="dajalnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierru"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="178"/>
      <feat att="naglašena_beseda_1" val="Piêrru"/>
      <feat att="SAMPA_1" val="pi&quot;E:rru"/>
      <feat att="IPA_1" val="piˈɛːrru"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreu"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="1"/>
      <feat att="naglašena_beseda_1" val="Pierrêu"/>
      <feat att="SAMPA_1" val="piErr&quot;E:u"/>
      <feat att="IPA_1" val="piɛrrˈɛːu"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmetd"/>
    <feat att="število" val="ednina"/>
    <feat att="sklon" val="tožilnik"/>
    <feat att="živost" val="da"/>
    <feat att="živost" val="da"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierra"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="185"/>
      <feat att="naglašena_beseda_1" val="Piêrra"/>
      <feat att="SAMPA_1" val="pi&quot;E:rra"/>
      <feat att="IPA_1" val="piˈɛːrra"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrea"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="1"/>
      <feat att="naglašena_beseda_1" val="Pierrêa"/>
      <feat att="SAMPA_1" val="piErr&quot;E:a"/>
      <feat att="IPA_1" val="piɛrrˈɛːa"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmem"/>
    <feat att="število" val="ednina"/>
    <feat att="sklon" val="mestnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierru"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="106"/>
      <feat att="naglašena_beseda_1" val="Piêrru"/>
      <feat att="SAMPA_1" val="pi&quot;E:rru"/>
      <feat att="IPA_1" val="piˈɛːrru"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreu"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="3"/>
      <feat att="naglašena_beseda_1" val="Pierrêu"/>
      <feat att="SAMPA_1" val="piErr&quot;E:u"/>
      <feat att="IPA_1" val="piɛrrˈɛːu"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmeo"/>
    <feat att="število" val="ednina"/>
    <feat att="sklon" val="orodnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrom"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="347"/>
      <feat att="naglašena_beseda_1" val="Pierrom"/>
      <feat att="SAMPA_1" val="piErrOm"/>
      <feat att="IPA_1" val="piɛrrɔm"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreom"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêom"/>
      <feat att="SAMPA_1" val="piErr&quot;E:Om"/>
      <feat att="IPA_1" val="piɛrrˈɛːɔm"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmdi"/>
    <feat att="število" val="dvojina"/>
    <feat att="sklon" val="imenovalnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierra"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrra"/>
      <feat att="SAMPA_1" val="pi&quot;E:rra"/>
      <feat att="IPA_1" val="piˈɛːrra"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrea"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêa"/>
      <feat att="SAMPA_1" val="piErr&quot;E:a"/>
      <feat att="IPA_1" val="piɛrrˈɛːa"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmdr"/>
    <feat att="število" val="dvojina"/>
    <feat att="sklon" val="rodilnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrov"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrrov"/>
      <feat att="SAMPA_1" val="pi&quot;E:rrOv"/>
      <feat att="IPA_1" val="piˈɛːrrɔv"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreov"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêov"/>
      <feat att="SAMPA_1" val="piErr&quot;E:Ov"/>
      <feat att="IPA_1" val="piɛrrˈɛːɔv"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmdd"/>
    <feat att="število" val="dvojina"/>
    <feat att="sklon" val="dajalnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierroma"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrroma"/>
      <feat att="SAMPA_1" val="pi&quot;E:rrOma"/>
      <feat att="IPA_1" val="piˈɛːrrɔma"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreoma"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêoma"/>
      <feat att="SAMPA_1" val="piErr&quot;E:Oma"/>
      <feat att="IPA_1" val="piɛrrˈɛːɔma"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmdt"/>
    <feat att="število" val="dvojina"/>
    <feat att="sklon" val="tožilnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierra"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrra"/>
      <feat att="SAMPA_1" val="pi&quot;E:rra"/>
      <feat att="IPA_1" val="piˈɛːrra"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrea"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêa"/>
      <feat att="SAMPA_1" val="piErr&quot;E:a"/>
      <feat att="IPA_1" val="piɛrrˈɛːa"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmdm"/>
    <feat att="število" val="dvojina"/>
    <feat att="sklon" val="mestnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrih"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrrih"/>
      <feat att="SAMPA_1" val="pi&quot;E:rrix"/>
      <feat att="IPA_1" val="piˈɛːrrix"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreih"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêih"/>
      <feat att="SAMPA_1" val="piErr&quot;E:ix"/>
      <feat att="IPA_1" val="piɛrrˈɛːix"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmdo"/>
    <feat att="število" val="dvojina"/>
    <feat att="sklon" val="orodnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierroma"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrroma"/>
      <feat att="SAMPA_1" val="pi&quot;E:rrOma"/>
      <feat att="IPA_1" val="piˈɛːrrɔma"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreoma"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêoma"/>
      <feat att="SAMPA_1" val="piErr&quot;E:Oma"/>
      <feat att="IPA_1" val="piɛrrˈɛːɔma"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmmi"/>
    <feat att="število" val="množina"/>
    <feat att="sklon" val="imenovalnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierri"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="9"/>
      <feat att="naglašena_beseda_1" val="Piêrri"/>
      <feat att="SAMPA_1" val="pi&quot;E:rri"/>
      <feat att="IPA_1" val="piˈɛːrri"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrei"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêi"/>
      <feat att="SAMPA_1" val="piErr&quot;E:i"/>
      <feat att="IPA_1" val="piɛrrˈɛːi"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmmr"/>
    <feat att="število" val="množina"/>
    <feat att="sklon" val="rodilnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrov"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrrov"/>
      <feat att="SAMPA_1" val="pi&quot;E:rrOv"/>
      <feat att="IPA_1" val="piˈɛːrrɔv"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreov"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêov"/>
      <feat att="SAMPA_1" val="piErr&quot;E:Ov"/>
      <feat att="IPA_1" val="piɛrrˈɛːɔv"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmmd"/>
    <feat att="število" val="množina"/>
    <feat att="sklon" val="dajalnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrom"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrrom"/>
      <feat att="SAMPA_1" val="pi&quot;E:rrOm"/>
      <feat att="IPA_1" val="piˈɛːrrɔm"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreom"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêom"/>
      <feat att="SAMPA_1" val="piErr&quot;E:Om"/>
      <feat att="IPA_1" val="piɛrrˈɛːɔm"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmmt"/>
    <feat att="število" val="množina"/>
    <feat att="sklon" val="tožilnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierre"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrre"/>
      <feat att="SAMPA_1" val="pi&quot;E:rrE"/>
      <feat att="IPA_1" val="piˈɛːrrɛ"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierree"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêe"/>
      <feat att="SAMPA_1" val="piErr&quot;E:E"/>
      <feat att="IPA_1" val="piɛrrˈɛːɛ"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmmm"/>
    <feat att="število" val="množina"/>
    <feat att="sklon" val="mestnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrih"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Piêrrih"/>
      <feat att="SAMPA_1" val="pi&quot;E:rrix"/>
      <feat att="IPA_1" val="piˈɛːrrix"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierreih"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêih"/>
      <feat att="SAMPA_1" val="piErr&quot;E:ix"/>
      <feat att="IPA_1" val="piɛrrˈɛːix"/>
    </FormRepresentation>
  </WordForm>
  <WordForm>
    <feat att="msd" val="Slmmo"/>
    <feat att="število" val="množina"/>
    <feat att="sklon" val="orodnik"/>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierri"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="tip" val="C1a1g_s_2"/>
      <feat att="pogostnost" val="2"/>
      <feat att="naglašena_beseda_1" val="Piêrri"/>
      <feat att="SAMPA_1" val="pi&quot;E:rri"/>
      <feat att="IPA_1" val="piˈɛːrri"/>
    </FormRepresentation>
    <FormRepresentation>
      <feat att="zapis_oblike" val="Pierrei"/>
      <feat att="SPSP" val="C1a1g"/>
      <feat att="norma" val="nestandardno"/>
      <feat att="tip" val="C1a1g_n_2"/>
      <feat att="pogostnost" val="0"/>
      <feat att="naglašena_beseda_1" val="Pierrêi"/>
      <feat att="SAMPA_1" val="piErr&quot;E:i"/>
      <feat att="IPA_1" val="piɛrrˈɛːi"/>
    </FormRepresentation>
  </WordForm>
  <RelatedForm>
    <feat att="idref" val="LE_51458bd527bbbd949f991eea706f708e"/>
    <feat att="besedna_vrsta" val="pridevnik"/>
    <feat att="lema" val="Pierrov"/>
  </RelatedForm>
</LexicalEntry>
*/

type SlolexFeat struct {
	// <feat att="lema" val="Pierrov"/>
	Att   string `xml:"att,attr"`
	Value string `xml:"val,attr"`
}

type SlolexLemma struct {
	//   <Lemma>
	//     <feat att="zapis_oblike" val="Pierre"/>
	//     <feat att="naglašena_beseda_1" val="Pierrè"/>
	//   </Lemma>
	Feats []SlolexFeat `xml:"feat"`
}

func (sl SlolexLemma) FindLema() string {
	for _, f := range sl.Feats {
		if f.Att == "zapis_oblike" {
			return f.Value
		}
	}
	return ""
}

type SlolexWordForm struct {
	Feats           []SlolexFeat               `xml:"feat"`
	Representations []SlolexFormRepresentation `xml:"FormRepresentation"`
}

func (swf SlolexWordForm) FindFormRepresentations() []string {
	var res []string
	for _, r := range swf.Representations {
		repr := r.FindRepresentation()
		if repr != "" {
			res = append(res, repr)
		}
	}
	return res
}

type SlolexFormRepresentation struct {
	Feats []SlolexFeat `xml:"feat"`
}

func (srf SlolexFormRepresentation) FindRepresentation() string {
	for _, f := range srf.Feats {
		if f.Att == "zapis_oblike" {
			return f.Value
		}
	}
	return ""
}

type SlolexLexicalEntry struct {
	Feats []SlolexFeat     `xml:"feat"`
	Lema  SlolexLemma      `xml:"Lemma"`
	Forms []SlolexWordForm `xml:"WordForm"`
}

func (sle SlolexLexicalEntry) FindFormRepresentations() [][]string {
	var res [][]string
	for _, f := range sle.Forms {
		res = append(res, f.FindFormRepresentations())
	}
	return res
}

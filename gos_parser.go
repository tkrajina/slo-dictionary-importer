package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type WordFrequency struct {
	Lemma        string
	PartOfSpeech string
	Frequency    float64
}

func main() {
	fn := "data/GOS1.0-words-all/GOS1.0-words-all-lemmas-parts_of_speech-taxonomy-entire.tsv"

	f, err := os.Open(fn)
	panicIfErr(err)

	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	var lemmasStarted bool
	var count int
	for scanner.Scan() {
		line := scanner.Text()
		if lemmasStarted {
			parts := strings.Split(line, "\t")
			freq, err := strconv.ParseFloat(strings.Replace(strings.Replace(strings.Trim(parts[4], `"`), ",", ".", 1), " %", "", 1), 32)
			panicIfErr(err)
			f := WordFrequency{
				Lemma:        strings.Trim(parts[0], `"`),
				PartOfSpeech: strings.Trim(parts[2], `"`),
				Frequency:    freq,
			}
			fmt.Println(count, f.Lemma, f.PartOfSpeech, f.Frequency)
			count++
		}
		if strings.HasPrefix(line, `"Lemma"`) {
			lemmasStarted = true
			fmt.Println("Lemmas started after:", line)
			continue
		}
	}
}

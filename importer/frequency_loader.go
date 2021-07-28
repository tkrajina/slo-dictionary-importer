package importer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type WordFrequency struct {
	Lemma        string
	PartOfSpeech string
	Frequency    float64
	Count        int
}

type WordFrequencyOrError struct {
	Entry WordFrequency
	Err   error
}

func LoadFrequencyChan() <-chan WordFrequencyOrError {
	res := make(chan WordFrequencyOrError)

	fn := "data/GOS1.0-words-all/GOS1.0-words-all-lemmas-parts_of_speech-taxonomy-entire.tsv"
	f, err := os.Open(fn)
	if err != nil {
		res <- WordFrequencyOrError{Err: err}
		return res
	}

	go func() {
		defer close(res)

		scanner := bufio.NewScanner(f)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		var lemmasStarted bool
		var count int
		for scanner.Scan() {
			line := scanner.Text()
			if lemmasStarted {
				parts := strings.Split(line, "\t")
				freq, err := strconv.ParseFloat(strings.Replace(strings.Replace(strings.Trim(parts[4], `"`), ",", ".", 1), " %", "", 1), 32)
				if err != nil {
					res <- WordFrequencyOrError{Err: err}
					return
				}
				count++
				res <- WordFrequencyOrError{Entry: WordFrequency{
					Lemma:        strings.Trim(parts[0], `"`),
					PartOfSpeech: strings.Trim(parts[2], `"`),
					Frequency:    freq,
					Count:        count,
				}}
				if count%1000 == 0 {
					fmt.Println(count, "word frequencies")
				}
			}
			if strings.HasPrefix(line, `"Lemma"`) {
				lemmasStarted = true
				fmt.Println("Lemmas started after:", line)
				continue
			}
		}
	}()

	return res
}

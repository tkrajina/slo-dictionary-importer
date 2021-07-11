package main

import (
	"fmt"

	"github.com/tkrajina/slo-dictionary-importer/importer"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	entries, err := importer.SlolexLoader()
	panicIfErr(err)
	fmt.Println("Loaded %d slolex entries\n", len(entries))
}

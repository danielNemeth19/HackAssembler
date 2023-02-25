package main

import (
	"github.com/pkg/profile"
	"log"
	"os"
	"time"
)

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func parseArg() string {
	if len(os.Args) != 2 {
		log.Fatalf("Incorrect number of arguments. Got: %d", len(os.Args)-1)
	}
	input := os.Args[1]
	log.Printf("Source to translate: %s", input)
	return input
}

func main() {
	defer profile.Start(profile.ProfilePath("profiling/")).Stop()
	defer TimeTrack(time.Now(), "main")
	source := parseArg()

	var parser Parser
	var symbolTable SymbolTable
	var translator HackTranslator
	parser.Initialize(source)
	symbolTable.Initialize()
	translator.Initialize()

	for _, sf := range parser.FilesToTranslate {
		parser.TranslateFile(sf, &symbolTable, translator)
	}
}

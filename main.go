package main

import (
	"fmt"
	"github.com/pkg/profile"
	"log"
	"os"
	"path/filepath"
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
	parser := Parser{Source: source}
	var symbolTable SymbolTable
	var translator HackTranslator
	symbolTable.Initialize()
	translator.Initialize()

	if isDir := parser.IsSourceDir(); isDir == false {
		parser.TranslateFile(parser.Source, &symbolTable, translator)
	} else {
		fmt.Printf("Source is a directory: %s\n", parser.Source)
		e := filepath.Walk(parser.Source, func(path string, f os.FileInfo, err error) error {
			fmt.Printf("here: %s\n", path)
			if extension := filepath.Ext(f.Name()); extension == ".asm" {
				fmt.Printf("Input file needs to be translated: %s\n", f.Name())
				parser.Source = f.Name()
				parser.TranslateFile(f.Name(), &symbolTable, translator)
			}
			return err
		})
		CheckError(e)
	}
}

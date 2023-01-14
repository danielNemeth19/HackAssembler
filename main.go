package main

import (
	"fmt"
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
	inputFile := os.Args[1]
	if extension := filepath.Ext(inputFile); extension != ".asm" {
		log.Fatalf("Input file needs to be an asm file - received: %s", extension)
	}
	log.Printf("File to translate: %s", inputFile)
	return inputFile
}

func main() {
	defer TimeTrack(time.Now(), "main")
	fp := parseArg()
	parser := Parser{SourceFile: fp}
	var symbolTable SymbolTable
	symbolTable.Initialize()
	parser.GetAssemblyCode(&symbolTable)
	translator := HackTranslator{}
	translator.Initialize()

	df := parser.SetDestinationFile()
	f, err := os.Create(df)
	CheckError(err)
	for _, codeSnippet := range parser.AssemblyCode {
		if parser.IsAInstruction(codeSnippet) {
			address := parser.GetAddress(codeSnippet, &symbolTable)
			code := translator.TranslateAInstruction(address)
			fmt.Fprintln(f, code)
		} else {
			comp, dest, jmp := parser.ParseCInstruction(codeSnippet)
			compCode := translator.TranslateComp(comp)
			destCode := translator.TranslateDest(dest)
			jmpCode := translator.TranslateJmp(jmp)
			code := translator.TranslateCInstruction(compCode, destCode, jmpCode)
			fmt.Fprintln(f, code)
		}
	}
	f.Close()
}

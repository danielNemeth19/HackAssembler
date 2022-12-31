package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type AssemblyCode struct {
	lineNumber int
	codeLines  []string
}

func readHackFile(fp string) AssemblyCode {
	file, err := os.Open(fp)
	checkError(err)
	scanner := bufio.NewScanner(file)
	code := AssemblyCode{}
	for scanner.Scan() {
		line := scanner.Text()
		subs := strings.Split(line, "//")
		if len(subs[0]) > 0 {
			code.codeLines = append(code.codeLines, strings.TrimSpace(subs[0]))
			code.lineNumber = code.lineNumber + 1
		}
	}
	return code
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
	fp := parseArg()
	assemblyCode := readHackFile(fp)
	for index := range assemblyCode.codeLines {
		fmt.Println(assemblyCode.codeLines[index], index)
	}
	fmt.Printf("Number of code lines: %d\n", assemblyCode.lineNumber)
}

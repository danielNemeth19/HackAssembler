package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type SymbolTable struct {
	table       map[string]int
	nextFreeMem int
}

func (symbolTable *SymbolTable) initialize() {
	initValues := map[string]int{
		"SP": 0, "LCL": 1, "ARG": 2, "THIS": 3, "THAT": 4, "SCREEN": 16384, "KBD": 24576,
	}
	for i := 0; i < 16; i++ {
		key := "R" + strconv.Itoa(i)
		initValues[key] = i
	}
	symbolTable.table = initValues
	symbolTable.nextFreeMem = 16
}

func (symbolTable *SymbolTable) storeLabel(codeSnippet string, counter int) {
	label := codeSnippet[1 : len(codeSnippet)-1]
	symbolTable.table[label] = counter
}

func (symbolTable *SymbolTable) getAddress(symbol string) (int, bool) {
	address, found := symbolTable.table[symbol]
	return address, found
}

func (symbolTable *SymbolTable) storeVariable(variable string) int {
	address := symbolTable.nextFreeMem
	symbolTable.table[variable] = address
	symbolTable.nextFreeMem++
	fmt.Printf("nextMem: %d", symbolTable.nextFreeMem)
	return address
}

type Parser struct {
	symbolTable *SymbolTable
}

func (parser Parser) isAInstruction(token string) bool {
	if strings.HasPrefix(token, "@") {
		return true
	}
	return false
}

func (parser Parser) getAddress(token string) int {
	address, err := strconv.Atoi(token[1:])
	if err != nil {
		fmt.Printf("%s must be a symbol\n", token)
		address, found := parser.symbolTable.getAddress(token[1:])
		if found != true {
			fmt.Printf("Address: %d -- Found: %v\n", address, found)
			address = parser.symbolTable.storeVariable(token[1:])
		}
	}
	return address
}

type HackCode struct{}

func (code HackCode) translateAInstruction(address int) {
	mys := fmt.Sprintf("0%015b\n", address)
	fmt.Printf("Bit representation: %s", mys)
}

func readHackFile(fp string) ([]string, SymbolTable) {
	var symbolTable SymbolTable
	symbolTable.initialize()
	file, err := os.Open(fp)
	checkError(err)
	scanner := bufio.NewScanner(file)
	var assemblyCode []string
	counter := -1
	for scanner.Scan() {
		line := scanner.Text()
		subs := strings.Split(line, "//")
		if len(subs[0]) > 0 {
			codeSnippet := strings.TrimSpace(subs[0])
			if strings.HasPrefix(codeSnippet, "(") {
				symbolTable.storeLabel(codeSnippet, counter+1)
			} else {
				counter += 1
				assemblyCode = append(assemblyCode, codeSnippet)
			}
		}
	}
	return assemblyCode, symbolTable
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
	assemblyCode, symbolTable := readHackFile(fp)
	parser := Parser{symbolTable: &symbolTable}
	hackCode := HackCode{}
	for _, code := range assemblyCode {
		if parser.isAInstruction(code) {
			fmt.Printf("This is A instruction: %s\n", code)
			address := parser.getAddress(code)
			hackCode.translateAInstruction(address)
		}
	}
	fmt.Printf("Number items in symbol table: %d\n", len(symbolTable.table))
	fmt.Printf("symbol table: %v\n", symbolTable.table)
}

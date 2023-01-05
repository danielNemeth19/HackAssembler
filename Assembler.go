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
	return address
}

type Parser struct{}

func (parser Parser) isAInstruction(token string) bool {
	if strings.HasPrefix(token, "@") {
		return true
	}
	return false
}

func (parser Parser) getAddress(token string, table *SymbolTable) int {
	token = token[1:]
	address, err := strconv.Atoi(token)
	if err == nil {
		return address
	}
	address, found := table.getAddress(token)
	if found != true {
		address = table.storeVariable(token)
	}
	return address
}

func (parser Parser) parseCInstruction(codeSnippet string) (string, string, string) {
	var comp, dest, jmp string
	subs := strings.Split(codeSnippet, "=")
	if len(subs) == 2 {
		dest, comp, jmp = subs[0], subs[1], "null"
	} else {
		subs = strings.Split(codeSnippet, ";")
		dest, comp, jmp = "null", subs[0], subs[1]
	}
	return comp, dest, jmp
}

type HackTranslator struct{}

func (translator HackTranslator) translateAInstruction(address int) string {
	machineCode := fmt.Sprintf("0%015b", address)
	//fmt.Printf("Bit representation: %s\n", machineCode)
	return machineCode
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
	parser := Parser{}
	translator := HackTranslator{}
	var hackCode []string
	for _, codeSnippet := range assemblyCode {
		if parser.isAInstruction(codeSnippet) {
			//fmt.Printf("This is A instruction: %s\n", codeSnippet)
			address := parser.getAddress(codeSnippet, &symbolTable)
			code := translator.translateAInstruction(address)
			hackCode = append(hackCode, code)
		} else {
			comp, dest, jmp := parser.parseCInstruction(codeSnippet)
			fmt.Printf("%s snippet -> %s %s %s\n", codeSnippet, comp, dest, jmp)
		}
	}

	fmt.Printf("Number items in symbol table: %d\n", len(symbolTable.table))
	fmt.Printf("symbol table: %v\n", symbolTable.table)
	for _, c := range hackCode {
		fmt.Printf("%s\n", c)
	}
}

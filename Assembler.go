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
	return machineCode
}

func (translator HackTranslator) translateComp(comp string) string {
	compMap := map[string]string{
		"0": "0101010", "1": "0111111", "-1": "0111010", "D": "0001100", "A": "0110000", "M": "1110000",
		"!D": "0001101", "!A": "0110001", "!M": "1110001", "-D": "0001111", "-A": "0110011", "-M": "1110011",
		"D+1": "0011111", "A+1": "0110111", "M+1": "1110111", "D-1": "0001110", "A-1": "0110010", "M-1": "1110010",
		"D+A": "0000010", "D+M": "1000010", "D-A": "0010011", "D-M": "1010011", "A-D": "0000111",
		"M-D": "1000111", "D&A": "0000000", "D&M": "1000000", "D|A": "0010101", "D|M": "1010101",
	}
	hackCode, found := compMap[comp]
	if found != true {
		log.Fatalf("Comp %s invalid", comp)
	}
	return hackCode
}

func (translator HackTranslator) translateDest(dest string) string {
	destMap := map[string]string{
		"null": "000", "M": "001", "D": "010", "MD": "011",
		"A": "100", "AM": "101", "AD": "110", "AMD": "111",
	}
	hackCode, found := destMap[dest]
	if found != true {
		log.Fatalf("Dest %s invalid", dest)
	}
	return hackCode
}

func (translator HackTranslator) translateJmp(jmp string) string {
	jmpMap := map[string]string{
		"null": "000", "JGT": "001", "JEQ": "010", "JGE": "011",
		"JLT": "100", "JNE": "101", "JLE": "110", "JMP": "111",
	}
	jmpCode, found := jmpMap[jmp]
	if found != true {
		log.Fatalf("Jmp %s invalid", jmp)
	}
	return jmpCode
}

func (translator HackTranslator) translateCInstruction(comp, dest, jmp string) string {
	code := "111" + comp + dest + jmp
	return code
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
	
	f, err := os.Create("test.hack")
	checkError(err)
	for _, codeSnippet := range assemblyCode {
		if parser.isAInstruction(codeSnippet) {
			address := parser.getAddress(codeSnippet, &symbolTable)
			code := translator.translateAInstruction(address)
			fmt.Fprintln(f, code)
		} else {
			comp, dest, jmp := parser.parseCInstruction(codeSnippet)
			compCode := translator.translateComp(comp)
			destCode := translator.translateDest(dest)
			jmpCode := translator.translateJmp(jmp)
			code := translator.translateCInstruction(compCode, destCode, jmpCode)
			fmt.Fprintln(f, code)
		}
	}
	f.Close()
}

package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	SourceFile   string
	AssemblyCode []string
}

func (parser *Parser) SetDestinationFile() string {
	hackExt := ".hack"
	return parser.SourceFile[:len(parser.SourceFile)-len(".asm")] + hackExt
}

func (parser *Parser) GetAssemblyCode(symbolTable *SymbolTable) {
	file, err := os.Open(parser.SourceFile)
	CheckError(err)
	scanner := bufio.NewScanner(file)
	counter := -1
	for scanner.Scan() {
		line := scanner.Text()
		subs := strings.Split(line, "//")
		if len(subs[0]) > 0 {
			codeSnippet := strings.TrimSpace(subs[0])
			if parser.IsLabel(codeSnippet) {
				symbolTable.StoreLabel(codeSnippet, counter+1)
			} else {
				counter += 1
				parser.AssemblyCode = append(parser.AssemblyCode, codeSnippet)
			}
		}
	}
	return
}

func (parser *Parser) IsLabel(token string) bool {
	if strings.HasPrefix(token, "(") {
		return true
	}
	return false
}

func (parser *Parser) IsAInstruction(token string) bool {
	if strings.HasPrefix(token, "@") {
		return true
	}
	return false
}

func (parser *Parser) GetAddress(token string, table *SymbolTable) int {
	token = token[1:]
	address, err := strconv.Atoi(token)
	if err == nil {
		return address
	}
	address, found := table.GetAddress(token)
	if found != true {
		address = table.StoreVariable(token)
	}
	return address
}

func (parser *Parser) ParseCInstruction(codeSnippet string) (string, string, string) {
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

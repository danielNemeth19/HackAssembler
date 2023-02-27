package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Parser struct {
	FilesToTranslate []string
}

func (parser *Parser) Initialize(sourceInput string) {
	e := filepath.Walk(sourceInput, func(path string, f os.FileInfo, err error) error {
		if extension := filepath.Ext(f.Name()); extension == ".asm" {
			fullName, _ := filepath.Abs(path)
			parser.FilesToTranslate = append(parser.FilesToTranslate, fullName)
		}
		return err
	})
	CheckError(e)
}

func (parser *Parser) IsSourceFile(f os.FileInfo) bool {
	if extension := filepath.Ext(f.Name()); extension == ".asm" {
		return true
	}
	return false
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

func (parser *Parser) SetDestinationFile(fn string) string {
	hackExt := ".hack"
	return fn[:len(fn)-len(".asm")] + hackExt
}

func (parser *Parser) GetAssemblyCode(fn string, symbolTable *SymbolTable) []string {
	file, err := os.Open(fn)
	CheckError(err)
	scanner := bufio.NewScanner(file)
	counter := -1
	var assemblyCode []string
	for scanner.Scan() {
		line := scanner.Text()
		subs := strings.Split(line, "//")
		if len(subs[0]) > 0 {
			codeSnippet := strings.TrimSpace(subs[0])
			if parser.IsLabel(codeSnippet) {
				symbolTable.StoreLabel(codeSnippet, counter+1)
			} else {
				counter += 1
				assemblyCode = append(assemblyCode, codeSnippet)
			}
		}
	}
	return assemblyCode
}

func (parser *Parser) GetAddress(token string, table *SymbolTable) int {
	token = token[1:]
	address, err := strconv.Atoi(token)
	if err == nil {
		return address
	}
	address, found := table.GetSymbolAddress(token)
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

func (parser *Parser) TranslateFile(fn string, st *SymbolTable, tr HackTranslator) {
	assemblyCode := parser.GetAssemblyCode(fn, st)
	df := parser.SetDestinationFile(fn)
	f, err := os.Create(df)
	CheckError(err)
	buffer := bufio.NewWriter(f)
	for _, codeSnippet := range assemblyCode {
		if parser.IsAInstruction(codeSnippet) {
			address := parser.GetAddress(codeSnippet, st)
			code := tr.TranslateAInstruction(address)
			buffer.WriteString(code + "\n")
		} else {
			comp, dest, jmp := parser.ParseCInstruction(codeSnippet)
			compCode := tr.TranslateComp(comp)
			destCode := tr.TranslateDest(dest)
			jmpCode := tr.TranslateJmp(jmp)
			code := tr.TranslateCInstruction(compCode, destCode, jmpCode)
			buffer.WriteString(code + "\n")
		}
	}
	err = buffer.Flush()
	CheckError(err)
	f.Close()
}

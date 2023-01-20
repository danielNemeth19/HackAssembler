package main

import (
	"testing"
)

type inputTable struct {
	code           string
	isLabel        bool
	isAInstruction bool
}

type addressTable struct {
	token   string
	address int
}

func makeParser() Parser {
	return Parser{SourceFile: "test.asm"}
}

func makeSymbolTable() SymbolTable {
	var symbolTable SymbolTable
	symbolTable.Initialize()
	return symbolTable
}

func TestParser_SetDestinationFile(t *testing.T) {
	p := makeParser()
	expectedPath := "test.hack"
	destPath := p.SetDestinationFile()
	if destPath != expectedPath {
		t.Errorf("Path incorrect: got %s, expected: %s", destPath, expectedPath)
	}
}

func TestParser_IsLabel_False(t *testing.T) {
	p := makeParser()
	res := p.IsLabel("not a label")
	if res == true {
		t.Errorf("Result incorrect: got %v, expected false\n", res)
	}
}

func TestParser_IsLabel_True(t *testing.T) {
	p := makeParser()
	res := p.IsLabel("(thisIsALabel)")
	if res != true {
		t.Errorf("Result incorrect: got %v, expected true\n", res)
	}
}

func TestParser_IsLabel_TestTable(t *testing.T) {
	p := makeParser()
	table := []inputTable{
		{code: "(label1)", isLabel: true},
		{code: "A=M", isLabel: false},
		{code: "(label2)", isLabel: true},
		{code: "0;JMP", isLabel: false},
	}
	for _, data := range table {
		res := p.IsLabel(data.code)
		if res != data.isLabel {
			t.Errorf("Result incorrect: got %v, expected %v\n", res, data.isLabel)
		}
	}
}

func TestParser_IsAInstruction(t *testing.T) {
	p := makeParser()
	table := []inputTable{{code: "@256", isAInstruction: true}, {code: "0;JMP", isAInstruction: false}}
	for _, data := range table {
		res := p.IsAInstruction(data.code)
		if res != data.isAInstruction {
			t.Errorf("Result incorrect: got %v, expected %v\n", res, data.isAInstruction)
		}
	}
}

func TestParser_GetAddress(t *testing.T) {
	st := makeSymbolTable()
	p := makeParser()
	resultTable := []addressTable{
		{token: "@SP", address: 0},
		{token: "@99", address: 99},
		{token: "@first_var", address: 16},
		{token: "@second_var", address: 17},
	}
	for _, data := range resultTable {
		res := p.GetAddress(data.token, &st)
		if res != data.address {
			t.Errorf("Result incorrect: got %v, expected %v\n", res, data.address)
		}
	}
}

func TestParser_ParseCInstruction_No_Jump(t *testing.T) {
	p := makeParser()
	code := "M=D"
	comp, dest, jmp := p.ParseCInstruction(code)
	if comp != "D" || dest != "M" || jmp != "null" {
		t.Errorf("Got %s %s %s - expected D M, null\n", comp, dest, jmp)
	}
}

func TestParser_ParseCInstruction_Jump(t *testing.T) {
	p := makeParser()
	code := "D;JGT"
	comp, dest, jmp := p.ParseCInstruction(code)
	if comp != "D" || dest != "null" || jmp != "JGT" {
		t.Errorf("Got %s %s %s - expected D null JGT\n", comp, dest, jmp)
	}
}

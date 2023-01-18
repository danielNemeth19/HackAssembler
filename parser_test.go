package main

import (
	"testing"
)

type inputTable struct {
	code           string
	isLabel        bool
	isAInstruction bool
}

func makeParser() Parser {
	return Parser{SourceFile: "test.asm"}
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

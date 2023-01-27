package main

import (
	"io/fs"
	"testing"
	"time"
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

func makeSymbolTable() SymbolTable {
	var symbolTable SymbolTable
	symbolTable.Initialize()
	return symbolTable
}

type FakeFile struct {
	name     string
	contents string
	mode     fs.FileMode
	size     int64
}

func (ff FakeFile) Name() string {
	return ff.name
}

func (ff FakeFile) IsDir() bool {
	return false
}

func (ff FakeFile) ModTime() time.Time {
	return time.Time{}
}

func (ff FakeFile) Mode() fs.FileMode {
	return ff.mode
}

func (ff FakeFile) Size() int64 {
	return ff.size
}

func (ff FakeFile) Sys() any {
	return nil
}

func TestParser_IsSourceFile(t *testing.T) {
	p := Parser{}
	fFile := FakeFile{name: "file.asm"}
	res := p.IsSourceFile(fFile)
	if res != true {
		t.Errorf("Result incorrect: got %v, expected true\n", res)
	}

	fFile = FakeFile{name: "notsource.txt"}
	res = p.IsSourceFile(fFile)
	if res == true {
		t.Errorf("Result incorrect: got %v, expected false\n", res)
	}
}

func TestParser_IsLabel_False(t *testing.T) {
	p := Parser{}
	if res := p.IsLabel("not a label"); res == true {
		t.Errorf("Result incorrect: got %v, expected false\n", res)
	}
}

func TestParser_IsLabel_True(t *testing.T) {
	p := Parser{}
	if res := p.IsLabel("(thisIsLabel)"); res != true {
		t.Errorf("Result incorrect: got %v, expected true\n", res)
	}
}

func TestParser_IsLabel_TestTable(t *testing.T) {
	p := Parser{}
	table := []inputTable{
		{code: "(label1)", isLabel: true},
		{code: "A=M", isLabel: false},
		{code: "(label2)", isLabel: true},
		{code: "0;JMP", isLabel: false},
	}
	for _, data := range table {
		if res := p.IsLabel(data.code); res != data.isLabel {
			t.Errorf("Result incorrect: got %v, expected %v\n", res, data.isLabel)
		}
	}
}

func TestParser_IsAInstruction(t *testing.T) {
	p := Parser{}
	table := []inputTable{{code: "@256", isAInstruction: true}, {code: "0;JMP", isAInstruction: false}}
	for _, data := range table {
		if res := p.IsAInstruction(data.code); res != data.isAInstruction {
			t.Errorf("Result incorrect: got %v, expected %v\n", res, data.isAInstruction)
		}
	}
}

func TestParser_SetDestinationFile(t *testing.T) {
	p := Parser{}
	expectedPath := "test.hack"
	destPath := p.SetDestinationFile("test.asm")
	if destPath != expectedPath {
		t.Errorf("Path incorrect: got %s, expected: %s", destPath, expectedPath)
	}
}

func TestParser_GetAddress(t *testing.T) {
	st := makeSymbolTable()
	p := Parser{}
	resultTable := []addressTable{
		{token: "@SP", address: 0},
		{token: "@99", address: 99},
		{token: "@first_var", address: 16},
		{token: "@second_var", address: 17},
	}
	for _, data := range resultTable {
		if res := p.GetAddress(data.token, &st); res != data.address {
			t.Errorf("Result incorrect: got %v, expected %v\n", res, data.address)
		}
	}
}

func TestParser_ParseCInstruction_No_Jump(t *testing.T) {
	p := Parser{}
	code := "M=D"
	comp, dest, jmp := p.ParseCInstruction(code)
	if comp != "D" || dest != "M" || jmp != "null" {
		t.Errorf("Got %s %s %s - expected D M, null\n", comp, dest, jmp)
	}
}

func TestParser_ParseCInstruction_Jump(t *testing.T) {
	p := Parser{}
	code := "D;JGT"
	comp, dest, jmp := p.ParseCInstruction(code)
	if comp != "D" || dest != "null" || jmp != "JGT" {
		t.Errorf("Got %s %s %s - expected D null JGT\n", comp, dest, jmp)
	}
}

func BenchmarkParser_ParseCInstruction(b *testing.B) {
	p := Parser{}
	for i := 0; i < b.N; i++ {
		code := "D;JGT"
		p.ParseCInstruction(code)
	}
}

func BenchmarkParser_GetAddress(b *testing.B) {
	p := Parser{}
	st := makeSymbolTable()
	for i := 0; i < b.N; i++ {
		p.GetAddress("@newVar", &st)
	}
}

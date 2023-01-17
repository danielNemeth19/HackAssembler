package main

import "testing"

func TestParser_IsLabel_False(t *testing.T) {
	var p Parser
	p = Parser{SourceFile: "test.asm"}
	res := p.IsLabel("not a label")
	if res == true {
		t.Errorf("Result incorrect: got %v, expected false\n", res)
	}
}

func TestParser_IsLabel_True(t *testing.T) {
	var p Parser
	p = Parser{SourceFile: "test.asm"}
	res := p.IsLabel("(thisIsALabel)")
	if res != true {
		t.Errorf("Result incorrect: got %v, expected true\n", res)
	}
}

func TestParser_IsLabel_TestTable(t *testing.T) {
	var p Parser
	p = Parser{SourceFile: "test.asm"}
	table := []struct {
		code    string
		isLabel bool
	}{
		{"(label1)", true},
		{"A=M", false},
		{"(label2)", true},
		{"0;JMP", false},
	}
	for _, data := range table {
		res := p.IsLabel(data.code)
		if res != data.isLabel {
			t.Errorf("Result incorrect: got %v, expected %v\n", res, data.isLabel)
		}
	}
}

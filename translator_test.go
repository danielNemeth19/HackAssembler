package main

import (
	"testing"
)

var compMap = map[string]string{
	"0": "0101010", "1": "0111111", "-1": "0111010", "D": "0001100", "A": "0110000", "M": "1110000",
	"!D": "0001101", "!A": "0110001", "!M": "1110001", "-D": "0001111", "-A": "0110011", "-M": "1110011",
	"D+1": "0011111", "A+1": "0110111", "M+1": "1110111", "D-1": "0001110", "A-1": "0110010", "M-1": "1110010",
	"D+A": "0000010", "D+M": "1000010", "D-A": "0010011", "D-M": "1010011", "A-D": "0000111",
	"M-D": "1000111", "D&A": "0000000", "D&M": "1000000", "D|A": "0010101", "D|M": "1010101",
}

var destMap = map[string]string{
	"null": "000", "M": "001", "D": "010", "MD": "011",
	"A": "100", "AM": "101", "AD": "110", "AMD": "111",
}

var jmpMap = map[string]string{
	"null": "000", "JGT": "001", "JEQ": "010", "JGE": "011",
	"JLT": "100", "JNE": "101", "JLE": "110", "JMP": "111",
}

func TestHackTranslator_Initialize_CompMap(t *testing.T) {
	var translator HackTranslator
	translator.Initialize()

	for k, v := range compMap {
		if cValue := translator.compMap[k]; cValue != v {
			t.Errorf("Result incorrect: got %s, expected %s\n", cValue, v)
		}
	}
}

func TestHackTranslator_Initialize_DestMap(t *testing.T) {
	var translator HackTranslator
	translator.Initialize()

	for k, v := range destMap {
		if dValue := translator.destMap[k]; dValue != v {
			t.Errorf("Result incorrect: got %s, expected %s\n", dValue, v)
		}
	}
}

func TestHackTranslator_Initialize_JumpMap(t *testing.T) {
	var translator HackTranslator
	translator.Initialize()

	for k, v := range jmpMap {
		if jValue := translator.jmpMap[k]; jValue != v {
			t.Errorf("Result incorrect: got %s, expected %s\n", jValue, v)
		}
	}
}

func TestHackTranslator_TranslateAInstruction(t *testing.T) {
	tr := HackTranslator{}
	expected := "0000000000011010"
	if code := tr.TranslateAInstruction(26); code != expected {
		t.Errorf("Result incorrect: got %s, expected %s\n", code, expected)
	}
}

func TestHackTranslator_TranslateComp(t *testing.T) {
	tr := HackTranslator{}
	tr.Initialize()

	for c, v := range compMap {
		if code := tr.TranslateComp(c); code != v {
			t.Errorf("For comp %s result incorrect: got %s, expected %s\n", c, code, v)
		}
	}
}

func TestHackTranslator_TranslateDest(t *testing.T) {
	tr := HackTranslator{}
	tr.Initialize()

	for d, v := range destMap {
		if code := tr.TranslateDest(d); code != v {
			t.Errorf("For dest %s result incorrect: got %s, expected %s\n", d, code, v)
		}
	}
}

func TestHackTranslator_TranslateJmp(t *testing.T) {
	tr := HackTranslator{}
	tr.Initialize()

	for j, v := range jmpMap {
		if code := tr.TranslateJmp(j); code != v {
			t.Errorf("For jmp %s result incorrect: got %s, expected %s\n", j, code, v)
		}
	}
}

func TestHackTranslator_TranslateCInstruction(t *testing.T) {
	tr := HackTranslator{}
	tr.Initialize()

	code := tr.TranslateCInstruction("0110000", "010", "000")
	if code != "1110110000010000" {
		t.Errorf("Result incorrect: got %s, expected %s\n", code, "")
	}
}

package main

import (
	"fmt"
	"log"
)

type HackTranslator struct {
	compMap, destMap, jmpMap map[string]string
}

func (translator *HackTranslator) Initialize() {
	translator.compMap = map[string]string{
		"0": "0101010", "1": "0111111", "-1": "0111010", "D": "0001100", "A": "0110000", "M": "1110000",
		"!D": "0001101", "!A": "0110001", "!M": "1110001", "-D": "0001111", "-A": "0110011", "-M": "1110011",
		"D+1": "0011111", "A+1": "0110111", "M+1": "1110111", "D-1": "0001110", "A-1": "0110010", "M-1": "1110010",
		"D+A": "0000010", "D+M": "1000010", "D-A": "0010011", "D-M": "1010011", "A-D": "0000111",
		"M-D": "1000111", "D&A": "0000000", "D&M": "1000000", "D|A": "0010101", "D|M": "1010101",
	}
	translator.destMap = map[string]string{
		"null": "000", "M": "001", "D": "010", "MD": "011",
		"A": "100", "AM": "101", "AD": "110", "AMD": "111",
	}
	translator.jmpMap = map[string]string{
		"null": "000", "JGT": "001", "JEQ": "010", "JGE": "011",
		"JLT": "100", "JNE": "101", "JLE": "110", "JMP": "111",
	}
}

func (translator *HackTranslator) TranslateAInstruction(address int) string {
	machineCode := fmt.Sprintf("0%015b", address)
	return machineCode
}

func (translator *HackTranslator) TranslateComp(comp string) string {
	hackCode, found := translator.compMap[comp]
	if found != true {
		log.Fatalf("Comp %s invalid", comp)
	}
	return hackCode
}

func (translator *HackTranslator) TranslateDest(dest string) string {
	hackCode, found := translator.destMap[dest]
	if found != true {
		log.Fatalf("Dest %s invalid", dest)
	}
	return hackCode
}

func (translator *HackTranslator) TranslateJmp(jmp string) string {
	jmpCode, found := translator.jmpMap[jmp]
	if found != true {
		log.Fatalf("Jmp %s invalid", jmp)
	}
	return jmpCode
}

func (translator *HackTranslator) TranslateCInstruction(comp, dest, jmp string) string {
	code := "111" + comp + dest + jmp
	return code
}

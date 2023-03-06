package main

import (
	"testing"
)

func TestSymbolTable_Initialize(t *testing.T) {
	var st SymbolTable
	st.Initialize()
	if st.nextFreeMem != 16 {
		t.Errorf("Result incorrect: got %d, expected %d\n", st.nextFreeMem, 16)
	}
}

func TestSymbolTable_StoreLabel(t *testing.T) {
	var st SymbolTable
	st.Initialize()
	toStore := st.StoreLabel("(myLabel)", 10)
	if toStore != true {
		t.Errorf("Expected true, got %v\n", toStore)
	}
	toSkip := st.StoreLabel("(myLabel)", 10)
	if toSkip != false {
		t.Errorf("Expected false, got %v\n", toSkip)
	}
}

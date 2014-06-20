package main

import (
	"testing"
)

func TestNewVar(t *testing.T) {
	if err := CreateVariable(0, "testname", "testdata"); err != nil {
		t.Fatal("can't create var", err)
	}

	if value, err := ReadVariable(0, "testname"); err != nil {
		t.Fatal("Can't read var", err)
	} else {
		if value != "testdata" {
			t.Fatal("value =", value, " should =  testdata")
		}
	}
}

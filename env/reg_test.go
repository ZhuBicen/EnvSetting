package env

import (
	"testing"
)

func TestNewVar(t *testing.T) {
	DeleteVariable(0, "NotExisting")
	if _, err := ReadVariable(0, "NotExisting"); err == nil {
		t.Fatal("Can't read var NotExisting.")
	}
	if err := EditVariable(0, "testname", "testdata"); err != nil {
		t.Fatal("can't create var", err)
	}

	if value, err := ReadVariable(0, "testname"); err != nil {
		t.Fatal("Can't read var", err)
	} else {
		if value != "testdata" {
			t.Fatal("value =", value, " should =  testdata")
		}
	}

	if err := DeleteVariable(0, "testname"); err != nil {
		t.Fatal("Can't delete the var", err)
	}
}

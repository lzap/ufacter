package ufacter

import (
	"strings"
	"testing"
)

func TestSimpleKey(t *testing.T) {
	f := NewJSONFormatter()
	if f == nil {
		t.Fail()
	}
	f.Add(NewFact(1, false, "key"))
	json := f.JSONString()
	if strings.Compare(json, `{"key":1}`) != 0 {
		t.Logf("Returned: %v", json)
		t.Fail()
	}
}

func TestOneLevel(t *testing.T) {
	f := NewJSONFormatter()
	if f == nil {
		t.Fail()
	}
	f.Add(NewFact(1, false, "node", "leaf"))
	json := f.JSONString()
	if strings.Compare(json, `{"node":{"leaf":1}}`) != 0 {
		t.Logf("Returned: %v", json)
		t.Fail()
	}
}

func TestTwoLevels(t *testing.T) {
	f := NewJSONFormatter()
	if f == nil {
		t.Fail()
	}
	f.Add(NewFact(1, false, "key"))
	f.Add(NewFact(1, false, "node", "node", "leaf"))
	json := f.JSONString()
	if strings.Compare(json, `{"key":1,"node":{"node":{"leaf":1}}}`) != 0 {
		t.Logf("Returned: %v", json)
		t.Fail()
	}
}

func TestNodeOverwrite(t *testing.T) {
	f := NewJSONFormatter()
	if f == nil {
		t.Fail()
	}
	f.Add(NewFact(1, false, "node"))
	f.Add(NewFact(1, false, "node", "leaf"))
	json := f.JSONString()
	if strings.Compare(json, `{"node":{"leaf":1}}`) != 0 {
		t.Logf("Returned: %v", json)
		t.Fail()
	}
}

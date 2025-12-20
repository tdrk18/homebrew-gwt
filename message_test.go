package main

import (
	"fmt"
	"testing"
)

func TestErrorMsgSetsError(t *testing.T) {
	m := Model{}

	err := errorMsg{Err: fmt.Errorf("boom")}
	m2, _ := m.Update(err)
	m = m2.(Model)

	if m.Error != "boom" {
		t.Fatalf("error = %q", m.Error)
	}
}

func TestSuccessMsgQuits(t *testing.T) {
	m := Model{}
	_, cmd := m.Update(successMsg{})

	if cmd == nil {
		t.Fatalf("expected quit cmd")
	}
}

package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCursorMoveDownUp(t *testing.T) {
	m := Model{
		Contexts: []Context{
			{Path: "/a"},
			{Path: "/b"},
			{Path: "/c"},
		},
		Cursor: 0,
	}

	// j
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = m2.(Model)

	if m.Cursor != 1 {
		t.Fatalf("cursor = %d, want 1", m.Cursor)
	}

	// k
	m2, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = m2.(Model)

	if m.Cursor != 0 {
		t.Fatalf("cursor = %d, want 0", m.Cursor)
	}
}

func TestEnterSelectsPath(t *testing.T) {
	m := Model{
		Contexts: []Context{
			{Path: "/a"},
			{Path: "/b"},
		},
		Cursor: 1,
	}

	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = m2.(Model)

	if m.SelectedPath != "/b" {
		t.Fatalf("selected = %q, want /b", m.SelectedPath)
	}
}

func TestNewBranchInputFlow(t *testing.T) {
	m := Model{
		Contexts: []Context{{Path: "/a"}},
	}

	// n → input mode
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = m2.(Model)

	if m.InputMode != InputNewBranch {
		t.Fatalf("mode = %v, want InputNewBranch", m.InputMode)
	}

	// input "foo"
	m2, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")})
	m = m2.(Model)
	m2, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("o")})
	m = m2.(Model)

	if m.InputText != "fo" {
		t.Fatalf("input = %q, want fo", m.InputText)
	}

	// enter → command発行（cmdの中身までは見ない）
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("expected addWorktreeCmd, got nil")
	}
}

func TestConfirmDeleteCancel(t *testing.T) {
	m := Model{
		Contexts: []Context{
			{Path: "/a"},
			{Path: "/b"},
		},
		Cursor:      1,
		InputMode:   InputConfirmDelete,
		DeleteIndex: 1,
	}

	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = m2.(Model)

	if m.InputMode != InputNone {
		t.Fatalf("mode = %v, want InputNone", m.InputMode)
	}
}

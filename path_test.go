package main

import (
	"path/filepath"
	"testing"
)

func TestRelativePath(t *testing.T) {
	root := "/repo"
	path := "/repo/worktrees/foo"

	got := relativePath(root, path)
	want := filepath.Join("worktrees", "foo")

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestInitialCursor(t *testing.T) {
	contexts := []Context{
		{Path: "/a", IsCurrent: false},
		{Path: "/b", IsCurrent: true},
		{Path: "/c", IsCurrent: false},
	}

	if got := initialCursor(contexts); got != 1 {
		t.Fatalf("cursor = %d, want 1", got)
	}
}

func TestInitialCursor_DefaultZero(t *testing.T) {
	contexts := []Context{
		{Path: "/a"},
		{Path: "/b"},
	}

	if got := initialCursor(contexts); got != 0 {
		t.Fatalf("cursor = %d, want 0", got)
	}
}

package main

import (
	"os"
	"path/filepath"
)

func relativePath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path // Fallback to original path if relative path calculation fails
	}
	return rel
}

func samePath(a, b string) bool {
	ap, err1 := filepath.EvalSymlinks(a)
	bp, err2 := filepath.EvalSymlinks(b)

	if err1 != nil || err2 != nil {
		return false
	}

	return ap == bp
}

func getInitialCWD() (string, error) {
	if cwd := os.Getenv("WT_CWD"); cwd != "" {
		return cwd, nil
	}
	return os.Getwd()
}

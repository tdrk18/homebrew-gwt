package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	repoRoot, err := detectRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	_ = repoRoot
}

func detectRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or git not installed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

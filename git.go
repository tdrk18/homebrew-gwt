package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func detectRepoRoot() (string, error) {
	cmd := execCommand("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or git not installed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func listWorktrees() ([]Worktree, error) {
	cmd := execCommand("git", "worktree", "list", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")

	var result []Worktree
	var current *Worktree

	for _, line := range lines {
		if line == "" {
			if current != nil {
				result = append(result, *current)
				current = nil
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			if current != nil {
				result = append(result, *current)
			}
			current = &Worktree{
				Path: strings.TrimPrefix(line, "worktree "),
			}
			continue
		}

		if current == nil {
			continue
		}

		if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch refs/heads/")
			continue
		}

		if line == "detached" {
			current.IsDetached = true
			continue
		}
	}

	if current != nil {
		result = append(result, *current)
	}

	return result, nil
}

func isWorktreeDirty(path string) bool {
	cmd := execCommand("git", "-C", path, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}

func addWorktree(branch string) error {
	_ = os.MkdirAll(".worktree", 0755)

	wtPath := filepath.Join(".worktree", branch)

	args := []string{"worktree", "add"}
	if !branchExists(branch) {
		// Create a new branch from HEAD at the given worktree path:
		// git worktree add -b <branch> <path>
		args = append(args, "-b", branch, wtPath)
	} else {
		// Add a worktree for an existing branch:
		// git worktree add <path> <branch>
		args = append(args, wtPath, branch)
	}

	cmd := execCommand("git", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func removeWorktree(ctx Context) error {
	var stderr bytes.Buffer

	cmd1 := execCommand("git", "worktree", "remove", ctx.Path)
	cmd1.Stderr = &stderr
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}

	if ctx.Branch != "" && !ctx.IsDetached {
		stderr.Reset()
		cmd2 := execCommand("git", "branch", "-D", ctx.Branch)
		cmd2.Stderr = &stderr
		if err := cmd2.Run(); err != nil {
			return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
		}
	}

	return nil
}

func branchExists(branch string) bool {
	cmd := execCommand(
		"git",
		"show-ref",
		"--verify",
		"--quiet",
		"refs/heads/"+branch,
	)
	return cmd.Run() == nil
}

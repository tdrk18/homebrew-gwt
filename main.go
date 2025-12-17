package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	_, err := detectRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	worktrees, err := listWorktrees()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cwd, err := getInitialCWD()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	contexts := buildContexts(worktrees, cwd)

	// 確認用（あとで消す）
	for _, ctx := range contexts {
		fmt.Fprintf(os.Stderr, "%+v\n", ctx)
	}
}

func detectRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or git not installed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func listWorktrees() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
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

func getInitialCWD() (string, error) {
	if cwd := os.Getenv("WT_CWD"); cwd != "" {
		return cwd, nil
	}
	return os.Getwd()
}

func buildContexts(
	worktrees []Worktree,
	currentCWD string,
) []Context {
	var contexts []Context

	for _, wt := range worktrees {
		isCurrent := samePath(wt.Path, currentCWD)

		contexts = append(contexts, Context{
			Path:       wt.Path,
			Branch:     wt.Branch,
			IsDetached: wt.IsDetached,
			IsCurrent:  isCurrent,
		})
	}

	return contexts
}

func samePath(a, b string) bool {
	ap, err1 := filepath.EvalSymlinks(a)
	bp, err2 := filepath.EvalSymlinks(b)

	if err1 != nil || err2 != nil {
		return false
	}

	return ap == bp
}

type Worktree struct {
	Path       string
	Branch     string // empty if detached
	IsDetached bool
}

type Context struct {
	Path       string
	Branch     string
	IsDetached bool
	IsCurrent  bool
}

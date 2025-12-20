package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// テスト終了時に必ず戻す
func withFakeExec(t *testing.T, fn func(name string, args ...string) *exec.Cmd) {
	t.Helper()

	orig := execCommand
	execCommand = fn

	t.Cleanup(func() {
		execCommand = orig
	})
}

// helper process
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	cmd := strings.Join(args, " ")

	switch {
	case strings.Contains(cmd, "rev-parse --show-toplevel"):
		os.Stdout.WriteString("/repo\n")
		os.Exit(0)

	case strings.Contains(cmd, "worktree list"):
		os.Stdout.WriteString(`
worktree /repo
branch refs/heads/main

worktree /repo/.worktree/feature
branch refs/heads/feature

worktree /repo/.worktree/detached
detached
`)
		os.Exit(0)

	case strings.Contains(cmd, "status --porcelain"):
		// パスで挙動を切り替える
		switch {
		case strings.Contains(cmd, "/clean"):
			// clean repo
			os.Exit(0)

		case strings.Contains(cmd, "/dirty"):
			os.Stdout.WriteString(" M file.txt\n")
			os.Exit(0)

		default:
			// git error
			os.Exit(1)
		}

	case strings.Contains(cmd, "worktree add"):
		switch {
		case strings.Contains(cmd, "success"):
			os.Exit(0)
		case strings.Contains(cmd, "failure"):
			os.Stderr.WriteString("fatal: a branch named 'foo' already exists\n")
			os.Exit(1)
		}

	case strings.Contains(cmd, "worktree remove"):
		switch {
		case strings.Contains(cmd, "/fail-remove"):
			os.Stderr.WriteString("fatal: cannot remove worktree\n")
			os.Exit(1)
		default:
			os.Exit(0)
		}

	case strings.Contains(cmd, "branch -D"):
		switch {
		case strings.Contains(cmd, "fail-branch"):
			os.Stderr.WriteString("error: failed to delete branch\n")
			os.Exit(1)
		default:
			os.Exit(0)
		}

	case strings.Contains(cmd, "show-ref --verify"):
		// branch exists
		os.Exit(0)
	}

	os.Exit(1)
}

func fakeExecCommand(name string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, args...)

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")

	return cmd
}

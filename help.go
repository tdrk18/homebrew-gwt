package main

import "fmt"

func printHelp() {
	fmt.Println(`gwt-bin is an internal command for gwt.

This command launches a TUI to select a git worktree and prints
the selected path to stdout.

Normally, you should NOT run this command directly.
Instead, use the shell function:

  gwt

Setup example (zsh):

  source ./shell/gwt.zsh`)
}

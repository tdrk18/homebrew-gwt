package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	var (
		showHelp    bool
		showVersion bool
	)

	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.BoolVar(&showHelp, "h", false, "show help (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "show version")

	// flag.Parse() は TUI より前！
	flag.Parse()

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("%s %s\n", AppName, AppVersion)
		os.Exit(0)
	}

	// TUI
	repoRoot, err := detectRepoRoot()
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

	contexts := buildContexts(worktrees, cwd, repoRoot)
	model := newModel(contexts)

	p := tea.NewProgram(
		model,
		tea.WithOutput(os.Stderr),
	)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	m := finalModel.(Model)
	if m.SelectedPath != "" {
		fmt.Print(m.SelectedPath)
	}
}

func buildContexts(
	worktrees []Worktree,
	currentCWD string,
	repoRoot string,
) []Context {
	var contexts []Context

	for _, wt := range worktrees {
		isCurrent := samePath(wt.Path, currentCWD)

		contexts = append(contexts, Context{
			Path:        wt.Path,
			DisplayPath: relativePath(repoRoot, wt.Path),
			Branch:      wt.Branch,
			IsDetached:  wt.IsDetached,
			IsCurrent:   isCurrent,
			IsDirty:     isWorktreeDirty(wt.Path),
		})
	}

	return contexts
}

func initialCursor(contexts []Context) int {
	for i, ctx := range contexts {
		if ctx.IsCurrent {
			return i
		}
	}
	return 0
}

func refreshCmd(contexts []Context) tea.Cmd {
	return func() tea.Msg {
		for i := range contexts {
			contexts[i].IsDirty = isWorktreeDirty(contexts[i].Path)
		}
		return contexts
	}
}

func addWorktreeCmd(branch string) tea.Cmd {
	return func() tea.Msg {
		if err := addWorktree(branch); err != nil {
			return errorMsg{Err: err}
		}
		return successMsg{}
	}
}

func removeWorktreeCmd(ctx Context) tea.Cmd {
	return func() tea.Msg {
		if err := removeWorktree(ctx); err != nil {
			return errorMsg{Err: err}
		}
		return successMsg{}
	}
}

type Worktree struct {
	Path       string
	Branch     string // empty if detached
	IsDetached bool
}

type Context struct {
	Path        string
	DisplayPath string

	Branch     string
	IsDetached bool
	IsCurrent  bool
	IsDirty    bool
}

type InputMode int

const (
	InputNone InputMode = iota
	InputNewBranch
	InputConfirmDelete
)

type successMsg struct{}

type errorMsg struct {
	Err error
}

var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	currentStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	dirtyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	detachedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
	model := newModel(contexts)

	p := tea.NewProgram(model)
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

func initialCursor(contexts []Context) int {
	for i, ctx := range contexts {
		if ctx.IsCurrent {
			return i
		}
	}
	return 0
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

type Model struct {
	Contexts     []Context
	Cursor       int
	SelectedPath string
}

func newModel(contexts []Context) Model {
	return Model{
		Contexts: contexts,
		Cursor:   initialCursor(contexts),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if len(m.Contexts) > 0 {
				m.SelectedPath = m.Contexts[m.Cursor].Path
			}
			return m, tea.Quit

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "j":
				if m.Cursor < len(m.Contexts)-1 {
					m.Cursor++
				}
			case "k":
				if m.Cursor > 0 {
					m.Cursor--
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	for i, ctx := range m.Contexts {
		prefix := "  "
		if ctx.IsCurrent {
			prefix = "> "
		}

		line := fmt.Sprintf("%s%s", prefix, ctx.Path)
		if i == m.Cursor {
			line = "[" + line + "]"
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

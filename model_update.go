package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// --- Input mode: new branch ---
	if m.InputMode == InputNewBranch {
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.Type {

			case tea.KeyEsc:
				m.InputMode = InputNone
				m.InputText = ""
				return m, nil

			case tea.KeyEnter:
				branch := strings.TrimSpace(m.InputText)
				m.InputMode = InputNone
				m.InputText = ""

				if branch != "" {
					return m, addWorktreeCmd(branch)
				}
				return m, nil

			case tea.KeyBackspace:
				if len(m.InputText) > 0 {
					m.InputText = m.InputText[:len(m.InputText)-1]
				}
				return m, nil

			case tea.KeyRunes:
				m.InputText += string(msg.Runes)
				return m, nil
			}
		}
		return m, nil
	}

	// --- Input mode: confirm delete ---
	if m.InputMode == InputConfirmDelete {
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch string(msg.Runes) {

			case "y":
				ctx := m.Contexts[m.DeleteIndex]
				m.InputMode = InputNone
				return m, removeWorktreeCmd(ctx)

			case "n":
				m.InputMode = InputNone
				return m, nil
			}

			if msg.Type == tea.KeyEsc {
				m.InputMode = InputNone
				return m, nil
			}
		}
		return m, nil
	}

	// --- Normal mode ---
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEsc:
			// clear error if exists
			if m.Error != "" {
				m.Error = ""
				return m, nil
			}
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

			case "r":
				return m, refreshCmd(m.Contexts)

			case "n":
				m.InputMode = InputNewBranch
				m.InputText = ""

			case "d":
				ctx := m.Contexts[m.Cursor]
				if ctx.IsCurrent {
					return m, nil
				}
				m.InputMode = InputConfirmDelete
				m.DeleteIndex = m.Cursor
			}
		}

	case []Context:
		m.Contexts = msg
		return m, nil

	case successMsg:
		return m, tea.Quit

	case errorMsg:
		m.Error = msg.Err.Error()
		return m, nil
	}

	return m, nil
}

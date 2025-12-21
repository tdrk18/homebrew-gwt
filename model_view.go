package main

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(
		"j/k: move  enter: cd  n: new  r: refresh  esc: quit   * current  ! dirty  @ detached\n",
	)

	if m.Error != "" {
		b.WriteString(fmt.Sprintf("error: %s\n", m.Error))
	} else if m.InputMode == InputNewBranch {
		b.WriteString(fmt.Sprintf("new branch: %s\n", m.InputText))
	} else if m.InputMode == InputConfirmDelete {
		ctx := m.Contexts[m.DeleteIndex]
		b.WriteString(
			fmt.Sprintf("delete %s (%s)? y/n\n", ctx.Path, ctx.Branch),
		)
	} else {
		b.WriteString("\n")
	}

	for i, ctx := range m.Contexts {
		// cursor
		cursor := "  "
		if m.Cursor == i {
			cursor = cursorStyle.Render("> ")
		}

		// current marker
		current := "  "
		if ctx.IsCurrent {
			current = "* "
		}

		status := "  "
		if ctx.IsDetached && ctx.IsDirty {
			status = detachedStyle.Render("@!")
		} else if ctx.IsDetached {
			status = detachedStyle.Render("@ ")
		} else if ctx.IsDirty {
			status = dirtyStyle.Render("! ")
		}

		branch := ctx.Branch
		if ctx.IsDetached {
			branch = "(detached)"
		}

		line := fmt.Sprintf(
			"%s%s%-12s %s %s",
			cursor,
			current,
			branch,
			status,
			ctx.DisplayPath,
		)

		if ctx.IsCurrent {
			line = currentStyle.Render(line)
		}

		b.WriteString(line + "\n")
	}

	return b.String()
}

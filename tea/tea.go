package tea

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	branches []string
	cursor   int
	selected map[int]struct{}
}

func initialModel(branches []string) model {
	return model{
		branches: branches,
		selected: make(map[int]struct{}),
	}
}

func TeaMe(branches []string) (tea.Model, error) {
	p := tea.NewProgram(initialModel(branches))
	t, err := p.Run()
	return t, err
}

// tea required Init(), but we don't do anything with it.
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.branches)-1 {
				m.cursor++
			}

		// The spacebar (a literal space) toggle the selected state for
		// the item that the cursor is pointing at.
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		// The enter key deletes the selected branches
		case "enter":
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Select branches to be removed:\n\n"

	// Iterate over our branches
	for i, choice := range m.branches {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

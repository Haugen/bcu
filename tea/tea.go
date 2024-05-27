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

func GetBranches(branches []string) []string {
	p := tea.NewProgram(initialModel(branches))
	t, err := p.Run()

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	var selectedBranches []string
	for k := range t.(model).selected {
		selectedBranches = append(selectedBranches, t.(model).branches[k])
	}

	return selectedBranches
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

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

		// The enter key quits the branch selection.
		case "enter":
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := "Select branches to be removed:\n\n"

	for i, choice := range m.branches {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

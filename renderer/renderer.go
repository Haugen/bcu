package renderer

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

const (
	qKey      = byte(113)
	ctrlC     = byte(3)
	upArrow   = byte(38)
	jKey      = byte(106)
	downArrow = byte(40)
	kKey      = byte(107)
	spaceKey  = byte(32)
	enterKey  = byte(13)
)

// Renderer manages the interactive branch selection UI
type Renderer struct {
	state       *State
	firstRender bool
}

// State holds the current state of the branch selector
type State struct {
	cursorPos int
	list      []string
	selected  map[int]bool
}

// NewState creates a new state with the given branch list
func NewState(branches []string) *State {
	return &State{
		cursorPos: 0,
		list:      branches,
		selected:  make(map[int]bool, len(branches)),
	}
}

// MoveCursorUp moves the cursor up by one position if possible
func (s *State) MoveCursorUp() bool {
	if s.cursorPos > 0 {
		s.cursorPos--
		return true
	}
	return false
}

// MoveCursorDown moves the cursor down by one position if possible
func (s *State) MoveCursorDown() bool {
	if s.cursorPos < len(s.list)-1 {
		s.cursorPos++
		return true
	}
	return false
}

// ToggleSelection toggles the selection state of the current item
func (s *State) ToggleSelection() {
	s.selected[s.cursorPos] = !s.selected[s.cursorPos]
}

// GetSelectedBranches returns a slice of selected branch names
func (s *State) GetSelectedBranches() []string {
	var branches []string
	for i, branch := range s.list {
		if s.selected[i] {
			branches = append(branches, branch)
		}
	}
	return branches
}

// GetOutputLines returns the lines to be rendered (without ANSI codes)
func (s *State) GetOutputLines() []string {
	lines := []string{"Select branches to delete (use ↑/↓ or j/k to navigate, Space to select, Enter to confirm, q to quit):", ""}

	for i, item := range s.list {
		cursor := "  "
		if s.cursorPos == i {
			cursor = "> "
		}

		checkbox := "[ ] "
		if s.selected[i] {
			checkbox = "[x] "
		}

		lines = append(lines, cursor+checkbox+item)
	}

	return lines
}

func readInput() (byte, error) {
	var char byte
	buf := make([]byte, 3)
	_, err := os.Stdin.Read(buf)
	char = buf[0]

	// Map Up and Down arrow keys to the Javascript key codes since they don't fit in a single byte
	if buf[0] == 27 && buf[1] == 91 {
		switch buf[2] {
		case 65:
			char = 38
		case 66:
			char = 40
		}
	}

	return char, err
}

func NewRenderer(branches []string) *Renderer {
	return &Renderer{
		state:       NewState(branches),
		firstRender: true,
	}
}

func (r *Renderer) Run() []string {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Initial render
	r.render()

	for {
		input, err := readInput()
		if err != nil {
			log.Fatal(err)
		}

		switch input {
		case upArrow, kKey:
			if r.state.MoveCursorUp() {
				r.render()
			}
		case downArrow, jKey:
			if r.state.MoveCursorDown() {
				r.render()
			}
		case spaceKey:
			r.state.ToggleSelection()
			r.render()
		case enterKey:
			return r.state.GetSelectedBranches()
		case qKey, ctrlC:
			return []string{}
		}
	}
}

func (r *Renderer) render() {
	// Move cursor up to overwrite previous output (skip on first render)
	if !r.firstRender {
		fmt.Printf("\033[%dA", len(r.state.list)+2)
	}
	r.firstRender = false

	// Clear from cursor down
	fmt.Print("\033[J")

	// Get output lines and print them with proper line endings
	lines := r.state.GetOutputLines()
	for _, line := range lines {
		fmt.Print(line + "\r\n")
	}
}

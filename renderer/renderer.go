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

type Renderer struct {
	cursorPos   int
	list        []string
	selected    map[int]bool
	firstRender bool
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
		cursorPos:   0,
		list:        branches,
		selected:    make(map[int]bool, len(branches)),
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
			if r.cursorPos > 0 {
				r.cursorPos--
				r.render()
			}
		case downArrow, jKey:
			if r.cursorPos < len(r.list)-1 {
				r.cursorPos++
				r.render()
			}
		case spaceKey:
			r.selected[r.cursorPos] = !r.selected[r.cursorPos]
			r.render()
		case enterKey:
			// User confirmed selection
			return r.getSelectedBranches()
		case qKey, ctrlC:
			// User quit without confirming
			return []string{}
		}
	}
}

func (r *Renderer) getSelectedBranches() []string {
	var branches []string
	for i, branch := range r.list {
		if r.selected[i] {
			branches = append(branches, branch)
		}
	}
	return branches
}

func (r *Renderer) render() {
	// Move cursor up to overwrite previous output (skip on first render)
	if !r.firstRender {
		fmt.Printf("\033[%dA", len(r.list)+2)
	}
	r.firstRender = false

	// Clear from cursor down
	fmt.Print("\033[J")

	fmt.Print("Select branches to delete (use ↑/↓ or j/k to navigate, Space to select, Enter to confirm, q to quit):\r\n\r\n")

	for i, item := range r.list {
		cursor := "  "
		if r.cursorPos == i {
			cursor = "> "
		}

		checkbox := "[ ] "
		if r.selected[i] {
			checkbox = "[x] "
		}

		fmt.Print(cursor + checkbox + item + "\r\n")
	}
}

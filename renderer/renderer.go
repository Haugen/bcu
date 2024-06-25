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
	cursorPos int
	list      []string
	selected  map[int]bool
}

func readInput() (byte, error) {
	var char byte
	buf := make([]byte, 3)
	_, err := os.Stdin.Read(buf)
	char = buf[0]

	// Map Up and Down arrow keys to the Javascript key codes since they don't fit in a single byte
	if buf[0] == 27 && buf[1] == 91 {
		if buf[2] == 65 {
			char = 38
		} else if buf[2] == 66 {
			char = 40
		}
	}

	return char, err
}

func NewRenderer(branches []string) *Renderer {
	return &Renderer{
		cursorPos: 0,
		list:      branches,
		selected:  make(map[int]bool, len(branches)),
	}
}

func (r *Renderer) Run() map[int]bool {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		r.render()
		input, err := readInput()
		if err != nil {
			log.Fatal(err)
		}

		switch input {
		case upArrow:
		case kKey:
			if r.cursorPos > 0 {
				r.cursorPos--
			}
		case downArrow:
		case jKey:
			if r.cursorPos < len(r.list)-1 {
				r.cursorPos++
			}
		case spaceKey:
			r.selected[r.cursorPos] = !r.selected[r.cursorPos]
		}

		if input == qKey || input == ctrlC {
			break
		}
	}

	return r.selected
}

func (r *Renderer) render() {
	fmt.Print("\033[H\033[2J")
	for i, item := range r.list {
		if r.cursorPos == i {
			fmt.Print("> ")
		} else {
			fmt.Print("  ")
		}
		if r.selected[i] {
			fmt.Print("[x] ")
		} else {
			fmt.Print("[ ] ")
		}
		fmt.Println(item)
	}
}

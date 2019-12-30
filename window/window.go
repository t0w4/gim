package window

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	prompt "github.com/c-bata/go-prompt"

	"golang.org/x/crypto/ssh/terminal"
)

type Size struct {
	Row    int
	Column int
}

var asciiSequences = []*prompt.ASCIICode{
	{Key: prompt.Escape, ASCIICode: []byte{0x1b}},
	{Key: prompt.Up, ASCIICode: []byte{0x1b, 0x5b, 0x41}},
	{Key: prompt.Down, ASCIICode: []byte{0x1b, 0x5b, 0x42}},
	{Key: prompt.Right, ASCIICode: []byte{0x1b, 0x5b, 0x43}},
	{Key: prompt.Left, ASCIICode: []byte{0x1b, 0x5b, 0x44}},
	{Key: prompt.ControlC, ASCIICode: []byte{0x3}},
	{Key: prompt.Delete, ASCIICode: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{Key: prompt.Backspace, ASCIICode: []byte{0x7f}},
	{Key: prompt.Enter, ASCIICode: []byte{0xd}},
}

const (
	normalMode = iota
	insertMode
	commandMode
)

type Window struct {
	Size
	Input        *os.File  // Adopts os.File to use Fd () , ex) Stdin
	Output       io.Writer // ex) Stdout
	FileContents [][]byte
	position     Position
	mode         int // ex) insert mode
	command      []byte
}

func NewWindow(input *os.File, output io.Writer) *Window {
	return &Window{
		Input:        input,
		Output:       output,
		FileContents: nil,
		position:     Position{X: 1, Y: 1},
		mode:         normalMode,
		command:      []byte{},
	}
}

func (w *Window) IsInsertMode() bool {
	if w.mode == insertMode {
		return true
	}
	return false
}

func (w *Window) IsNormalMode() bool {
	if w.mode == normalMode {
		return true
	}
	return false
}

func (w *Window) IsCommandMode() bool {
	if w.mode == commandMode {
		return true
	}
	return false
}

func (w *Window) IsCommandNotTyped() bool {
	if len(w.command) == 0 {
		return true
	}
	return false
}

func (w *Window) SetInsertMode() {
	w.mode = insertMode
}

func (w *Window) SetNormalMode() {
	w.mode = normalMode
}

func (w *Window) SetCommandMode() {
	w.mode = commandMode
}

func (w *Window) AddCommand(b []byte) {
	w.command = append(w.command, b...)
}

func (w *Window) RemoveCommand() {
	if len(w.command) > 0 {
		w.command = w.command[:len(w.command)-1]
	}
}

func (w *Window) ResetCommand() {
	w.command = []byte{}
}

func (w *Window) TypedCommand() string {
	return string(w.command)
}

func (w *Window) ExecuteCommand() {
}

func (w *Window) InputtedUp() {
	// if cursor is top, don't move
	if w.position.Y == 1 {
		return
	}
	// If the number of characters in the line above is smaller than the current X,
	// the cursor moves to the last column
	if len(w.FileContents[w.position.Y-2]) < w.position.X {
		if len(w.FileContents[w.position.Y-2]) == 0 {
			w.position.X = 1
		} else {
			var limitX int
			if w.IsInsertMode() {
				limitX = len(w.FileContents[w.position.Y-2]) + 1
			} else {
				limitX = len(w.FileContents[w.position.Y-2])
			}
			w.position.X = limitX
		}
	}
	w.position.MoveUp(1)
	fmt.Fprintf(w.Output, "\033[%d;%dH> X: %d, Y: %d, Up    ", w.Row, 0, w.position.X, w.position.Y)
	w.MoveCursorToCurrentPosition()
}

func (w *Window) InputtedDown() {
	if len(w.FileContents) == w.position.Y {
		return
	}
	if len(w.FileContents[w.position.Y]) < w.position.X {
		if len(w.FileContents[w.position.Y]) == 0 {
			w.position.X = 1
		} else {
			var limitX int
			if w.IsInsertMode() {
				limitX = len(w.FileContents[w.position.Y]) + 1
			} else {
				limitX = len(w.FileContents[w.position.Y])
			}
			w.position.X = limitX
		}
	}
	w.position.MoveDown(1)
	fmt.Fprintf(w.Output, "\033[%d;%dH> X: %d, Y: %d, Down  ", w.Row, 0, w.position.X, w.position.Y)
	w.MoveCursorToCurrentPosition()
}

func (w *Window) InputtedLeft() {
	w.position.MoveLeft(1)
	fmt.Fprintf(w.Output, "\033[%d;%dH> X: %d, Y: %d, Left  ", w.Row, 0, w.position.X, w.position.Y)
	w.MoveCursorToCurrentPosition()
}

func (w *Window) InputtedRight() {
	var limitX int
	if w.IsInsertMode() {
		limitX = len(w.FileContents[w.position.Y-1]) + 1
	} else {
		limitX = len(w.FileContents[w.position.Y-1])
	}
	if limitX <= w.position.X {
		w.MoveCursorToCurrentPosition()
		return
	}
	w.position.MoveRight(1)
	fmt.Fprintf(w.Output, "\033[%d;%dH> X: %d, Y: %d, Right", w.Row, 0, w.position.X, w.position.Y)
	w.MoveCursorToCurrentPosition()
}

func (w *Window) InputtedOther(b []byte) {
	switch w.mode {
	case normalMode:
		if string(b) == "i" {
			w.SetInsertMode()
			return
		}
		if string(b) == ":" {
			w.SetCommandMode()
			fmt.Fprintf(w.Output, "\033[%d;%dH:", w.Size.Row, 0)
			return
		}
		fmt.Fprintf(w.Output, "\033[%d;%dH> X: %d, Y: %d, input: %s     ", w.Row, 0, w.position.X, w.position.Y, string(b))
		w.MoveCursorToCurrentPosition()
	case insertMode:
		be := w.FileContents[w.position.Y-1][:w.position.X-1]
		af := w.FileContents[w.position.Y-1][w.position.X-1:]
		w.FileContents[w.position.Y-1] = []byte(string(be) + string(b) + string(af))
		w.position.MoveRight(1)
		fmt.Fprintf(w.Output, "\033[%d;%dH%s", w.position.Y, 0, string(w.FileContents[w.position.Y-1]))
		w.MoveCursorToCurrentPosition()
	}
}

func (w *Window) GetKey(b []byte) prompt.Key {
	for _, k := range asciiSequences {
		if bytes.Equal(k.ASCIICode, b) {
			return k.Key
		}
	}
	return prompt.NotDefined
}

func (w *Window) SetSize() error {
	var err error
	w.Column, w.Row, err = terminal.GetSize(int(w.Input.Fd()))
	if err != nil {
		return err
	}
	return nil
}

// PrintFileContents outputs the contents passed by fc (usually the contents of the read file)
// to the location specified by Output.
// The file contents are not printed  on the last line.
func (w *Window) PrintFileContents() {
	fmt.Fprint(w.Output, "\033[H\033[2J")
	for i := 0; i < w.Row-1; i++ {
		if len(w.FileContents) <= i {
			fmt.Fprintln(w.Output, "")
		} else {
			fmt.Fprintf(w.Output, "%s\n", w.FileContents[i])
		}
	}
	w.MoveCursorToCurrentPosition()
}

func (w *Window) MoveCursorToCurrentPosition() {
	fmt.Fprintf(w.Output, "\033[%d;%dH", w.position.Y, w.position.X)
}

func (w *Window) ReadBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)

	reader := bufio.NewReader(w.Input)
	for {
		if n, err := reader.Read(buf); err == nil {
			bufCh <- buf[:n]
		}
	}
}

func (w *Window) SetFileContents(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		w.FileContents = append(w.FileContents, []byte(sc.Text()))
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

type Position struct {
	X int
	Y int
}

func (p *Position) MoveDown(num int) {
	p.Y += num
}

func (p *Position) MoveUp(num int) {
	if p.Y == 1 {
		return
	}
	p.Y -= num
}

func (p *Position) MoveRight(num int) {
	p.X += num
}

func (p *Position) MoveLeft(num int) {
	if p.X == 1 {
		return
	}
	p.X -= num
}

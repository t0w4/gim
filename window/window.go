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
}

type Window struct {
	Size
	Input        *os.File  // Adopts os.File to use Fd () , ex) Stdin
	Output       io.Writer // ex) Stdout
	FileContents [][]byte
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
	fmt.Fprint(w.Output, "\033[H")
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

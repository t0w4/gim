package window

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type Size struct {
	Row    int
	Column int
}

type Window struct {
	Size
	Input        io.Reader
	Output       io.Writer
	FileContents [][]byte
}

func (w *Window) SetSize(fd int) error {
	var err error
	w.Column, w.Row, err = terminal.GetSize(fd)
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

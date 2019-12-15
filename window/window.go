package window

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh/terminal"
)

type Size struct {
	Row    int
	Column int
}

type Window struct {
	Size
	Output io.Writer
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
func (w *Window) PrintFileContents(fc [][]byte) {
	fmt.Fprint(w.Output, "\033[H\033[2J")
	for i := 0; i < w.Row-1; i++ {
		if len(fc) <= i {
			fmt.Fprintln(w.Output, "")
		} else {
			fmt.Fprintf(w.Output, "%s\n", fc[i])
		}
	}
	fmt.Fprint(w.Output, "\033[H")
}

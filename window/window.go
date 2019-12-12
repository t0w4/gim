package window

import "golang.org/x/crypto/ssh/terminal"

type Size struct {
	Row    int
	Column int
}

type Window struct {
	Size
}

func GetSize(fd int) (*Size, error) {
	ws := &Window{}
	var err error
	ws.Row, ws.Column, err = terminal.GetSize(fd)
	if err != nil {
		return nil, err
	}
	return &ws.Size, nil
}

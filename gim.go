package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gim/window"
	"io"
	"os"
	"os/signal"
	"syscall"

	prompt "github.com/c-bata/go-prompt"

	"golang.org/x/crypto/ssh/terminal"
)

const NotTerminalWarning = `Gim: Warning: Input is not from a terminal
Gim: Error reading input, exiting...
Gim: Finished.`

const ManyFileEditWarning = `Gim: Warning: Only one file can be edited
Gim: Trying to edit more than one file...
Gim: Finished.`

const (
	ExitOk = iota
	ExitError
)

var fileContents [][]byte
var normalState *terminal.State
var insertMode = false

type position struct {
	X int
	Y int
}

func (p *position) moveDown(num int) {
	p.Y += num
}

func (p *position) moveUp(num int) {
	if p.Y == 0 {
		return
	}
	p.Y -= num
}

func (p *position) moveRight(num int) {
	p.X += num
}

func (p *position) moveLeft(num int) {
	if p.X == 0 {
		return
	}
	p.X -= num
}

func main() {
	if !terminal.IsTerminal(syscall.Stdin) {
		fmt.Println(NotTerminalWarning)
		os.Exit(ExitError)
	}
	switch len(os.Args) {
	case 1:
		fmt.Println("no arg")
	case 2:
		fileName := os.Args[1]
		_, err := os.Stat(fileName)
		if err == os.ErrNotExist {
			fmt.Printf("%s is not exist\n", fileName)
			os.Exit(ExitError)
		} else if err != nil {
			fmt.Printf("file stat error: %v\n", err)
			os.Exit(ExitError)
		}
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("file open error: %v\n", err)
			os.Exit(ExitError)
		}
		defer file.Close()

		rd := bufio.NewReader(file)
		for {
			line, _, err := rd.ReadLine()
			if err == io.EOF {
				break
			}
			fileContents = append(fileContents, line)
		}

		signalChan := make(chan os.Signal, 1)
		// catch SIGINT(Ctrl+C), KILL signal, and window size changes
		signal.Notify(
			signalChan,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGWINCH,
		)
		ws, err := window.GetSize(syscall.Stdin)
		if err != nil {
			fmt.Printf("get window sieze error: %v", err)
			os.Exit(ExitError)
		}
		makeFileWindow(ws.Column)

		exitChan := make(chan int)
		go func() {
			for {
				s := <-signalChan
				switch s {
				// SIGINT(Ctrl+C)
				case syscall.SIGINT:
					exitChan <- 130

				// kILL signal
				case syscall.SIGTERM:
					exitChan <- 143

				case syscall.SIGWINCH:
					// In raw mode, the file content view will be corrupted,
					// so return to normal mode.
					terminal.Restore(syscall.Stdin, normalState)
					ws, err := window.GetSize(syscall.Stdin)
					if err != nil {
						fmt.Printf("get window sieze error: %v", err)
						os.Exit(ExitError)
					}
					makeFileWindow(ws.Column)

				default:
					exitChan <- 1
				}
			}
		}()

		bufCh := make(chan []byte, 128)
		p := position{X: 0, Y: 0}
		go readBuffer(bufCh)
		go func() {
			for {
				normalState, err = terminal.MakeRaw(syscall.Stdin)
				if err != nil {
					fmt.Printf("make raw error: %v\n", err)
					os.Exit(ExitError)
				}
				b := <-bufCh
				switch GetKey(b) {
				case prompt.Up:
					p.moveUp(1)
					fmt.Printf("\033[%d;%dH", p.Y, p.X)
				case prompt.Down:
					p.moveDown(1)
					fmt.Printf("\033[%d;%dH", p.Y, p.X)
				case prompt.Left:
					p.moveLeft(1)
					fmt.Printf("\033[%d;%dH", p.Y, p.X)
				case prompt.Right:
					p.moveRight(1)
					fmt.Printf("\033[%d;%dH", p.Y, p.X)
				case prompt.ControlC:
					exitChan <- 130
				case prompt.NotDefined:
					if string(b) == "i" && !insertMode {
						insertMode = true
						continue
					}
					if insertMode {
						fmt.Print(string(b))
					}
				}
			}
		}()
		code := <-exitChan
		os.Exit(code)

	default:
		fmt.Println(ManyFileEditWarning)
		os.Exit(ExitError)
	}
	os.Exit(ExitOk)
}

func makeFileWindow(column int) {
	fmt.Print("\033[H\033[2J")
	for i := 0; i < column-1; i++ {
		if len(fileContents) <= i {
			fmt.Println("")
		} else {
			fmt.Printf("%s\n", fileContents[i])
		}
	}
	fmt.Print("\033[H")
}

func readBuffer(bufCh chan []byte) {
	buf := make([]byte, 1024)

	reader := bufio.NewReader(os.Stdin)
	for {
		if n, err := reader.Read(buf); err == nil {
			bufCh <- buf[:n]
		}
	}
}

func GetKey(b []byte) prompt.Key {
	for _, k := range asciiSequences {
		if bytes.Equal(k.ASCIICode, b) {
			return k.Key
		}
	}
	return prompt.NotDefined
}

var asciiSequences = []*prompt.ASCIICode{
	{Key: prompt.Escape, ASCIICode: []byte{0x1b}},
	{Key: prompt.Up, ASCIICode: []byte{0x1b, 0x5b, 0x41}},
	{Key: prompt.Down, ASCIICode: []byte{0x1b, 0x5b, 0x42}},
	{Key: prompt.Right, ASCIICode: []byte{0x1b, 0x5b, 0x43}},
	{Key: prompt.Left, ASCIICode: []byte{0x1b, 0x5b, 0x44}},

	{Key: prompt.ControlC, ASCIICode: []byte{0x3}},

	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	{Key: prompt.Up, ASCIICode: []byte{0x1b, 0x4f, 0x41}},
	{Key: prompt.Down, ASCIICode: []byte{0x1b, 0x4f, 0x42}},
	{Key: prompt.Right, ASCIICode: []byte{0x1b, 0x4f, 0x43}},
	{Key: prompt.Left, ASCIICode: []byte{0x1b, 0x4f, 0x44}},
}

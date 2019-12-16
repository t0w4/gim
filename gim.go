package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gim/position"
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

		// create window
		win := window.Window{Output: os.Stdout}

		err = win.SetSize(syscall.Stdin)
		if err != nil {
			fmt.Printf("set window sieze error: %v", err)
			os.Exit(ExitError)
		}
		win.PrintFileContents(fileContents)

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
					err := win.SetSize(syscall.Stdin)
					if err != nil {
						fmt.Printf("set window sieze error: %v", err)
						os.Exit(ExitError)
					}
					win.PrintFileContents(fileContents)
					normalState, err = terminal.MakeRaw(syscall.Stdin)
					if err != nil {
						fmt.Printf("make raw error: %v\n", err)
						os.Exit(ExitError)
					}
				default:
					exitChan <- 1
				}
			}
		}()

		bufCh := make(chan []byte, 128)
		p := position.Position{X: 1, Y: 1}
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
					if p.Y == 1 {
						continue
					}
					if len(fileContents[p.Y-2]) < p.X {
						if len(fileContents[p.Y-2]) == 0 {
							p.X = 1
						} else {
							p.X = len(fileContents[p.Y-2])
						}
					}
					p.MoveUp(1)
					fmt.Printf("\033[%d;%dH> X: %d, Y: %d, Up    ", win.Row, 0, p.X, p.Y)
					fmt.Printf("\033[%d;%dH", p.Y, p.X)
				case prompt.Down:
					if len(fileContents) == p.Y {
						continue
					}
					if len(fileContents[p.Y]) < p.X {
						if len(fileContents[p.Y]) == 0 {
							p.X = 1
						} else {
							p.X = len(fileContents[p.Y])
						}
					}
					p.MoveDown(1)
					fmt.Printf("\033[%d;%dH> X: %d, Y: %d, Down  ", win.Row, 0, p.X, p.Y)
					fmt.Printf("\033[%d;%dH", p.Y, p.X)
				case prompt.Left:
					p.MoveLeft(1)
					fmt.Printf("\033[%d;%dH> X: %d, Y: %d, Left  ", win.Row, 0, p.X, p.Y)
					fmt.Printf("\033[%d;%dH", p.Y, p.X)
				case prompt.Right:
					if len(fileContents[p.Y-1]) <= p.X {
						fmt.Printf("\033[%d;%dH", p.Y, p.X)
						continue
					}
					p.MoveRight(1)
					fmt.Printf("\033[%d;%dH> X: %d, Y: %d, Right", win.Row, 0, p.X, p.Y)
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
					} else {
						fmt.Printf("\033[%d;%dH> X: %d, Y: %d, input: %s     ", win.Row, 0, p.X, p.Y, string(b))
						fmt.Printf("\033[%d;%dH", p.Y, p.X)
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

	//// Tmux sends following keystrokes when control+arrow is pressed, but for
	//// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	//// it a normal arrow press, because that's more important.
	//{Key: prompt.Up, ASCIICode: []byte{0x1b, 0x4f, 0x41}},
	//{Key: prompt.Down, ASCIICode: []byte{0x1b, 0x4f, 0x42}},
	//{Key: prompt.Right, ASCIICode: []byte{0x1b, 0x4f, 0x43}},
	//{Key: prompt.Left, ASCIICode: []byte{0x1b, 0x4f, 0x44}},
}

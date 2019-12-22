package main

import (
	"fmt"
	"gim/window"
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
		signalChan := make(chan os.Signal, 1)
		// catch SIGINT(Ctrl+C), KILL signal, and window size changes
		signal.Notify(
			signalChan,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGWINCH,
		)

		// create window
		win := window.Window{Input: os.Stdin, Output: os.Stdout}

		fileName := os.Args[1]
		if err := win.SetFileContents(fileName); err != nil {
			fmt.Println(err)
			os.Exit(ExitError)
		}

		err := win.SetSize()
		if err != nil {
			fmt.Printf("set window sieze error: %v", err)
			os.Exit(ExitError)
		}
		win.PrintFileContents()

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
					err := win.SetSize()
					if err != nil {
						fmt.Printf("set window sieze error: %v", err)
						os.Exit(ExitError)
					}
					win.PrintFileContents()
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
		p := window.Position{X: 1, Y: 1}
		go win.ReadBuffer(bufCh)
		go func() {
			for {
				normalState, err = terminal.MakeRaw(syscall.Stdin)
				if err != nil {
					fmt.Printf("make raw error: %v\n", err)
					os.Exit(ExitError)
				}
				b := <-bufCh
				switch win.GetKey(b) {
				case prompt.Up:
					// if cursor is top, don't move
					if p.Y == 1 {
						continue
					}
					// If the number of characters in the line above is smaller than the current X,
					// the cursor moves to the last column
					if len(win.FileContents[p.Y-2]) < p.X {
						if len(win.FileContents[p.Y-2]) == 0 {
							p.X = 1
						} else {
							p.X = len(win.FileContents[p.Y-2])
						}
					}
					p.MoveUp(1)
					fmt.Printf("\033[%d;%dH> X: %d, Y: %d, Up    ", win.Row, 0, p.X, p.Y)
					fmt.Printf("\033[%d;%dH", p.Y, p.X)
				case prompt.Down:
					if len(win.FileContents) == p.Y {
						continue
					}
					if len(win.FileContents[p.Y]) < p.X {
						if len(win.FileContents[p.Y]) == 0 {
							p.X = 1
						} else {
							p.X = len(win.FileContents[p.Y])
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
					if len(win.FileContents[p.Y-1]) <= p.X {
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

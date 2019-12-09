package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

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
		ws := GetWindowSize(syscall.Stdin)
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
					ws := GetWindowSize(syscall.Stdin)
					makeFileWindow(ws.Column)

				default:
					exitChan <- 1
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

type Size struct {
	Row    int
	Column int
}

type Window struct {
	Size
}

func GetWindowSize(fd int) *Size {
	ws := &Window{}
	var err error
	ws.Row, ws.Column, err = terminal.GetSize(fd)
	if err != nil {
		fmt.Printf("get window sieze error: %v", err)
		os.Exit(ExitError)
	}
	return &ws.Size
}

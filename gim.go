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

		signalChan := make(chan os.Signal, 1)
		// catch SIGINT(Ctrl+C), KILL signal, and window size changes
		signal.Notify(
			signalChan,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGWINCH,
		)
		ws := GetWindowSize(syscall.Stdin)
		makeFileWindow(fileName, ws.Row, ws.Column)

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
					makeFileWindow(fileName, ws.Row, ws.Column)

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

func makeFileWindow(fileName string, row int, column int) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("file open error: %v\n", err)
		os.Exit(ExitError)
	}
	rd := bufio.NewReader(file)

	fmt.Print("\033[H\033[2J")
	for i := 0; i < column-1; i++ {
		line, _, err := rd.ReadLine()
		if err == io.EOF {
			fmt.Println("")
			continue
		}

		if err != nil {
			fmt.Fprintf(os.Stdout, "read line err : %v\n", err)
			os.Exit(ExitError)
		}
		fmt.Printf("%s\n", string(line))
	}
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

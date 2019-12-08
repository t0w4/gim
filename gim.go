package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

const NotTerminalWarning = `Gim: Warning: Input is not from a terminal
Gim: Error reading input, exiting...
Gim: Finished.`

const (
	ExitOk = iota
	ExitError
)

func main()  {
	if !terminal.IsTerminal(syscall.Stdin) {
		fmt.Println(NotTerminalWarning)
		os.Exit(ExitError)
	}
	fmt.Println("ok")
	os.Exit(ExitOk)
}

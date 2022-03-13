package utils

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func ReadPassword(prompt string) string {
	fmt.Fprint(os.Stderr, prompt)
	var fd int
	if terminal.IsTerminal(syscall.Stdin) {
		fd = syscall.Stdin
	} else {
		tty, err := os.Open("/dev/tty")
		if err != nil {
			return ""
		}
		defer tty.Close()
		fd = int(tty.Fd())
	}

	pass, _ := terminal.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	return string(pass)
}

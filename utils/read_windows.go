package utils

import (
	"fmt"
)

func ReadPassword(prompt string) string {
	var pwd string
	fmt.Printf(prompt)
	fmt.Scanln(&pwd)
	return pwd
}

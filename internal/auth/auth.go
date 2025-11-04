package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func VerifyPassword(passFile string) bool {
	fmt.Print("Enter master password: ")

	var pw []byte
	var err error

	if term.IsTerminal(int(os.Stdin.Fd())) {
		pw, err = term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
	} else {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		pw = []byte(strings.TrimSpace(input))
	}

	if err != nil {
		fmt.Println("error reading password:", err)
		return false
	}

	hashBytes, err := os.ReadFile(passFile)
	if err != nil {
		fmt.Println("error reading password file:", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword(hashBytes, pw)
	return err == nil
}

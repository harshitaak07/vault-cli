package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

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
		input = strings.TrimSpace(input)
		pw = []byte(input)
	}

	if err != nil {
		fmt.Println("error reading password:", err)
		return false
	}

	h := sha256.Sum256(pw)
	expected, err := os.ReadFile(passFile)
	if err != nil {
		fmt.Println("error reading password file:", err)
		return false
	}

	return hex.EncodeToString(h[:]) == strings.TrimSpace(string(expected))
}

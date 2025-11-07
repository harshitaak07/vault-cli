package auth

import (
    "bufio"
    "errors"
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

    ok, err := CheckPassword(passFile, string(pw))
    if err != nil {
        fmt.Println("error verifying password:", err)
        return false
    }
    return ok
}

func CheckPassword(passFile, password string) (bool, error) {
    password = strings.TrimSpace(password)
    if password == "" {
        return false, nil
    }
    if strings.TrimSpace(passFile) == "" {
        return false, errors.New("password file not configured")
    }
    hashBytes, err := os.ReadFile(passFile)
    if err != nil {
        return false, fmt.Errorf("read password file: %w", err)
    }
    if err := bcrypt.CompareHashAndPassword(hashBytes, []byte(password)); err != nil {
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return false, nil
        }
        return false, err
    }
    return true, nil
}

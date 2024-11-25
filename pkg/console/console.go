package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// Info prints an information message.
func Info(a ...interface{}) {
	message := fmt.Sprint(a...)
	fmt.Printf("  %s\n", message)
}

// Success prints a success message.
func Success(a ...interface{}) {
	message := fmt.Sprint(a...)
	fmt.Printf("%s✔%s %s\n", colorGreen, colorReset, message)
}

// Warning prints a warning message.
func Warning(a ...interface{}) {
	message := fmt.Sprint(a...)
	fmt.Printf("%s!%s %s\n", colorYellow, colorReset, message)
}

// Error prints an error message with a newline.
func Error(a ...interface{}) {
	message := fmt.Sprint(a...)
	fmt.Printf("%s✘%s %s\n", colorRed, colorReset, message)
}

// Input prints an input prompt.
func Input(a ...interface{}) {
	message := fmt.Sprint(a...)
	fmt.Printf("%s%s%s", colorYellow, message, colorReset)
}

// ReadLine reads a line from standard input.
func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// ReadPassword reads a password from standard input without echoing.
func ReadPassword() (string, error) {
	Input("Enter server user password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(password), nil
}

// Print prints a message to the console.
func Print(a ...interface{}) {
	fmt.Println(a...)
}

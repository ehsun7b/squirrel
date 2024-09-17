package app

import (
	"bufio"
	"fmt"
	"os"
	"squirrel/types"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

func ReadInput[T any](name string, desc string, mandatory bool, p types.Printer, t *T) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		description := descriptionOfField(desc)
		p("{0}{1}: ", name, description)

		if scanner.Scan() {
			input := scanner.Text()
			if mandatory && strings.TrimSpace(input) == "" {
				p("{red}The {0} field is mandatory and cannot be empty!{/red}\n", name)
				continue
			}

			switch v := any(t).(type) {
			case *string:
				*v = input
			case *int:
				num, err := strconv.Atoi(input)
				if err != nil {
					p("{red}Invalid input for {0}: {1}{/red}\n", name, err)
					continue
				}
				*v = num
			case *int64:
				num, err := strconv.ParseInt(input, 10, 64)
				if err != nil {
					p("{red}Invalid input for {0}: {1}{/red}\n", name, err)
					continue
				}
				*v = num
			default:
				p("{red}Unsupported type!{/red}\n")
				continue
			}
			return
		}

		if err := scanner.Err(); err != nil {
			p("{red}Reading {0} failed!{/red}{1}\n", name, err)
		}
	}
}

func ReadSecret(name string, desc string, mandatory bool, p types.Printer, t *string) {
	for {
		description := descriptionOfField(desc)
		p("{0} {1}: ", name, description)

		password, err := readPasswordWithMask()
		if err != nil {
			p("{red}Reading {0} failed!{/red}{1}\n", name, err)
			continue
		}

		if mandatory && strings.TrimSpace(password) == "" {
			p("{red}The {0} field is mandatory and cannot be empty!{/red}\n", name)
			continue
		}

		*t = password
		return
	}
}

func PrintSecret(p types.Printer, secret string, seconds int) {
	p("{gray}Your chosen master password: {0}{/gray}\n", secret)

	time.Sleep(time.Duration(seconds) * time.Second)

	// Clear the secret from the terminal
	clearLine()
}

func GetYesNoInput(p types.Printer, prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Print the prompt
		p("{0} (y/n): ", prompt)

		// Read the input
		input, err := reader.ReadString('\n')
		if err != nil {
			p("Error reading input. Please try again.\n")
			continue
		}

		// Trim whitespace and convert to lowercase
		input = strings.TrimSpace(strings.ToLower(input))

		// Check for valid input
		if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		} else {
			// If the input is not valid, prompt again
			p("Please enter 'y' or 'n'.\n")
		}
	}
}

// clearLine uses ANSI escape codes to clear the previous line
func clearLine() {
	// Move cursor up one line and clear the line
	fmt.Print("\033[1A\033[K")
}

func readPasswordWithMask() (string, error) {
	var password []byte
	var err error

	password, err = term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(password), nil
}

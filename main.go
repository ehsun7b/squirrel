package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"squirrel/app"
	"squirrel/data"
	l "squirrel/log"
	"squirrel/secure"
	"strings"
)

const appName = "Squirrel"
const appVersion = "0.1.0"

var (
	ErrWrongPassword error = errors.New("wrong password")
)

var (
	state         data.State = data.State{}
	password      []byte     = nil
	encryptionKey []byte     = nil
)

type RunMode int8

const (
	Normal RunMode = iota
	Safe
)

var commands = map[string]app.Command{
	"version": app.VersionCommand(l.Println, appName, appVersion),

	"list": app.ListCommand(l.Print, decryptor),
	"ls":   app.ListCommand(l.Print, decryptor),

	"new":    app.NewCommand(l.Print, encryptor),
	"add":    app.NewCommand(l.Print, encryptor),
	"create": app.NewCommand(l.Print, encryptor),

	"delete": app.DeleteCommand(l.Print),
	"del":    app.DeleteCommand(l.Print),
	"remove": app.DeleteCommand(l.Print),

	"show": app.ShowCommand(l.Print, decryptor),

	"search": app.SearchCommand(l.Print, decryptor),

	"edit": app.EditCommand(l.Print, encryptor, decryptor),

	"help": app.HelpCommand(l.Print),
}

func main() {
	app.Logo()
	l.Println("{brightWhite}{0} v{1}{/brightWhite}", appName, appVersion)
	printLow("Loading...\n")

	key, err := signInOrInitialize()
	if err != nil {
		switch err {
		case ErrWrongPassword:
			l.Println("{bgRed}WRONG PASSWORD!{/bgRed}")
			os.Exit(1)
		default:
			l.Println("{red}Uknown error. {0}{/red}")
			os.Exit(1)
		}
	}

	encryptionKey = key

	runMode(os.Args)
}

func runMode(args []string) {
	if len(args) == 1 {
		interactiveMode(int8(Normal))
	} else {
		os.Exit(1)
	}
}

func interactiveMode(mode int8) {
	state = app.ReadState()
	printLow("There are {0} entries.\n", state.Count)

	reader := bufio.NewReader(os.Stdin)
	for {
		prompt()
		// Read user input
		input, err := getCommand(reader)

		if err != nil {
			l.Println("{red}Error reading input{/red}")
			continue
		}

		// Clean up the input string (remove newline characters)
		input = strings.TrimSpace(input)

		// Handle exit command
		if input == "exit" {
			fmt.Println("Exiting...")
			break
		} else {
			processInput(input)
		}
	}

	// -command loop
}

func prompt() {
	l.Print("{brightGreen}🐿️ ❯{/brightGreen} ")
}

func signInOrInitialize() ([]byte, error) {
	firstRun := !data.HasPassVerifyFile()

	if firstRun {
		l.Println("Initializing master password...")
		l.Println("{magenta}Choose a secure password and make sure to remember it. Without this password, your data will not be recoverable, and there will be no way to reset it.{/magenta}")

		var pass, veryfy string

		for {
			app.ReadSecret("Master password", "", true, l.Print, &pass)
			app.ReadSecret("Verify password", "", true, l.Print, &veryfy)

			if pass != veryfy {
				l.Println("Password did not match!")
			} else {
				show := app.GetYesNoInput(l.Print, "Do you need to see your password for 5 seconds?")

				if show {
					app.PrintSecret(l.Print, pass, 5)
				}
				break
			}
		}

		password = []byte(pass)

		salt, err := secure.DeriveKeyScrypt(password, []byte(""))

		if err != nil {
			l.Println("{red}Can't generate a salt!{/red} {0}", err)
			os.Exit(1)
		}

		key := secure.DeriveKeyPBKDF2(password, salt)
		e, err := secure.EncryptAES("squirrel", key)
		if err != nil {
			l.Println("{red}Can't encrypt sample text with given password!{/red} {0}", e)
			os.Exit(1)
		}

		err = data.SavePassVerify(e)
		if err != nil {
			l.Println("{red}Can't write to disk!{/red} {0}", e)
			os.Exit(1)
		}

		return key, nil
	} else {
		var pass string
		for {
			app.ReadSecret("Enter password", "", false, l.Print, &pass)
			password = []byte(pass)

			salt, err := secure.DeriveKeyScrypt(password, []byte(""))

			if err != nil {
				l.Println("{red}Can't generate a salt!{/red} {0}", err)
				os.Exit(1)
			}

			key := secure.DeriveKeyPBKDF2(password, salt)

			data, err := data.LoadPassVerify()
			if err != nil {
				l.Println("{red}Can't read from disk!{/red} {0}", err)
				os.Exit(1)
			}

			d, err := secure.DecryptAES(data, key)
			if err != nil {
				l.Println("{red}Can't decrypt!{/red} {0}", err)
				os.Exit(1)
			}

			if d == "squirrel" {
				return key, nil
			} else {
				l.Println("{brightWhite}Wrong password!{/white}")
			}
		}
	}
}

func printLow(template string, values ...interface{}) {
	l.Print("{gray}"+template+"{/gray}", values...)
}

func processInput(input string) {
	parts := strings.Split(input, " ")

	command, exists := commands[parts[0]]
	if exists {
		command(parts[1:]...)
	} else {
		command := strings.Split(input, " ")[0]
		l.Println("No {magenta}{0}{/magenta} command.Type {green}help{/green} for available commands.", command)
	}

}

func getCommand(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Clean up the input string (remove newline characters)
	input = strings.TrimSpace(input)
	return input, nil
}

func encryptor(value string) (string, error) {
	en, err := secure.EncryptAES(value, encryptionKey)
	if err != nil {
		l.Println("{red}Encryption failed! {0}{/red}", err)
		return value, err
	}
	return en, nil
}

func decryptor(value string) (string, error) {
	en, err := secure.DecryptAES(value, encryptionKey)
	if err != nil {
		l.Println("{red}Decryption failed! {0}{/red}", err)
		return value, err
	}
	return en, nil
}

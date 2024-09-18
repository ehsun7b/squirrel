package app

import (
	"fmt"
	"squirrel/types"
	"strings"
)

type helpLine struct {
	command     string
	aliases     []string
	description string
	examples    []string
}

func HelpCommand(p types.Printer) Command {
	return func(args ...string) {

		helpLines := []helpLine{
			{
				command:     "version",
				aliases:     []string{},
				description: "Shows version info.",
				examples:    []string{"version"},
			},
			{
				command:     "delete",
				aliases:     []string{"del", "remove"},
				description: "Deletes an entry.",
				examples:    []string{"delete 123", "del", "remove 32"},
			},
			{
				command:     "new",
				aliases:     []string{"create", "add"},
				description: "Creates a new entry.",
				examples:    []string{"new", "create", "add"},
			},
		}
		printHelpLines(helpLines, p)
	}
}

func printHelpLines(helpLines []helpLine, p types.Printer) {
	fmt.Println("Available Commands:")
	fmt.Println(strings.Repeat("=", 20))
	for _, line := range helpLines {
		// Print command and aliases
		p(fmt.Sprintf("Command: {brightWhite}%-10s{/brightWhite} Aliases: %-20s\n", line.command, strings.Join(line.aliases, ", ")))

		// Print description
		fmt.Printf("  Description: %s\n", line.description)

		// Print examples if any
		if len(line.examples) > 0 {
			fmt.Println("  Examples:")
			for _, example := range line.examples {
				fmt.Printf("    - %s\n", example)
			}
		}
		fmt.Println(strings.Repeat("-", 40)) // Divider for readability
	}
}

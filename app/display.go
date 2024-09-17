package app

import (
	"fmt"
	"squirrel/data"
	"squirrel/types"
	"strings"
)

func display(ent data.Entry, p types.Printer) {
	p(format(ent))
}

func format(entry data.Entry) string {
	fields := []struct {
		name  string
		value string
	}{
		{"Id", fmt.Sprintf("%d", entry.Id)},
		{"Title", entry.Title},
		{"Username", entry.Username},
		{"Password", "{gray}" + entry.Password + "{/gray}"},
		{"Address", entry.Address},
		{"Notes", entry.Notes},
	}

	maxFieldLength := 0
	for _, field := range fields {
		if len(field.value) > 0 && len(field.name) > maxFieldLength {
			maxFieldLength = len(field.name)
		}
	}

	// Build the formatted string, skipping empty fields
	var result strings.Builder
	for _, field := range fields {
		if field.value == "" {
			continue // Skip empty fields
		}
		// Calculate how many dots/spaces to add
		var adjust = 10

		if field.name == "password" {
			adjust = adjust - 13 // 13 is length of {gray}{/gray}
		}

		paddingDots := strings.Repeat(" ", maxFieldLength-len(field.name)+adjust) // Adjust +10 to set value column distance
		result.WriteString(fmt.Sprintf("%s%s%s\n", field.name, paddingDots, field.value))
	}

	return result.String()
}

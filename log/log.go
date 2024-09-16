package log

import (
	"fmt"
	"strings"
)

var colorMap = map[string]string{
	"black":           "\033[30m",
	"red":             "\033[31m",
	"green":           "\033[32m",
	"yellow":          "\033[33m",
	"blue":            "\033[34m",
	"magenta":         "\033[35m",
	"cyan":            "\033[36m",
	"white":           "\033[37m",
	"gray":            "\033[90m",
	"brightRed":       "\033[91m",
	"brightGreen":     "\033[92m",
	"brightYellow":    "\033[93m",
	"brightBlue":      "\033[94m",
	"brightMagenta":   "\033[95m",
	"brightCyan":      "\033[96m",
	"brightWhite":     "\033[97m",
	"bgBlack":         "\033[40m",
	"bgRed":           "\033[41m",
	"bgGreen":         "\033[42m",
	"bgYellow":        "\033[43m",
	"bgBlue":          "\033[44m",
	"bgMagenta":       "\033[45m",
	"bgCyan":          "\033[46m",
	"bgWhite":         "\033[47m",
	"bgGray":          "\033[100m",
	"bgBrightRed":     "\033[101m",
	"bgBrightGreen":   "\033[102m",
	"bgBrightYellow":  "\033[103m",
	"bgBrightBlue":    "\033[104m",
	"bgBrightMagenta": "\033[105m",
	"bgBrightCyan":    "\033[106m",
	"bgBrightWhite":   "\033[107m",
	"blink":           "\033[5m",
}

const Reset = "\033[0m"

func Print(template string, values ...interface{}) {
	fmt.Print(appluColor(template, values...))
}

func Println(template string, values ...interface{}) {
	fmt.Println(appluColor(template, values...))
}

// PrintTemplate processes the template with color tags and prints to standard output
func appluColor(template string, values ...interface{}) string {
	// Replace color tags in the template
	for colorName, colorCode := range colorMap {
		startTag := "{" + colorName + "}"
		endTag := "{/" + colorName + "}"
		template = strings.ReplaceAll(template, startTag, colorCode)
		template = strings.ReplaceAll(template, endTag, Reset)
	}

	// Ensure the terminal color is reset at the end
	template += Reset

	// Replace indexed placeholders `{0}`, `{1}`, etc., with actual values
	template = replacePlaceholders(template, values...)

	return template
}

// Replace indexed placeholders `{0}`, `{1}`, etc., with actual values.
func replacePlaceholders(template string, values ...interface{}) string {
	for i, value := range values {
		placeholder := fmt.Sprintf("{%d}", i)
		template = strings.ReplaceAll(template, placeholder, fmt.Sprint(value))
	}
	return template
}

package commands

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

type HelpCommand struct {
}

func (command HelpCommand) Execute() string {
	descriptions := map[string]string{
		"/help":        "show available commands",
		"/stats":       "show your message count",
		"/nick <name>": "change your nickname",
	}

	var items []string

	for name, description := range descriptions {
		items = append(items, name+": "+description)
	}

	sort.Strings(items)

	width := 0
	for _, item := range items {
		length := utf8.RuneCountInString(item)
		if length > width {
			width = length
		}
	}

	border := "\t+" + strings.Repeat("-", width+2) + "+\n"

	var result strings.Builder
	result.WriteString(border)
	for _, item := range items {
		result.WriteString(fmt.Sprintf("\t| %-*s |\n", width, item))
	}

	return result.String() + border
}

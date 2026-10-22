package commands

import (
	"sort"
	"strings"
)

type HelpCommand struct {
}

func (command HelpCommand) Execute() string {
	descriptions := map[string]string{
		"/help":  "show available commands",
		"/stats": "show your message count",
	}

	var items []string

	for name, description := range descriptions {
		items = append(items, name+": "+description)
	}

	sort.Strings(items)

	return strings.Join(items, " | ") + "\n"
}

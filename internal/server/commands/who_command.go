package commands

import (
	"fmt"
	"sort"
	"strings"
)

type WhoCommand struct {
	Nicknames []string
}

func (command WhoCommand) Execute() string {
	sort.Strings(command.Nicknames)

	var result strings.Builder
	result.WriteString(
		fmt.Sprintf("Clients online (%d):\n", len(command.Nicknames)),
	)

	for _, nickname := range command.Nicknames {
		result.WriteString("\t" + nickname + "\n")
	}

	return result.String()
}

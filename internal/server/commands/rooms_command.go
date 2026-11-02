package commands

import (
	"fmt"
	"sort"
	"strings"
)

type RoomsCommand struct {
	Rooms map[string]int
}

func (command RoomsCommand) Execute() string {
	names := make([]string, 0, len(command.Rooms))

	for name := range command.Rooms {
		names = append(names, name)
	}

	sort.Strings(names)

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Rooms (%d):\n", len(names)))

	for _, name := range names {
		result.WriteString(fmt.Sprintf("\t%s (%d)\n", name, command.Rooms[name]))
	}

	return result.String()
}

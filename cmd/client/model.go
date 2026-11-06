package main

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

type chatModel struct {
	input    textinput.Model
	viewport viewport.Model
	messages []string
}

package core

import (
	"github.com/jpinilloslr/actionai/internal/core/input"
	"github.com/jpinilloslr/actionai/internal/core/output"
)

type command struct {
	Model        string
	InputType    input.Type
	OutputType   output.Type
	Shortcut     string
	Instructions string
	Notify       bool
}

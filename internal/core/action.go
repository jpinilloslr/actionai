package core

import (
	"github.com/jpinilloslr/actionai/internal/core/input"
	"github.com/jpinilloslr/actionai/internal/core/output"
)

type action struct {
	Model        string
	Inputs       []input.Type
	Output       output.Type
	Shortcut     string
	Instructions string
	Notify       bool
}

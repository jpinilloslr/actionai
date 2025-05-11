package core

import "github.com/jpinilloslr/actionai/internal/core/input"

type AIModel interface {
	Run(
		model string,
		instructions string,
		inputs []input.Input,
	) (string, error)
}

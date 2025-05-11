package core

import "github.com/jpinilloslr/actionai/internal/core/input"

type AIModel interface {
	TextToSpeech(text string) error
	SpeechToText(audioFile string) (string, error)
	Run(
		model string,
		instructions string,
		inputs []input.Input,
	) (string, error)
}

package core

type AiModel interface {
	TextToSpeech(text string) error
	SpeechToText(audioFile string) (string, error)
	RunWithText(model string, instructions string, text string) (string, error)
	RunWithImage(model string, instructions string, data string) (string, error)
}

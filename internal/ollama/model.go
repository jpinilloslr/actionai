package ollama

import (
	"fmt"
	"log/slog"

	"github.com/jpinilloslr/actionai/internal/core"
)

type AIModel struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) (core.AiModel, error) {
	m := AIModel{
		logger: logger,
	}

	if err := m.init(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *AIModel) init() error {
	return fmt.Errorf("not implemented yet")
}

func (m *AIModel) RunWithText(
	model string,
	instructions string,
	text string,
) (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (m *AIModel) RunWithImage(
	model string,
	instructions string,
	data string,
) (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (m *AIModel) SpeechToText(audioFile string) (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (m *AIModel) TextToSpeech(text string) error {
	return fmt.Errorf("not implemented yet")
}

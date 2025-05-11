package ollama

import (
	"fmt"
	"log/slog"

	"github.com/jpinilloslr/actionai/internal/core"
)

type AiModel struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) (core.AiModel, error) {
	m := AiModel{
		logger: logger,
	}

	if err := m.init(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *AiModel) init() error {
	return fmt.Errorf("not implemented yet")
}

func (m *AiModel) Run(
	model string,
	instructions string,
	inputs []core.Input,
) (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (m *AiModel) SpeechToText(audioFile string) (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (m *AiModel) TextToSpeech(text string) error {
	return fmt.Errorf("not implemented yet")
}

package ollama

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jpinilloslr/actionai/internal/core"
)

type VoiceEngine struct {
	logger *slog.Logger
}

func NewVoiceEngine(logger *slog.Logger) (core.VoiceEngine, error) {
	m := VoiceEngine{
		logger: logger,
	}

	if err := m.init(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *VoiceEngine) init() error {
	return fmt.Errorf("not implemented yet")
}

func (m *VoiceEngine) Transcribe(audioFile string) (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (m *VoiceEngine) Speak(ctx context.Context, text string) error {
	return fmt.Errorf("not implemented yet")
}

package openai

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jpinilloslr/actionai/internal/core"
	"github.com/openai/openai-go"
)

type VoiceEngine struct {
	apiKey string
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
	m.logger.Info("Initializing OpenAI model")
	key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return fmt.Errorf("OPENAI_API_KEY not set")
	}

	m.apiKey = key
	return nil
}

func (m *VoiceEngine) Transcribe(audioFile string) (string, error) {
	client := openai.NewClient()
	ctx := context.Background()

	file, err := os.Open(audioFile)
	if err != nil {
		return "", err
	}

	transcription, err := client.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		Model: openai.AudioModelWhisper1,
		File:  file,
	})
	if err != nil {
		return "", err
	}

	return transcription.Text, nil
}

func (m *VoiceEngine) Speak(text string) error {
	return fmt.Errorf("not implemented")
}

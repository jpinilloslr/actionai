package openai

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/jpinilloslr/actionai/internal/core"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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
	client := openai.NewClient(
		option.WithAPIKey(m.apiKey),
	)
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
	client := openai.NewClient(
		option.WithAPIKey(m.apiKey),
	)
	ctx := context.Background()

	res, err := client.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Input:          text,
		Model:          openai.SpeechModelTTS1,
		Voice:          openai.AudioSpeechNewParamsVoiceNova,
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatPCM,
	})
	defer res.Body.Close()

	if err != nil {
		return err
	}

	op := &oto.NewContextOptions{}
	op.SampleRate = 24000
	op.ChannelCount = 1
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return err
	}

	<-readyChan

	player := otoCtx.NewPlayer(res.Body)
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	return player.Close()
}

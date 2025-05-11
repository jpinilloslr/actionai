package openai

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jpinilloslr/actionai/internal/core"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AiModel struct {
	apiKey string
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
	m.logger.Info("Initializing OpenAI model")
	key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return fmt.Errorf("OPENAI_API_KEY not set")
	}

	m.apiKey = key
	return nil
}

func (m *AiModel) Run(
	model string,
	instructions string,
	inputs []core.Input,
) (string, error) {
	m.logger.Info(
		"Running OpenAI model",
		slog.Any("model", model),
		slog.Any("instructions", instructions),
		slog.Any("inputs", inputs),
	)

	client := openai.NewClient(
		option.WithAPIKey(m.apiKey),
	)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(instructions),
	}

	for _, input := range inputs {
		if input.Text != nil {
			messages = append(messages, openai.UserMessage(*input.Text))
		}
		if input.ImageData != nil {
			messages = append(
				messages,
				openai.UserMessage(
					[]openai.ChatCompletionContentPartUnionParam{
						openai.ImageContentPart(
							openai.ChatCompletionContentPartImageImageURLParam{
								URL: *input.ImageData,
							}),
					},
				),
			)
		}
	}

	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Model:    model,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("empty response from model")
	}

	content := chatCompletion.Choices[0].Message.Content
	m.logger.Info("Response from OpenAI model",
		slog.Any("response", content))

	return content, nil
}

func (m *AiModel) SpeechToText(audioFile string) (string, error) {
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

func (m *AiModel) TextToSpeech(text string) error {
	return fmt.Errorf("not implemented")
}

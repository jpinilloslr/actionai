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

type AIModel struct {
	apiKey string
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
	m.logger.Info("Initializing OpenAI model")
	key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return fmt.Errorf("OPENAI_API_KEY not set")
	}

	m.apiKey = key
	return nil
}

func (m *AIModel) RunWithText(
	model string,
	instructions string,
	text string,
) (string, error) {
	m.logger.Info(
		"Running OpenAI model with text",
		slog.Any("model", model),
		slog.Any("instructions", instructions),
		slog.Any("text", text),
	)

	client := openai.NewClient(
		option.WithAPIKey(m.apiKey),
	)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(instructions),
		openai.UserMessage(text),
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
		return "", fmt.Errorf("Empty response from model")
	}

	content := chatCompletion.Choices[0].Message.Content
	m.logger.Info("Response from OpenAI model",
		slog.Any("response", content))

	return content, nil
}

func (m *AIModel) RunWithImage(
	model string,
	instructions string,
	data string,
) (string, error) {
	m.logger.Info(
		"Running OpenAI model with image",
		slog.Any("model", model),
		slog.Any("instructions", instructions),
	)

	client := openai.NewClient(
		option.WithAPIKey(m.apiKey),
	)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(instructions),
		openai.UserMessage(
			[]openai.ChatCompletionContentPartUnionParam{
				openai.ImageContentPart(
					openai.ChatCompletionContentPartImageImageURLParam{
						URL: data,
					}),
			},
		),
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
		return "", fmt.Errorf("Empty response from model")
	}

	content := chatCompletion.Choices[0].Message.Content
	m.logger.Info("Response from OpenAI model",
		slog.Any("response", content))

	return content, nil
}

func (m *AIModel) SpeechToText(audioFile string) (string, error) {
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

func (m *AIModel) TextToSpeech(text string) error {
	return fmt.Errorf("TextToSpeech not implemented")
}

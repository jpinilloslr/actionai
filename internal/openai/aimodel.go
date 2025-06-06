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

func NewAIModel(logger *slog.Logger) (core.AIModel, error) {
	m := AIModel{
		logger: logger,
	}

	if err := m.init(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *AIModel) init() error {
	m.logger.Info("Initializing OpenAI model for content generation")
	key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return fmt.Errorf("OPENAI_API_KEY not set")
	}

	m.apiKey = key
	return nil
}

func (m *AIModel) Run(
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

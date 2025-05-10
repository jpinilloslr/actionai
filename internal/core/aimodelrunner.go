package core

import (
	"fmt"
	"log/slog"

	"github.com/jpinilloslr/actionai/internal/core/input"
	"github.com/jpinilloslr/actionai/internal/core/output"
	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type AiModelRunner struct {
	aiModel    AiModel
	logger     *slog.Logger
	installer  *installer
	actionRepo *actionRepo
	inReceiver *input.Receiver
	outSender  *output.Sender
	notifier   platform.Notifier
}

func NewAiModelRunner(
	logger *slog.Logger,
	workDir *WorkDir,
	aiModel AiModel,
	dialog platform.Dialog,
	notifier platform.Notifier,
	clipboard platform.Clipboard,
	screenshotter platform.Screenshotter,
	speechRecorder platform.SpeechRecorder,
	selTextProvider platform.SelTextProvider,
	shortcutsCreator platform.ShortcutCreator,
) (*AiModelRunner, error) {
	cmdRepo, err := newActionRepo(logger, workDir.ActionsFile())
	if err != nil {
		return nil, err
	}

	inReceiver := input.New(
		dialog,
		clipboard,
		screenshotter,
		speechRecorder,
		selTextProvider,
	)
	outSender := output.New(dialog, clipboard)
	installer := newInstaller(logger, cmdRepo, shortcutsCreator)

	return &AiModelRunner{
		logger:     logger,
		aiModel:    aiModel,
		actionRepo: cmdRepo,
		notifier:   notifier,
		outSender:  outSender,
		installer:  installer,
		inReceiver: inReceiver,
	}, nil
}

func (r *AiModelRunner) InstallShortcuts() error {
	if err := r.installer.Install(); err != nil {
		return err
	}

	return nil
}

func (r *AiModelRunner) Run(actionId string) error {
	action, err := r.actionRepo.GetById(actionId)
	if err != nil {
		return err
	}

	input, err := r.inReceiver.Receive(action.Inputs[0])
	if err != nil {
		return err
	}

	resp, err := r.run(action, input)
	if err != nil {
		return err
	}

	if action.Notify {
		err := r.notifier.Notify("Action AI", actionId+" completed.")
		if err != nil {
			return err
		}
	}

	return r.outSender.Send(action.Output, resp)
}

func (r *AiModelRunner) run(action *action, input *input.Input) (string, error) {
	if input.Text != nil {
		return r.aiModel.RunWithText(action.Model, action.Instructions, *input.Text)
	}

	if input.ImageData != nil {
		return r.aiModel.RunWithImage(
			action.Model,
			action.Instructions,
			*input.ImageData,
		)
	}

	if input.SpeechFileName != nil {
		text, err := r.aiModel.SpeechToText(*input.SpeechFileName)
		if err != nil {
			return "", err
		}
		return r.aiModel.RunWithText(action.Model, action.Instructions, text)
	}

	return "", fmt.Errorf("unsupported input")
}

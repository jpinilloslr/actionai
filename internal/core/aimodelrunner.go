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
	cmdRepo    *commandRepo
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
	cmdRepo, err := newCommandRepo(logger, workDir.CommandsFile())
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
		cmdRepo:    cmdRepo,
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

func (r *AiModelRunner) Run(command string) error {
	cmd, err := r.cmdRepo.GetById(command)
	if err != nil {
		return err
	}

	input, err := r.inReceiver.Receive(cmd.InputType)
	if err != nil {
		return err
	}

	resp, err := r.run(cmd, input)
	if err != nil {
		return err
	}

	if cmd.Notify {
		err := r.notifier.Notify("AI Shortcuts", command+" completed.")
		if err != nil {
			return err
		}
	}

	return r.outSender.Send(cmd.OutputType, resp)
}

func (r *AiModelRunner) run(cmd *command, input *input.Input) (string, error) {
	if input.Text != nil {
		return r.aiModel.RunWithText(cmd.Model, cmd.Instructions, *input.Text)
	}

	if input.ImageData != nil {
		return r.aiModel.RunWithImage(cmd.Model, cmd.Instructions, *input.ImageData)
	}

	if input.SpeechFileName != nil {
		text, err := r.aiModel.SpeechToText(*input.SpeechFileName)
		if err != nil {
			return "", err
		}
		return r.aiModel.RunWithText(cmd.Model, cmd.Instructions, text)
	}

	return "", fmt.Errorf("unsupported input")
}

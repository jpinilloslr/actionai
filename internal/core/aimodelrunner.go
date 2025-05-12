package core

import (
	"log/slog"

	"github.com/jpinilloslr/actionai/internal/core/input"
	"github.com/jpinilloslr/actionai/internal/core/output"
	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type AIModelRunner struct {
	aiModel     AIModel
	voiceEngine VoiceEngine
	logger      *slog.Logger
	installer   *installer
	actionRepo  *actionRepo
	inReceiver  *input.Receiver
	outSender   *output.Sender
	notifier    platform.Notifier
}

func NewAIModelRunner(
	logger *slog.Logger,
	workDir *WorkDir,
	aiModel AIModel,
	voiceEngine VoiceEngine,
	dialog platform.Dialog,
	notifier platform.Notifier,
	clipboard platform.Clipboard,
	shortcutsMgr platform.ShortcutsMgr,
	screenshotter platform.Screenshotter,
	voiceRecorder platform.VoiceRecorder,
	selTextProvider platform.SelTextProvider,
) (*AIModelRunner, error) {
	cmdRepo, err := newActionRepo(logger, workDir.ActionsFile())
	if err != nil {
		return nil, err
	}

	inReceiver := input.New(
		dialog,
		clipboard,
		screenshotter,
		voiceRecorder,
		selTextProvider,
	)
	outSender := output.New(dialog, clipboard)
	installer := newInstaller(logger, cmdRepo, shortcutsMgr)

	return &AIModelRunner{
		logger:      logger,
		aiModel:     aiModel,
		actionRepo:  cmdRepo,
		notifier:    notifier,
		outSender:   outSender,
		installer:   installer,
		inReceiver:  inReceiver,
		voiceEngine: voiceEngine,
	}, nil
}

func (r *AIModelRunner) InstallShortcuts() error {
	if err := r.installer.Install(); err != nil {
		return err
	}

	return nil
}

func (r *AIModelRunner) Run(actionId string) error {
	action, err := r.actionRepo.GetById(actionId)
	if err != nil {
		return err
	}

	inputs, err := r.inReceiver.Receive(action.Inputs)
	if err != nil {
		return err
	}

	resp, err := r.run(action, inputs)
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

func (r *AIModelRunner) run(action *action, inputs []input.Input) (string, error) {
	r.processVoiceInput(inputs)
	return r.aiModel.Run(action.Model, action.Instructions, inputs)
}

func (r *AIModelRunner) processVoiceInput(inputs []input.Input) error {
	for i := range inputs {
		if inputs[i].VoiceFileName != nil {
			text, err := r.voiceEngine.Transcribe(*inputs[i].VoiceFileName)
			if err != nil {
				return err
			}

			inputs[i].Text = &text
			inputs[i].VoiceFileName = nil
		}
	}

	return nil
}

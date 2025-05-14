package core

import (
	"context"
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
	assetsMgr   *AssetsMgr
	audioPlayer platform.AudioPlayer
}

func NewAIModelRunner(
	logger *slog.Logger,
	assetsMgr *AssetsMgr,
	aiModel AIModel,
	voiceEngine VoiceEngine,
	dialog platform.Dialog,
	notifier platform.Notifier,
	clipboard platform.Clipboard,
	shortcutsMgr platform.ShortcutsMgr,
	screenshotter platform.Screenshotter,
	voiceRecorder platform.VoiceRecorder,
	selTextProvider platform.SelTextProvider,
	audioPlayer platform.AudioPlayer,
) (*AIModelRunner, error) {
	actionRepo, err := newActionRepo(logger, assetsMgr.ActionsFile())
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
	outSender := output.New(dialog, clipboard, voiceEngine.Speak)
	installer := newInstaller(logger, actionRepo, shortcutsMgr)

	return &AIModelRunner{
		logger:      logger,
		aiModel:     aiModel,
		notifier:    notifier,
		outSender:   outSender,
		installer:   installer,
		assetsMgr:   assetsMgr,
		inReceiver:  inReceiver,
		actionRepo:  actionRepo,
		voiceEngine: voiceEngine,
		audioPlayer: audioPlayer,
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.audioPlayer.PlayLoop(ctx, r.assetsMgr.SoundFile())
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

package core

import (
	"context"
	"log/slog"
	"os"
	"os/exec"

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

func (r *AIModelRunner) playSound(ctx context.Context) error {
	soundFile := r.assetsMgr.SoundFile()
	if _, err := os.Stat(soundFile); err != nil {
		return nil
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				r.logger.Info("Stopping sound playback")
				return
			default:
				r.logger.Info("Playing sound", "file", soundFile)
				cmd := exec.Command("aplay", soundFile)
				if err := cmd.Run(); err != nil {
					r.logger.Error("Error playing sound", "error", err)
				}
			}
		}
	}()

	return nil
}

func (r *AIModelRunner) run(action *action, inputs []input.Input) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.playSound(ctx)
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

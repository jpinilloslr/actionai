package core

import (
	"context"
	"log/slog"

	"github.com/jpinilloslr/actionai/internal/core/input"
	"github.com/jpinilloslr/actionai/internal/core/output"
	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type ActionRunner struct {
	actionRepo  *actionRepo
	aiModel     AIModel
	assetsMgr   *AssetsMgr
	audioPlayer platform.AudioPlayer
	inReceiver  *input.Receiver
	logger      *slog.Logger
	notifier    platform.Notifier
	outSender   *output.Sender
	voiceEngine VoiceEngine
}

func NewActionRunner(
	logger *slog.Logger,
	assetsMgr *AssetsMgr,
	aiModel AIModel,
	voiceEngine VoiceEngine,
	dialog platform.Dialog,
	notifier platform.Notifier,
	clipboard platform.Clipboard,
	audioPlayer platform.AudioPlayer,
	screenshotter platform.Screenshotter,
	voiceRecorder platform.VoiceRecorder,
	selTextProvider platform.SelTextProvider,
) (*ActionRunner, error) {
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

	return &ActionRunner{
		actionRepo:  actionRepo,
		aiModel:     aiModel,
		assetsMgr:   assetsMgr,
		audioPlayer: audioPlayer,
		inReceiver:  inReceiver,
		logger:      logger,
		notifier:    notifier,
		outSender:   outSender,
		voiceEngine: voiceEngine,
	}, nil
}

func (r *ActionRunner) RunFromActionRepo(ctx context.Context, actionId string) error {
	action, err := r.actionRepo.GetById(actionId)
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return r.RunFromAction(ctx, action)
}

func (r *ActionRunner) RunFromAction(ctx context.Context, action *action) error {
	inputs, err := r.inReceiver.Receive(ctx, action.Inputs)
	if err != nil {
		return err
	}

	resp, err := r.run(ctx, action, inputs)
	if err != nil {
		return err
	}

	if action.Notify {
		err := r.notifier.Notify("Action AI", "Action completed.")
		if err != nil {
			return err
		}
	}

	return r.outSender.Send(action.Output, resp)
}

func (r *ActionRunner) run(ctx context.Context, action *action, inputs []input.Input) (string, error) {
	soundCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	r.audioPlayer.PlayLoop(soundCtx, r.assetsMgr.SoundFile())
	r.processVoiceInput(ctx, inputs)
	return r.aiModel.Run(action.Model, action.Instructions, inputs)
}

func (r *ActionRunner) processVoiceInput(ctx context.Context, inputs []input.Input) error {
	for i := range inputs {
		if inputs[i].VoiceFileName != nil {
			text, err := r.voiceEngine.Transcribe(ctx, *inputs[i].VoiceFileName)
			if err != nil {
				return err
			}

			inputs[i].Text = &text
			inputs[i].VoiceFileName = nil
		}
	}

	return nil
}

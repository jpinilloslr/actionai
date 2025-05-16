package input

import (
	"context"
	"fmt"

	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type Receiver struct {
	dialog          platform.Dialog
	clipboard       platform.Clipboard
	screenshotter   platform.Screenshotter
	voiceRecorder   platform.VoiceRecorder
	selTextProvider platform.SelTextProvider
}

func New(
	dialog platform.Dialog,
	clipboard platform.Clipboard,
	screenshotter platform.Screenshotter,
	voiceRecorder platform.VoiceRecorder,
	selTextProvider platform.SelTextProvider,
) *Receiver {
	return &Receiver{
		dialog:          dialog,
		clipboard:       clipboard,
		screenshotter:   screenshotter,
		voiceRecorder:   voiceRecorder,
		selTextProvider: selTextProvider,
	}
}

func (r *Receiver) Receive(ctx context.Context, types []Type) ([]Input, error) {
	result := []Input{}

	for _, inType := range types {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			in, err := r.process(inType)
			if err != nil {
				return nil, err
			}
			result = append(result, *in)
		}
	}

	return result, nil
}

func (r *Receiver) process(inType Type) (*Input, error) {
	switch inType {
	case Clipboard:
		return r.getClipboard()
	case Window:
		return r.getFromWindow()
	case SelectedText:
		return r.getSelectedText()
	case Screen:
		return r.getScreen()
	case ScreenSection:
		return r.getScreenSection()
	case Voice:
		return r.getVoice()
	default:
		return nil, fmt.Errorf("unsupported input source: %s", inType)
	}
}

func (r *Receiver) getClipboard() (*Input, error) {
	isText, err := r.clipboard.IsText()
	if err != nil {
		return nil, err
	}

	if isText {
		text, err := r.clipboard.GetText()
		if err != nil {
			return nil, err
		}

		return &Input{
			Text: &text,
		}, nil
	}

	b64Data, err := r.clipboard.GetBase64()
	if err != nil {
		return nil, err
	}

	return &Input{
		ImageData: &b64Data,
	}, nil
}

func (r *Receiver) getFromWindow() (*Input, error) {
	text, err := r.dialog.Prompt()
	if err != nil {
		return nil, err
	}

	return &Input{
		Text: &text,
	}, nil
}

func (r *Receiver) getSelectedText() (*Input, error) {
	text, err := r.selTextProvider.Get()
	if err != nil {
		return nil, err
	}

	return &Input{
		Text: &text,
	}, nil
}

func (r *Receiver) getScreen() (*Input, error) {
	data, err := r.screenshotter.GetScreenB64()
	if err != nil {
		return nil, err
	}

	return &Input{
		ImageData: &data,
	}, nil
}

func (r *Receiver) getScreenSection() (*Input, error) {
	data, err := r.screenshotter.GetSectionB64()
	if err != nil {
		return nil, err
	}

	return &Input{
		ImageData: &data,
	}, nil
}

func (r *Receiver) getVoice() (*Input, error) {
	fileName, err := r.voiceRecorder.Record()
	if err != nil {
		return nil, err
	}

	return &Input{
		VoiceFileName: &fileName,
	}, nil
}

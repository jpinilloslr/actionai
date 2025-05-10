package input

import (
	"fmt"

	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type Receiver struct {
	dialog          platform.Dialog
	clipboard       platform.Clipboard
	screenshotter   platform.Screenshotter
	speechRecorder  platform.SpeechRecorder
	selTextProvider platform.SelTextProvider
}

func New(
	dialog platform.Dialog,
	clipboard platform.Clipboard,
	screenshotter platform.Screenshotter,
	speechRecorder platform.SpeechRecorder,
	selTextProvider platform.SelTextProvider,
) *Receiver {
	return &Receiver{
		dialog:          dialog,
		clipboard:       clipboard,
		screenshotter:   screenshotter,
		speechRecorder:  speechRecorder,
		selTextProvider: selTextProvider,
	}
}

func (r *Receiver) Receive(inType Type) (*Input, error) {
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
	recordFile, err := r.speechRecorder.Record()
	if err != nil {
		return nil, err
	}

	return &Input{
		SpeechFileName: &recordFile,
	}, nil
}

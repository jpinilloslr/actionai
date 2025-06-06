package output

import (
	"context"
	"fmt"

	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type Sender struct {
	dialog    platform.Dialog
	clipboard platform.Clipboard
	speak     func(context.Context, string) error
}

func New(
	dialog platform.Dialog,
	clipboard platform.Clipboard,
	speak func(context.Context, string) error,
) *Sender {
	return &Sender{
		dialog:    dialog,
		clipboard: clipboard,
		speak:     speak,
	}
}

func (s *Sender) Send(outType Type, value string) error {
	switch outType {
	case Stdout:
		fmt.Println(value)
		return nil
	case Clipboard:
		if err := s.setClipboard(value); err != nil {
			return err
		}
		return nil
	case Window:
		if err := s.setWindow(value); err != nil {
			return err
		}
		return nil
	case Voice:
		if err := s.setVoice(value); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported output type: %s", outType)
	}
}

func (s *Sender) setClipboard(value string) error {
	if err := s.clipboard.SetText(value); err != nil {
		return err
	}
	return nil
}

func (s *Sender) setWindow(value string) error {
	return s.dialog.ShowMultiline(value)
}

func (s *Sender) setVoice(value string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- s.dialog.ShowCancellableDialog(ctx, "Speaking...")
	}()

	go func() {
		errCh <- s.speak(ctx, value)
	}()

	err := <-errCh
	cancel()
	<-errCh
	return err
}

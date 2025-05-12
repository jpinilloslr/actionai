package platform

import "context"

type Dialog interface {
	Prompt() (string, error)
	ShowMultiline(text string) error
	ShowCancellableDialog(ctx context.Context, text string) error
}

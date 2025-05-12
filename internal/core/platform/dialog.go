package platform

import "context"

type Dialog interface {
	Prompt() (string, error)
	Show(text string) error
	ShowInfo(ctx context.Context, text string) error
}

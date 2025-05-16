package core

import "context"

type VoiceEngine interface {
	Speak(ctx context.Context, text string) error
	Transcribe(ctx context.Context, audioFile string) (string, error)
}

package core

import "context"

type VoiceEngine interface {
	Speak(ctx context.Context, text string) error
	Transcribe(audioFile string) (string, error)
}

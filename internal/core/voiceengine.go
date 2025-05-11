package core

type VoiceEngine interface {
	Speak(text string) error
	Transcribe(audioFile string) (string, error)
}

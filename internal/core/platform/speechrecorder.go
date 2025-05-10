package platform

type SpeechRecorder interface {
	Record() (string, error)
}

package platform

type VoiceRecorder interface {
	Record() (string, error)
}

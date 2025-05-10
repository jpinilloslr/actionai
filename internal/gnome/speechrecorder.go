package gnome

import (
	"os"
	"os/exec"
)

type SpeechRecorder struct {
}

func NewSpeechRecorder() *SpeechRecorder {
	return &SpeechRecorder{}
}

func (p *SpeechRecorder) Record() (string, error) {
	fileName := "/tmp/audio.mp3"

	recCmd := exec.Command(
		"ffmpeg",
		"-y",
		"-f", "alsa",
		"-i", "default",
		"-acodec", "libmp3lame",
		fileName,
	)
	if err := recCmd.Start(); err != nil {
		return "", err
	}

	uiCmd := exec.Command(
		"zenity",
		"--info",
		"--text=Recording...",
	)

	if err := uiCmd.Run(); err != nil {
		_ = recCmd.Process.Kill()
		return "", err
	}

	if err := recCmd.Process.Signal(os.Interrupt); err != nil {
		_ = recCmd.Process.Kill()
	}

	_ = recCmd.Wait()

	return fileName, nil
}

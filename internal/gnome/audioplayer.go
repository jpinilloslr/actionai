package gnome

import (
	"context"
	"os"
	"os/exec"
)

type AudioPlayer struct {
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{}
}

func (a *AudioPlayer) PlayLoop(ctx context.Context, fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		return
	}
	cmd := exec.Command("aplay", fileName)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cmd.Start(); err != nil {
					return
				}
			}
		}
	}()
}

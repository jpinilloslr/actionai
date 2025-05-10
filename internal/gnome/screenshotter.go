package gnome

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
)

type Screenshotter struct {
}

func NewScreenshotter() *Screenshotter {
	return &Screenshotter{}
}

func (p *Screenshotter) GetSectionB64() (string, error) {
	return p.getScreenshot(false)
}

func (p *Screenshotter) GetScreenB64() (string, error) {
	return p.getScreenshot(true)
}

func (p *Screenshotter) getScreenshot(fullScreen bool) (string, error) {
	tmpImg := "/tmp/screenshot.png"

	args := []string{"gnome-screenshot"}

	if !fullScreen {
		args = append(args, "-a")
	}

	args = append(args, "-f", tmpImg)

	err := exec.Command(args[0], args[1:]...).Run()
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpImg)

	data, err := os.ReadFile(tmpImg)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	imgData := fmt.Sprintf("data:image/png;base64,%s", b64)

	return imgData, nil
}

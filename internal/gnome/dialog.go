package gnome

import (
	"os/exec"
	"strings"
)

type Dialog struct {
}

func NewDialog() *Dialog {
	return &Dialog{}
}

func (d *Dialog) Prompt() (string, error) {
	cmd := exec.Command("zenity",
		"--text-info",
		"--editable",
		"--title=Action AI",
		"--width=500",
		"--height=500",
	)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	text := strings.TrimSpace(string(output))
	return text, nil
}

func (d *Dialog) Show(text string) error {
	cmd := exec.Command("zenity",
		"--text-info",
		"--title=Action AI",
		"--width=500",
		"--height=500",
	)
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

package gnome

import "os/exec"

type SelTextProvider struct {
}

func NewSelTextProvider() *SelTextProvider {
	return &SelTextProvider{}
}

func (p *SelTextProvider) Get() (string, error) {
	out, err := exec.Command("wl-paste", "--primary").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

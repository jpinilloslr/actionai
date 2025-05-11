package gnome

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
)

type Clipboard struct {
}

func NewClipboard() *Clipboard {
	return &Clipboard{}
}

func (c *Clipboard) SetText(text string) error {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = strings.NewReader(text)

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Clipboard) GetText() (string, error) {
	out, err := c.getContent()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (c *Clipboard) IsText() (bool, error) {
	mimeTypes, err := c.getMimeTypes()
	if err != nil {
		return false, err
	}

	for _, mimeType := range mimeTypes {
		if strings.Contains(mimeType, "text/") {
			return true, nil
		}
	}

	return false, nil
}

func (c *Clipboard) GetBase64() (string, error) {
	mimeTypes, err := c.getMimeTypes()
	if err != nil {
		return "", err
	}

	if len(mimeTypes) == 0 {
		return "", fmt.Errorf("unable to detect clipboard content mime type")
	}

	content, err := c.getContent()
	if err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString(content)
	encData := fmt.Sprintf("data:%s;base64,%s", mimeTypes[0], b64)
	return encData, nil
}

func (c *Clipboard) getMimeTypes() ([]string, error) {
	mimeTypes, err := exec.Command("wl-paste", "--list-types").Output()
	if err != nil {
		return nil, err
	}

	mimeTypesList := strings.Split(string(mimeTypes), "\n")
	return mimeTypesList, nil
}

func (c *Clipboard) getContent() ([]byte, error) {
	return exec.Command("wl-paste").Output()
}

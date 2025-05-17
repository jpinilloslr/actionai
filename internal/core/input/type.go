package input

import "fmt"

type Type string

const (
	Clipboard     Type = "clipboard"
	Screen        Type = "screen"
	ScreenSection Type = "screen-section"
	SelectedText  Type = "selected-text"
	Voice         Type = "voice"
	Window        Type = "window"
)

var validTypes = map[string]Type{
	"clipboard":      Clipboard,
	"screen":         Screen,
	"screen-section": ScreenSection,
	"selected-text":  SelectedText,
	"voice":          Voice,
	"window":         Window,
}

func ParseType(s string) (Type, error) {
	if t, ok := validTypes[s]; ok {
		return t, nil
	}
	return "", fmt.Errorf("invalid input type: %s", s)
}

func ParseTypeList(input []string) ([]Type, error) {
	result := make([]Type, 0, len(input))
	for _, s := range input {
		t, err := ParseType(s)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

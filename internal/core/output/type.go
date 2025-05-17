package output

import "fmt"

type Type string

const (
	Clipboard Type = "clipboard"
	Stdout    Type = "stdout"
	Voice     Type = "voice"
	Window    Type = "window"
)

var validTypes = map[string]Type{
	"clipboard": Clipboard,
	"stdout":    Stdout,
	"voice":     Voice,
	"window":    Window,
}

func ParseType(s string) (Type, error) {
	if t, ok := validTypes[s]; ok {
		return t, nil
	}
	return "", fmt.Errorf("invalid output type: %s", s)
}

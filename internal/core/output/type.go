package output

type Type string

const (
	Stdout    Type = "stdout"
	Clipboard Type = "clipboard"
	Window    Type = "window"
	Voice     Type = "voice"
)

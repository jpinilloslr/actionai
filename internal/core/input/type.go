package input

type Type string

const (
	Clipboard     Type = "clipboard"
	Screen        Type = "screen"
	ScreenSection Type = "screen-section"
	SelectedText  Type = "selected-text"
	Voice         Type = "voice"
	Window        Type = "window"
)

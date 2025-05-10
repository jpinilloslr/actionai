package platform

type Dialog interface {
	Prompt() (string, error)
	Show(text string) error
}

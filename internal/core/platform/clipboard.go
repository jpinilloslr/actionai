package platform

type Clipboard interface {
	SetText(text string) error
	GetText() (string, error)
	IsText() (bool, error)
	GetBase64() (string, error)
}

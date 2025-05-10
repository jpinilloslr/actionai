package platform

type Screenshotter interface {
	GetScreenB64() (string, error)
	GetSectionB64() (string, error)
}

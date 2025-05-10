package platform

type SelTextProvider interface {
	Get() (string, error)
}

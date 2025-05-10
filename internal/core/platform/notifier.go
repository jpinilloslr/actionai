package platform

type Notifier interface {
	Notify(title string, text string) error
}

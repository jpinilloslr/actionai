package gnome

import "os/exec"

type Notifier struct{}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) Notify(title string, text string) error {
	cmd := exec.Command("notify-send", title, text)
	return cmd.Run()
}

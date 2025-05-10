package platform

type ShortcutCreator interface {
	Create(id, command, binding string) error
}

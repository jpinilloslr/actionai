package platform

type ShortcutsMgr interface {
	Create(id, command, binding string) error
}

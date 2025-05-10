package gnome

import (
	"fmt"
	"os/exec"
	"strings"
)

type ShortcutCreator struct {
}

func NewShortcutCreator() *ShortcutCreator {
	return &ShortcutCreator{}
}

func (s *ShortcutCreator) Create(id, command, binding string) error {
	out, err := exec.Command("gsettings",
		"get",
		"org.gnome.settings-daemon.plugins.media-keys",
		"custom-keybindings",
	).Output()
	if err != nil {
		return fmt.Errorf("failed to get existing custom shortcuts: %w", err)
	}
	list := strings.TrimSpace(string(out))

	var updated string
	path := fmt.Sprintf("/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/%s/", id)
	if list == "@as []" || list == "[]" {
		updated = fmt.Sprintf("['%s']", path)
	} else {
		inner := strings.TrimSuffix(strings.TrimPrefix(list, "["), "]")
		updated = fmt.Sprintf("[%s, '%s']", inner, path)
	}

	if err := exec.Command("gsettings",
		"set",
		"org.gnome.settings-daemon.plugins.media-keys",
		"custom-keybindings",
		updated,
	).Run(); err != nil {
		return fmt.Errorf("set list: %w", err)
	}

	schema := fmt.Sprintf(
		"org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:%s",
		path,
	)
	for key, val := range map[string]string{
		"name":    id,
		"command": command,
		"binding": binding,
	} {
		if err := exec.Command("gsettings",
			"set", schema, key, val,
		).Run(); err != nil {
			return fmt.Errorf("set %s: %w", key, err)
		}
	}
	return nil
}

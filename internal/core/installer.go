package core

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type installer struct {
	logger    *slog.Logger
	cmdRepo   *commandRepo
	shortcuts platform.ShortcutCreator
}

func newInstaller(
	logger *slog.Logger,
	cmdRepo *commandRepo,
	shortcuts platform.ShortcutCreator,
) *installer {
	return &installer{
		logger:    logger,
		cmdRepo:   cmdRepo,
		shortcuts: shortcuts,
	}
}

func (i *installer) Install() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	cmds := i.cmdRepo.GetAll()

	for id, cmd := range cmds {
		command := fmt.Sprintf("%s --command %s", execPath, id)
		if err := i.shortcuts.Create(id, command, cmd.Shortcut); err != nil {
			i.logger.Error("failed to set shortcut", "id", id, "err", err)
			continue
		}
		i.logger.Info("installed shortcut", "id", id, "shortcut", cmd.Shortcut)
	}

	return nil
}

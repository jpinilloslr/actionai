package core

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type installer struct {
	actionRepo   *actionRepo
	logger       *slog.Logger
	shortcutsMgr platform.ShortcutsMgr
}

func newInstaller(
	logger *slog.Logger,
	actionRepo *actionRepo,
	shortcutsMgr platform.ShortcutsMgr,
) *installer {
	return &installer{
		logger:       logger,
		actionRepo:   actionRepo,
		shortcutsMgr: shortcutsMgr,
	}
}

func (i *installer) Install() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	actions := i.actionRepo.GetAll()

	for id, action := range actions {
		command := fmt.Sprintf("%s --action %s", execPath, id)
		if err := i.shortcutsMgr.Create(id, command, action.Shortcut); err != nil {
			i.logger.Error("failed to set shortcut", "id", id, "err", err)
			continue
		}
		i.logger.Info("installed shortcut", "id", id, "shortcut", action.Shortcut)
	}

	return nil
}

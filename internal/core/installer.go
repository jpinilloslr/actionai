package core

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jpinilloslr/actionai/internal/core/platform"
)

type Installer struct {
	actionRepo   *actionRepo
	logger       *slog.Logger
	shortcutsMgr platform.ShortcutsMgr
}

func NewInstaller(
	logger *slog.Logger,
	assetsMgr *AssetsMgr,
	shortcutsMgr platform.ShortcutsMgr,
) (*Installer, error) {
	actionRepo, err := newActionRepo(logger, assetsMgr.ActionsFile())
	if err != nil {
		return nil, err
	}

	return &Installer{
		logger:       logger,
		actionRepo:   actionRepo,
		shortcutsMgr: shortcutsMgr,
	}, nil
}

func (i *Installer) Install() error {
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

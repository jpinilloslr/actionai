package core

import (
	"fmt"
	"os"
	"path/filepath"
)

type AssetsMgr struct {
	workdir string
}

func NewAssetsMgr() (*AssetsMgr, error) {
	w := AssetsMgr{}

	if err := w.init(); err != nil {
		return nil, err
	}

	return &w, nil
}

func (w *AssetsMgr) ActionsFile() string {
	return filepath.Join(w.workdir, actionsFile)
}

func (w *AssetsMgr) LogsFile() string {
	return filepath.Join(w.workdir, logsFile)
}

func (w *AssetsMgr) SoundFile() string {
	return filepath.Join(w.workdir, soundFile)
}

func (w *AssetsMgr) init() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir = filepath.Join(dir, "actionai")

	if w.isAssetsMgr(dir) {
		w.setDir(dir)
		return nil
	}

	dir, err = os.Getwd()
	if err != nil {
		return err
	}

	if w.isAssetsMgr(dir) {
		w.setDir(dir)
		return nil
	}

	exeFile, err := os.Executable()
	if err != nil {
		return err
	}

	dir = filepath.Dir(exeFile)

	if w.isAssetsMgr(dir) {
		w.setDir(dir)
		return nil
	}

	return fmt.Errorf("could not resolve directory")
}

func (w *AssetsMgr) isAssetsMgr(dir string) bool {
	testFile := filepath.Join(dir, actionsFile)
	if _, err := os.Stat(testFile); err != nil {
		return false
	}
	return true
}

func (w *AssetsMgr) setDir(dir string) {
	fmt.Printf("Using directory %s\n", dir)
	w.workdir = dir
}

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

func (w *AssetsMgr) SoundFile() string {
	return filepath.Join(w.workdir, soundFile)
}

func (w *AssetsMgr) init() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir = filepath.Join(dir, "actionai")
	w.setDir(dir)
	return nil
}

func (w *AssetsMgr) setDir(dir string) {
	fmt.Printf("Using directory %s\n", dir)
	w.workdir = dir
}

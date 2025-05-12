package core

import (
	"fmt"
	"os"
	"path/filepath"
)

type WorkDir struct {
	dir string
}

func NewWorkDir() (*WorkDir, error) {
	w := WorkDir{}

	if err := w.init(); err != nil {
		return nil, err
	}

	return &w, nil
}

func (w *WorkDir) ActionsFile() string {
	return filepath.Join(w.dir, actionsFile)
}

func (w *WorkDir) LogsFile() string {
	return filepath.Join(w.dir, logsFile)
}

func (w *WorkDir) init() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir = filepath.Join(dir, "actionai")

	if w.isWorkDir(dir) {
		w.setDir(dir)
		return nil
	}

	dir, err = os.Getwd()
	if err != nil {
		return err
	}

	if w.isWorkDir(dir) {
		w.setDir(dir)
		return nil
	}

	exeFile, err := os.Executable()
	if err != nil {
		return err
	}

	dir = filepath.Dir(exeFile)

	if w.isWorkDir(dir) {
		w.setDir(dir)
		return nil
	}

	return fmt.Errorf("could not resolve directory")
}

func (w *WorkDir) isWorkDir(dir string) bool {
	testFile := filepath.Join(dir, actionsFile)
	if _, err := os.Stat(testFile); err != nil {
		return false
	}
	return true
}

func (w *WorkDir) setDir(dir string) {
	fmt.Printf("Using directory %s\n", dir)
	w.dir = dir
}

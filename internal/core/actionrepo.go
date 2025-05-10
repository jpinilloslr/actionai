package core

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type actionRepo struct {
	data        map[string]action
	logger      *slog.Logger
	actionsFile string
}

func newActionRepo(logger *slog.Logger, actionsFile string) (*actionRepo, error) {
	r := actionRepo{
		logger:      logger,
		actionsFile: actionsFile,
		data:        make(map[string]action),
	}

	if err := r.load(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *actionRepo) GetAll() map[string]action {
	return r.data
}

func (r *actionRepo) GetById(id string) (*action, error) {
	if action, ok := r.data[id]; ok {
		return &action, nil
	}
	return nil, fmt.Errorf("action %s not found", id)
}

func (r *actionRepo) load() error {
	fileName := r.actionsFile

	r.logger.Info("Loading actions from file", "file", fileName)
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var actions map[string]action
	err = json.Unmarshal(file, &actions)
	if err != nil {
		return err
	}

	if len(actions) == 0 {
		return fmt.Errorf("%s is empty", fileName)
	}

	r.data = actions
	return nil
}

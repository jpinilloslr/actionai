package core

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type commandRepo struct {
	data         map[string]command
	logger       *slog.Logger
	commandsFile string
}

func newCommandRepo(logger *slog.Logger, commandsFile string) (*commandRepo, error) {
	r := commandRepo{
		logger:       logger,
		commandsFile: commandsFile,
		data:         make(map[string]command),
	}

	if err := r.load(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *commandRepo) GetAll() map[string]command {
	return r.data
}

func (r *commandRepo) GetById(id string) (*command, error) {
	if cmd, ok := r.data[id]; ok {
		return &cmd, nil
	}
	return nil, fmt.Errorf("command %s not found", id)
}

func (r *commandRepo) load() error {
	fileName := r.commandsFile

	r.logger.Info("Loading commands from file", "file", fileName)
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var cmds map[string]command
	err = json.Unmarshal(file, &cmds)
	if err != nil {
		return err
	}

	if len(cmds) == 0 {
		return fmt.Errorf("%s is empty", fileName)
	}

	r.data = cmds
	return nil
}

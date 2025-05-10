package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jpinilloslr/actionai/internal/core"
	"github.com/jpinilloslr/actionai/internal/gnome"
	"github.com/jpinilloslr/actionai/internal/openai"
)

func main() {
	command, installShortcuts := getOptions()

	workDir, err := core.NewWorkDir()
	if err != nil {
		fmt.Printf("Error resolving working directory: %v\n", err)
	}

	logger, err := core.NewLogger(workDir.LogsFile())
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	model, err := openai.New(logger)
	if err != nil {
		logger.Error("Error initializing model", "error", err)
		os.Exit(1)
	}

	runner, err := core.NewAiModelRunner(
		logger,
		workDir,
		model,
		gnome.NewDialog(),
		gnome.NewNotifier(),
		gnome.NewClipboard(),
		gnome.NewScreenshotter(),
		gnome.NewSpeechRecorder(),
		gnome.NewSelTextProvider(),
		gnome.NewShortcutCreator(),
	)

	if err != nil {
		logger.Error("Error creating the model runner", "error", err)
		os.Exit(1)
	}

	if installShortcuts {
		err := runner.InstallShortcuts()
		if err != nil {
			logger.Error("Error installing shortcuts", "error", err)
			os.Exit(1)
		}
		return
	}

	if command == "" {
		logger.Error("No command provided")
		os.Exit(1)
	}

	err = runner.Run(command)
	if err != nil {
		logger.Error("Error running the model", "error", err)
		os.Exit(1)
	}
}

func getOptions() (string, bool) {
	command := flag.String("command", "", "Command to send to the model")
	installShortcuts := flag.Bool("install-shortcuts", false, "Install shortcuts for the model")
	flag.Parse()
	return *command, *installShortcuts
}

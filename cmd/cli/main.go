package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jpinilloslr/actionai/internal/core"
	"github.com/jpinilloslr/actionai/internal/gnome"
	"github.com/jpinilloslr/actionai/internal/openai"
)

func main() {
	actionId, install := getOptions()

	assetsMgr, err := core.NewAssetsMgr()
	if err != nil {
		fmt.Printf("Error resolving working directory: %v\n", err)
	}

	logger, err := core.NewLogger(assetsMgr.LogsFile())
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	model, err := openai.NewAIModel(logger)
	if err != nil {
		logger.Error("Error initializing model", "error", err)
		os.Exit(1)
	}

	voiceEngine, err := openai.NewVoiceEngine(logger)
	if err != nil {
		logger.Error("Error initializing voice engine", "error", err)
		os.Exit(1)
	}

	runner, err := core.NewActionRunner(
		logger,
		assetsMgr,
		model,
		voiceEngine,
		gnome.NewDialog(),
		gnome.NewNotifier(),
		gnome.NewClipboard(),
		gnome.NewAudioPlayer(),
		gnome.NewScreenshotter(),
		gnome.NewVoiceRecorder(),
		gnome.NewSelTextProvider(),
	)

	if err != nil {
		logger.Error("Error creating the model runner", "error", err)
		os.Exit(1)
	}

	if install {
		installer, err := core.NewInstaller(logger, assetsMgr, gnome.NewShortcutsMgr())
		if err != nil {
			logger.Error("Error creating the installer", "error", err)
			os.Exit(1)
		}

		err = installer.Install()
		if err != nil {
			logger.Error("Error installing shortcuts", "error", err)
			os.Exit(1)
		}
		return
	}

	if actionId == "" {
		logger.Error("No action provided")
		os.Exit(1)
	}

	err = runner.RunFromActionRepo(context.Background(), actionId)
	if err != nil {
		logger.Error("Error running the model", "error", err)
		os.Exit(1)
	}
}

func getOptions() (string, bool) {
	actionId := flag.String("action", "", "Action ID to run")
	install := flag.Bool("install", false, "Install shortcuts for actions")
	flag.Parse()

	return *actionId, *install
}

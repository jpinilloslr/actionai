package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jpinilloslr/actionai/internal/core"
	"github.com/jpinilloslr/actionai/internal/gnome"
	"github.com/jpinilloslr/actionai/internal/openai"
	"github.com/spf13/cobra"
)

var (
	inputs       []string
	output       string
	instructions string
	model        string
)

var validInputs = []string{
	"screen",
	"screen-section",
	"selected-text",
	"voice",
	"window",
}

var validOutputs = []string{
	"clipboard",
	"stdout",
	"voice",
	"window",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an action",
	Long:  "Run an action",
	RunE: func(cmd *cobra.Command, args []string) error {
		parsedInputs, err := core.ParseInputList(inputs)
		if err != nil {
			return err
		}

		parsedOutput, err := core.ParseOutput(output)
		if err != nil {
			return err
		}

		action := core.Action{
			Model:        model,
			Inputs:       parsedInputs,
			Output:       parsedOutput,
			Instructions: instructions,
		}
		fmt.Printf(
			"Running action: %v\n",
			action,
		)

		return runAction(&action)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceVarP(
		&inputs,
		"input",
		"i",
		nil,
		fmt.Sprintf("List of input types (%s)", strings.Join(validInputs, ", ")),
	)
	runCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"",
		fmt.Sprintf("Output type (%s)", strings.Join(validOutputs, ", ")),
	)
	runCmd.Flags().StringVarP(&model, "model", "m", "gpt-4.1-mini", "AI model to use")
	runCmd.Flags().StringVarP(&instructions, "instructions", "n", "", "Instructions to pass to the AI model")

	runCmd.MarkFlagRequired("input")
	runCmd.MarkFlagRequired("output")
}

func runAction(action *core.Action) error {
	logger := slog.Default()
	assetsMgr, err := core.NewAssetsMgr()
	if err != nil {
		return fmt.Errorf("Error resolving working directory: %v", err)
	}

	model, err := openai.NewAIModel(logger)
	if err != nil {
		return fmt.Errorf("Error initializing model: %v", err)
	}

	voiceEngine, err := openai.NewVoiceEngine(logger)
	if err != nil {
		return fmt.Errorf("Error initializing voice engine: %v", err)
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
		return fmt.Errorf("Error creating the model runner: %v", err)
	}

	err = runner.RunFromAction(context.Background(), action)
	if err != nil {
		return fmt.Errorf("Error running the model: %v", err)
	}

	return nil
}

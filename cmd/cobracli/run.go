package main

import (
	"fmt"
	"slices"
	"strings"

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
		for _, input := range inputs {
			if !slices.Contains(validInputs, input) {
				return fmt.Errorf("invalid input: %s", input)
			}
		}

		if !slices.Contains(validOutputs, output) {
			return fmt.Errorf("invalid output: %s", output)
		}

		fmt.Println("Running with:")
		fmt.Println("Inputs:", inputs)
		fmt.Println("Output:", output)
		fmt.Println("Instructions:", instructions)
		fmt.Println("Model:", model)

		return nil
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

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "actionai",
	Short: "ActionAI is a command line tool for executing AI actions",
	Long:  "ActionAI is a flexible command line tool designed to execute AI-driven actions seamlessly.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

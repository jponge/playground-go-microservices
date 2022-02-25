package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:     "temperature-generator",
	Short:   "A temperature updates generator",
	Version: "0.1",
}

var runStarted = false

func Execute() error {
	return rootCmd.Execute()
}

func RunStarted() bool {
	return runStarted
}

func init() {
	// We could to Viper init work here as well
}

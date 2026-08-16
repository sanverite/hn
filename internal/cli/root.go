package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "hn",
	Short: "Hacker News from your terminal",
	Long: "A fast, clean Hacker News client for the terminal.",
}

func Execute() error {
	return rootCmd.Execute()
}

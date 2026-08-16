package cli

import (
	"fmt"

	"github.com/sanverite/hn/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("git sha: %s\n", version.GitSHA)
		fmt.Printf("built at: %s\n", version.BuildTime)
		fmt.Printf("go version: %s\n", version.GoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

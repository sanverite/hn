package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sanverite/hn/internal/api"
	"github.com/sanverite/hn/internal/formatter"
	"github.com/spf13/cobra"
)

var newLimit int

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Show new stories from Hacker News.",
	RunE:  runNew,
}

func init() {
	newCmd.Flags().IntVarP(&newLimit, "limit", "l", 10, "number of stories to fetch")
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
	if newLimit < 1 || newLimit > 500 {
		return fmt.Errorf("limit must be between 1 and 500, got: %d", newLimit)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := api.New()

	stories, err := client.TopStories(ctx, "newstories", newLimit)
	if err != nil {
		return fmt.Errorf("fetching stories: %w", err)
	}

	formatter.PrintStories(os.Stdout, stories)
	return nil
}

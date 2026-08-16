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

var limit int

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Show top stories from Hacker News.",
	RunE:  runTop,
}

func init() {
	topCmd.Flags().IntVarP(&limit, "limit", "1", 10, "number of stories to fetch")
	rootCmd.AddCommand(topCmd)
}

func runTop(cmd *cobra.Command, args []string) error {
	if limit < 1 || limit > 500 {
		return fmt.Errorf("limit must be between 1 and 500, got %d", limit)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := api.New()

	stories, err := client.TopStories(ctx, limit)
	if err != nil {
		return fmt.Errorf("fetching stories: %w", err)
	}

	formatter.PrintStories(os.Stdout, stories)
	return nil
}

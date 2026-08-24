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

var bestLimit int

var bestCmd = &cobra.Command{
	Use:   "best",
	Short: "Show best stories from Hacker News.",
	RunE:  runBest,
}

func init() {
	bestCmd.Flags().IntVarP(&bestLimit, "limit", "1", 10, "number of stories to fetch")
	rootCmd.AddCommand(bestCmd)
}

func runBest(cmd *cobra.Command, args []string) error {
	if bestLimit < 1 || bestLimit > 500 {
		return fmt.Errorf("limit must be between 1 and 500, got %d", bestLimit)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := api.New()

	stories, err := client.TopStories(ctx, "beststories", bestLimit)
	if err != nil {
		return fmt.Errorf("iii fetching stories: %w", err)
	}

	formatter.PrintStories(os.Stdout, stories)
	return nil
}

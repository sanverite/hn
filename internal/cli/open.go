package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/sanverite/hn/internal/api"
	"github.com/sanverite/hn/internal/browser"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <id>",
	Short: "Open a story in your browser.",
	Args:  cobra.ExactArgs(1),
	RunE:  runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) error {
	id, err := parseID(args[0])
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := api.New()

	story, err := client.FetchStory(ctx, id)
	if err != nil {
		return fmt.Errorf("fetching story %d: %w", id, err)
	}

	url := story.URL
	if url == "" {
		// Text posts and Ask HN have no URL - fallback to discussions page
		url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ID)
	}

	fmt.Printf("Opening: %s\n", url)

	if err := browser.Open(url); err != nil {
		return fmt.Errorf("opening browser: %w", err)
	}

	return nil
}

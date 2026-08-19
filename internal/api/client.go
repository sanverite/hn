package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const baseURL = "https://hacker-news.firebaseio.com/v0"

// Client handles all the communication with the HN API.
// It is the only place in the codebase that knows about HTTP.
type Client struct {
	http    *http.Client
	baseURL string
}

// New returns a Client ready for use.
// Timeout is set deliberately - we never make a network call
// without a timeout in production code.
func New() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

// NewWithBaseURL creates a Client pointed at a custom base URL.
// Used in tests to inject a fake server.
func NewWithBaseURL(baseURL string) *Client {
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Story represents a single HN story.
type Story struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Score int    `json:"score"`
	By    string `json:"by"`
	Time  int64  `json:"time"`
}

// TopStories returns the top n story IDs from HN.
func (c *Client) TopStories(ctx context.Context, limit int) ([]Story, error) {
	ids, err := c.fetchIDs(ctx, "topstories")
	if err != nil {
		return nil, fmt.Errorf("fetching top stories: %w", err)
	}

	if limit > len(ids) {
		limit = len(ids)
	}

	return c.fetchStories(ctx, ids[:limit])
}

// fetchIDs retrieves a list of story IDs for a given feed.
func (c *Client) fetchIDs(ctx context.Context, feed string) ([]int, error) {
	url := fmt.Sprintf("%s/%s.json", c.baseURL, feed)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return ids, nil
}

// fetchStories retrieves full story details for a list of IDs.
func (c *Client) fetchStories(ctx context.Context, ids []int) ([]Story, error) {
	// Make a slice upfront with exact capacity we need.
	// This is important - we're going to write into specific indexes
	// from multiple goroutines. Pre-allocating means no slice grows,
	// no hidden allocations, and each goroutine writes to its own index
	// so there's no data race.
	stories := make([]Story, len(ids))

	// Buffered channel with capacity 5 acts as a semaphore.
	// Only 5 goroutines can hold a token at once.
	// The rest block on the send until a slot opens up.
	sem := make(chan struct{}, 5)

	// errgroup gives us two things:
	// - g.Go() launches a goroutine and tracks it
	// - g.Wait() blocks until all goroutines finish and returns
	// 	 the first non-nil error any of them returned
	//
	// ctx passed into errgroup.WithContext means: if any goroutine
	// fails, the context is cancelled and all other goroutines
	// should stop early when they next check ctx.
	g, ctx := errgroup.WithContext(ctx)

	for i, id := range ids {
		g.Go(func() error {
			// Acquire a token - blocks if 5 goroutines are already running.
			// struct{} is used because it takes zero bytes of memory.
			// We only care about the slot, not the value.
			sem <- struct{}{}

			// Release the token when this goroutine finishes,
			// regardless of success of failure.
			// defer runs even if we return an error.
			defer func() { <-sem }()

			// Each goroutine fetches one story.
			// If ctx is already cancelled (another goroutine failed),
			// FetchStory returns immediately with a context error.
			story, err := c.FetchStory(ctx, id)
			if err != nil {
				// Wrap with story ID so we know which one failed
				return fmt.Errorf("story %d: %w", id, err)
			}

			// Write to a specific index - no two goroutines
			// ever write to the same index, so this is safe
			// without a mutex.
			stories[i] = story
			return nil
		})
	}

	// Wait blocks until every g.Go goroutine has returned.
	// If any returned a non-nil error, Wait returns that error.
	// The other goroutines are not killed - they run to completion,
	// but their results are discarded since we're returning an error.
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("fetching stories: %w", err)
	}

	return stories, nil
}

// fetchStory retrieves a single Story by ID.
func (c *Client) FetchStory(ctx context.Context, id int) (Story, error) {
	url := fmt.Sprintf("%s/item/%d.json", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Story{}, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Story{}, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Story{}, fmt.Errorf("unexpected status %d:", resp.StatusCode)
	}

	var story Story
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return Story{}, fmt.Errorf("decoding story: %w", err)
	}

	return story, nil
}

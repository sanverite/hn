package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	stories := make([]Story, 0, len(ids))

	for _, id := range ids {
		story, err := c.FetchStory(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("fetching story %d: %w", id, err)
		}
		stories = append(stories, story)
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

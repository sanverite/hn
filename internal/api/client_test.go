package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanverite/hn/internal/api"
)

func TestTopStories(t *testing.T) {
	// We spin up a fake HTTP server so tests never hit the real API.
	// Real API calls in tests = flaky tests. Flaky tests = useless tests.

	ids := []int{1, 2}
	stories := map[int]api.Story{
		1: {ID: 1, Title: "Story One", Score: 100, By: "user1"},
		2: {ID: 2, Title: "Story Two", Score: 200, By: "user2"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v0/topstories.json":
			json.NewEncoder(w).Encode(ids)
		case "/v0/item/1.json":
			json.NewEncoder(w).Encode(stories[1])
		case "/v0/item/2.json":
			json.NewEncoder(w).Encode(stories[2])
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewWithBaseURL(server.URL + "/v0")

	got, err := client.TopStories(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 stories, got %d", len(got))
	}

	if got[0].Title != "Story One" {
		t.Errorf("expected 'Story One', got %q", got[0].Title)
	}
}

func TestTopStoriesLimitCapped(t *testing.T) {
	ids := []int{1}
	story := api.Story{ID: 1, Title: "Only Story", Score: 50, By: "user1"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v0/topstories.json":
			json.NewEncoder(w).Encode(ids)
		case "/v0/item/1.json":
			json.NewEncoder(w).Encode(story)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := api.NewWithBaseURL(server.URL + "/v0")

	// Request 10 but only 1 exists — should return 1, not error
	got, err := client.TopStories(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Errorf("expected 1 story, got %d", len(got))
	}
}

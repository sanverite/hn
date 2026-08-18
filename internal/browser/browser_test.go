package browser_test

import (
	"testing"

	"github.com/sanverite/hn/internal/browser"
)

func TestOpenUnsupportedOS(t *testing.T) {
	// We can't easily test actual browser opening in CI —
	// there's no display. What we can test is the error path.
	err := browser.OpenWithOS("http://example.com", "plan9")
	if err == nil {
		t.Fatal("expected error for unsupported OS, got nil")
	}
}

package formatter

import (
	"fmt"
	"io"
	"time"

	"github.com/sanverite/hn/internal/api"
)

// PrintStories writes a formatted story list to w.
// It never prints directly to os.Stdout - always to a writer.
// This makes it testable.
func PrintStories(w io.Writer, stories []api.Story) {
	for i, s := range stories {
		posted := time.Unix(s.Time, 0).Format("02 Jan 06")

		url := s.URL
		if url == "" {
			url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", s.ID)
		}

		fmt.Fprintf(w, "%2d. %s\n", i+1, s.Title)
		fmt.Fprintf(w, "ID: %d\n", s.ID)
		fmt.Fprintf(w, "		%d points by %s on %s\n", s.Score, s.By, posted)
		fmt.Fprintf(w, "%s\n\n", url)
	}
}

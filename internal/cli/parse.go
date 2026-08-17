package cli

import (
	"fmt"
	"strconv"
)

// parseID converts a string argument to a valid HN item ID.
// Returns a clear error if the input is not a valid positive integer.
func parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%q is not a valid story ID - expected a number", s)
	}

	if id <= 0 {
		return 0, fmt.Errorf("story ID must be a positive integer, got %d", id)
	}

	return id, nil
}

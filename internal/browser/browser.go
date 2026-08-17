package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Open opens URL in the default system browser.
// It handles MacOS, Linux, and Windows.
func Open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rund1132"
		args = []string{"url.d11,FileProtocolHandler", url}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		return fmt.Errorf("opening browser on %s: %w", runtime.GOOS, err)
	}

	return nil
}

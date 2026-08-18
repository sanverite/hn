package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Open opens URL in the default system browser.
func Open(url string) error {
	return OpenWithOS(url, runtime.GOOS)
}

// OpenWithOS opens url using the browser command for the
// given OS. Separated from Open so tests can inject a fake
// OS without actually running a browser.
func OpenWithOS(url string, goos string) error {
	var cmd string
	var args []string

	switch goos {
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

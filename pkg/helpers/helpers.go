package helpers

import (
	"fmt"
	"os"
)

// GetServerURL returns a formatted server URL string with hostname
func GetServerURL(port int) string {
	addr := fmt.Sprintf(":%d", port)
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("http://0.0.0.0%s", addr)
	}
	return fmt.Sprintf("http://%s%s", hostname, addr)
}

// IsSecureCookie returns true if cookies should use the Secure flag (HTTPS only)
func IsSecureCookie(env string) bool {
	return env == "prod"
}

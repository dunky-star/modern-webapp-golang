package main

import (
	"fmt"
	"os"
)

// getServerURL returns a formatted server URL string with hostname
func getServerURL(port int) string {
	addr := fmt.Sprintf(":%d", port)
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("http://0.0.0.0%s", addr)
	}
	return fmt.Sprintf("http://%s%s", hostname, addr)
}

// isSecureCookie returns true if cookies should use the Secure flag (HTTPS only)
func isSecureCookie(env string) bool {
	return env == "prod"
}

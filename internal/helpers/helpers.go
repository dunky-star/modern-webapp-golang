package helpers

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/dunky-star/modern-webapp-golang/internal/config"
)

var app *config.AppConfig

// GetServerURL returns a formatted server URL string with hostname
func GetServerURL(port int) string {
	addr := fmt.Sprintf(":%d", port)
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("http://0.0.0.0%s", addr)
	}
	return fmt.Sprintf("http://%s%s", hostname, addr)
}

func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError writes a client error response (4xx status codes)
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Printf("Client error with status of %d", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError writes a server error response (500) and logs the error with stack trace
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Printf("ERROR\t %s", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dunky-star/modern-webapp-golang/pkg/helpers"
)

const version = "1.0.0"

// Config, to allow the server to be configured at startup dynamically
type config struct {
	port int
	env  string
}

type application struct {
	config  config
	logger  *log.Logger
	session *scs.SessionManager
}

var startTime = time.Now()

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.Parse()

	// General application logger (stdout)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config:  cfg,
		logger:  logger,
		session: newSessionManager(cfg.env),
	}

	// Create the HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),     // Set the routes for the server
		ReadTimeout:  10 * time.Second, // Maximum duration for reading the entire request, including the body
		WriteTimeout: 30 * time.Second, // Maximum duration before timing out writes of the response
		IdleTimeout:  time.Minute,      // Maximum amount of time to wait for the next request when keep-alives are enabled
	}

	// Ensure request logger is closed on shutdown
	defer closeRequestLogger()

	logger.Printf("Server is running on %s\n", helpers.GetServerURL(cfg.port))

	// Start the server and log any error if it fails
	err := srv.ListenAndServe()
	logger.Fatal(err)
}

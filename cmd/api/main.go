package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/config"
	"github.com/dunky-star/modern-webapp-golang/internal/data"
	"github.com/dunky-star/modern-webapp-golang/internal/handlers"
	"github.com/dunky-star/modern-webapp-golang/internal/render"
	"github.com/dunky-star/modern-webapp-golang/pkg/helpers"
)

const appVersion = "1.0.0"

var app config.AppConfig
var appStartTime = time.Now()

func main() {
	var port int
	var env string
	flag.IntVar(&port, "port", 3000, "API server port")
	flag.StringVar(&env, "env", "dev", "Environment (dev|stage|prod)")
	flag.Parse()

	err := run(port, env)
	if err != nil {
		app.InfoLog.Fatal(err)
	}

	app.InfoLog.Printf("Server is running on %s\n", helpers.GetServerURL(port))

	// Create the HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes(),         // Set the routes for the server
		ReadTimeout:  10 * time.Second, // Maximum duration for reading the entire request, including the body
		WriteTimeout: 30 * time.Second, // Maximum duration before timing out writes of the response
		IdleTimeout:  time.Minute,      // Maximum amount of time to wait for the next request when keep-alives are enabled
	}

	// Ensure request logger is closed on shutdown
	defer closeRequestLogger()

	// Start the server and log any error if it fails
	err = srv.ListenAndServe()
	app.InfoLog.Fatal(err)
}

func run(port int, env string) error {
	// Register types for session storage
	gob.Register(data.Reservation{})

	// Initialize application configuration
	cfg := config.New(port, env)

	// Create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		cfg.InfoLog.Fatal("cannot create template cache")
		return err
	}

	// Set template cache and use cache flag
	cfg.TemplateCache = tc
	cfg.UseCache = (cfg.Env != "dev")

	app = *cfg

	// Initialize handlers repository
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// Initialize render package with app config
	render.NewTemplates(&app)

	// Store appStartTime in handlers package for health check
	handlers.SetAppStartTime(appStartTime)
	handlers.SetAppVersion(appVersion)

	return nil
}

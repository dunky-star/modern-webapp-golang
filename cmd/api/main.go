package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dunky-star/modern-webapp-golang/internal/config"
	"github.com/dunky-star/modern-webapp-golang/internal/data"
	"github.com/dunky-star/modern-webapp-golang/internal/driver"
	"github.com/dunky-star/modern-webapp-golang/internal/handlers"
	"github.com/dunky-star/modern-webapp-golang/internal/helpers"
	"github.com/dunky-star/modern-webapp-golang/internal/render"
	"github.com/joho/godotenv"
)

const appVersion = "1.0.0"

var app config.AppConfig
var appStartTime = time.Now()

func main() {
	var port int
	var env string
	var dsn string
	godotenv.Load(".env")
	flag.IntVar(&port, "port", 3000, "API server port")
	flag.StringVar(&env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&dsn, "db-dsn", os.Getenv("DB_DSN"), "DB connection string")
	flag.Parse()

	err := run(port, env, dsn)
	if err != nil {
		app.ErrorLog.Fatal(err)
	}

	// Close database connection when application exits
	defer driver.Close()
	// Close email channel when application exits
	defer close(app.MailChan)
	// Listen for mail messages using goroutine to send emails
	listerForMail()

	msg := data.MailData{
		To:      "dunky@do.us",
		From:    "me@here.com",
		Subject: "Test email",
		Content: "",
	}
	app.MailChan <- msg

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
	app.ErrorLog.Fatal(err)
}

func run(port int, env string, dsn string) error {
	// Create a channel for sending emails
	mailChan := make(chan data.MailData)
	app.MailChan = mailChan

	// Initialize application configuration
	cfg := config.New(port, env, dsn)

	// Create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		cfg.ErrorLog.Fatal("cannot create template cache")
		return err
	}

	// Set template cache and use cache flag
	cfg.TemplateCache = tc
	cfg.UseCache = (cfg.Env != "dev")

	app = *cfg

	// Validate DSN is set
	if cfg.DSN == "" {
		cfg.ErrorLog.Fatal("db-dsn flag or DB_DSN environment variable must be set")
	}

	// Connect to database with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := driver.Init(ctx, cfg.DSN)
	if err != nil {
		cfg.ErrorLog.Fatal(err)
	}
	cfg.InfoLog.Println("Database connection pool established successfully")

	// Initialize handlers repository
	repo := handlers.NewRepo(&app, dbPool)
	handlers.NewHandlers(repo)

	// Initialize render package with app config
	render.NewRender(&app)

	// Initialize helpers package with app config
	helpers.NewHelpers(&app)

	// Store appStartTime in handlers package for health check
	handlers.SetAppStartTime(appStartTime)
	handlers.SetAppVersion(appVersion)

	return nil
}

package dbrepo

import (
	"github.com/dunky-star/modern-webapp-golang/internal/config"
	"github.com/dunky-star/modern-webapp-golang/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConnection struct {
	App *config.AppConfig
	DB  *pgxpool.Pool
}

func NewDBConnection(conn *pgxpool.Pool, app *config.AppConfig) repository.DatabaseConn {
	return &DBConnection{
		App: app,
		DB:  conn,
	}
}

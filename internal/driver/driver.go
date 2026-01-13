package driver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool    *pgxpool.Pool
	once    sync.Once
	initErr error
)

// Init initializes the global connection pool using the provided DSN and returns it.
// Subsequent calls reuse the same pool.
func Init(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn must not be empty")
	}

	once.Do(func() {
		cfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			initErr = fmt.Errorf("unable to parse config: %w", err)
			return
		}

		cfg.MaxConns = 5
		cfg.MinConns = 2
		cfg.MaxConnLifetime = time.Hour
		cfg.HealthCheckPeriod = 30 * time.Second

		p, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			initErr = fmt.Errorf("unable to create pool: %w", err)
			return
		}

		if err := p.Ping(ctx); err != nil {
			p.Close()
			initErr = fmt.Errorf("database ping failed: %w", err)
			return
		}

		var now time.Time
		if err := p.QueryRow(ctx, "SELECT now()").Scan(&now); err != nil {
			p.Close()
			initErr = fmt.Errorf("validation query failed: %w", err)
			return
		}

		fmt.Println("Connected, server time is:", now)
		pool = p
	})

	return pool, initErr
}

// Close closes the global pool if it has been initialized.
func Close() {
	if pool != nil {
		pool.Close()
	}
}

// Pool returns the initialized pool instance.
func Pool() *pgxpool.Pool {
	return pool
}

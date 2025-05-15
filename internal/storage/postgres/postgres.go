package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"test_golang_user_api/internal/config"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Postgres) (*Storage, error) {
	db, err := sql.Open("postgres", buildUri(cfg))

	if err != nil {
		return nil, fmt.Errorf("failed to configure connection to postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Storage{db: db}, nil
}

func buildUri(cfg config.Postgres) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)
}

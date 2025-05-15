package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"test_golang_user_api/internal/config"
	"time"
)

type User struct {
	ID        uuid.UUID
	Firstname string
	Lastname  string
	Email     string
	Age       int
	Created   time.Time
}

func NewUser(firstname, lastname, email string, age int) *User {
	return &User{
		ID:        uuid.New(),
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Age:       age,
		Created:   time.Now(),
	}
}

type Storage struct {
	db *sql.DB
}

func New(cfg config.Postgres) (*Storage, error) {
	db, err := sql.Open("postgres", buildUri(cfg))

	if err != nil {
		return nil, fmt.Errorf("failed to configure connection to postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	if err := runMigrations(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
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

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up error: %w", err)
	}

	return nil
}

func (s *Storage) CreateUser(user *User) error {
	query := `INSERT INTO users (id, firstname, lastname, email, age, created) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(query, user.ID, user.Firstname, user.Lastname, user.Email, user.Age, user.Created)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (s *Storage) GetUser(id uuid.UUID) (User, error) {
	query := `SELECT id, firstname, lastname, email, age, created FROM users WHERE id = $1`

	var user User
	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Age,
		&user.Created,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found: %w", err)
		}
		return user, fmt.Errorf("query failed: %w", err)
	}

	return user, nil
}

func (s *Storage) EditUser(user User) (User, error) {
	query := `UPDATE users SET firstname = $1, lastname = $2, email = $3, age = $4 WHERE id = $5`

	_, err := s.db.Exec(query, user.Firstname, user.Lastname, user.Email, user.Age, user.ID)

	if err != nil {
		return User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *Storage) DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

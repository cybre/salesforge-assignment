package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config contains the configuration for connecting to a database.
type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// NewPostgresDB creates a new Postgres database connection.
func NewPostgresDB(config Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Name))
	if err != nil {
		return nil, err
	}

	return db, nil
}

// RunMigrations brings the database up to the latest version.
func RunMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		MigrationsTable:       "schema_migrations",
		MultiStatementEnabled: true,
		MultiStatementMaxSize: postgres.DefaultMultiStatementMaxSize,
	})
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create Migrate instance: %w", err)
	}

	// Run migrations up to the latest one.
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

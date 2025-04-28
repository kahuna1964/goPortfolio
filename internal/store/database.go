package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func Open() (*sql.DB, error) {
	cfg := PostgresConfig{
		Host:     "localHost",
		Port:     "5432",
		User:     "goDBA",
		Password: "gogogo",
		Database: "goPortfolio",
		SSLMode:  "disable",
	}

	fmt.Printf("   Connecting to database [%s]...", cfg.Database)
	db, err := sql.Open("pgx", cfg.String()) // don't use in production
	if err != nil {
		return nil, fmt.Errorf("db:open %w", err)
	}

	err = db.Ping() // make sure we really did connect
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected\n")
	return db, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

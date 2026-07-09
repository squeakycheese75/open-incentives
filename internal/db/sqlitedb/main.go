package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func NewDB(ctx context.Context, cfg Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON;`); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.ExecContext(ctx, `PRAGMA journal_mode = WAL;`); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.ExecContext(ctx, `PRAGMA busy_timeout = 5000;`); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Running database migrations...")
	if err := Migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	switch cfg.Bootstrap.Mode {
	case "none":
		log.Println("Bootstrap disabled")

	case "auto":
		empty, err := IsEmpty(ctx, db)
		if err != nil {
			db.Close()
			return nil, err
		}

		if empty {
			log.Println("Bootstrapping empty database...")
			if err := Seed(ctx, db, cfg.Bootstrap); err != nil {
				db.Close()
				return nil, err
			}
		} else {
			log.Println("Database already initialized; skipping bootstrap")
		}

	case "force":
		log.Println("Force bootstrapping database...")
		if err := Seed(ctx, db, cfg.Bootstrap); err != nil {
			db.Close()
			return nil, err
		}

	default:
		db.Close()
		return nil, fmt.Errorf("invalid bootstrap mode %q: must be auto, none, or force", cfg.Bootstrap.Mode)
	}

	return db, nil
}

func IsEmpty(ctx context.Context, db *sql.DB) (bool, error) {
	var count int

	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM organizations
	`).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

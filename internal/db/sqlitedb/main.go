package sqlitedb

import (
	"context"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func NewDB(ctx context.Context, path string, seed bool) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
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

	if seed {
		log.Println("Seeding database...")
		if err := Seed(ctx, db); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

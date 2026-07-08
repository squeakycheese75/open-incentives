package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func Seed(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM organizations
		)
	`).Scan(&exists); err != nil {
		return fmt.Errorf("check organizations: %w", err)
	}

	if exists {
		return tx.Commit()
	}

	orgPublicID := "org_" + uuid.NewString()
	projectPublicID := "proj_" + uuid.NewString()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO organizations (public_id, name)
		VALUES (?, ?)
	`, orgPublicID, "Default Organization")
	if err != nil {
		return fmt.Errorf("create default organization: %w", err)
	}

	orgID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get organization id: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO projects (
			public_id,
			org_id,
			name
		)
		VALUES (?, ?, ?)
	`, projectPublicID, orgID, "Default Project")
	if err != nil {
		return fmt.Errorf("create default project: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

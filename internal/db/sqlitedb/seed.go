package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/squeakycheese75/open-incentives/configs"
	"golang.org/x/crypto/bcrypt"
)

func Seed(ctx context.Context, db *sql.DB, cfg configs.BootstrapConfig) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var orgID int64

	err = tx.QueryRowContext(ctx, `
		SELECT id
		FROM organizations
		WHERE public_id = ?
		  AND deleted_at IS NULL
	`, cfg.OrgSlug).Scan(&orgID)

	if err == sql.ErrNoRows {
		res, err := tx.ExecContext(ctx, `
			INSERT INTO organizations (public_id, name)
			VALUES (?, ?)
		`, cfg.OrgSlug, cfg.OrgName)
		if err != nil {
			return fmt.Errorf("create organization: %w", err)
		}

		orgID, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("get organization id: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("lookup organization: %w", err)
	}

	var projectExists bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM projects
			WHERE org_id = ?
			  AND name = ?
			  AND deleted_at IS NULL
		)
	`, orgID, cfg.ProjectName).Scan(&projectExists); err != nil {
		return fmt.Errorf("check project: %w", err)
	}

	if !projectExists {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO projects (public_id, org_id, name)
			VALUES (?, ?, ?)
		`, "proj_"+uuid.NewString(), orgID, cfg.ProjectName)
		if err != nil {
			return fmt.Errorf("create project: %w", err)
		}
	}

	var userExists bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE org_id = ?
			  AND email = ?
			  AND deleted_at IS NULL
		)
	`, orgID, cfg.AdminEmail).Scan(&userExists); err != nil {
		return fmt.Errorf("check admin user: %w", err)
	}

	if !userExists {
		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(cfg.AdminPassword),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return fmt.Errorf("hash admin password: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO users (
				public_id,
				org_id,
				email,
				password_hash,
				role
			)
			VALUES (?, ?, ?, ?, ?)
		`, "user_"+uuid.NewString(), orgID, cfg.AdminEmail, string(passwordHash), "admin")
		if err != nil {
			return fmt.Errorf("create admin user: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// func Seed(ctx context.Context, db *sql.DB, cfg configs.BootstrapConfig) error {
// 	tx, err := db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("begin transaction: %w", err)
// 	}
// 	defer tx.Rollback()

// 	var orgID int64

// 	err = tx.QueryRowContext(ctx, `
// 		SELECT id
// 		FROM organizations
// 		WHERE slug = ?
// 		  AND deleted_at IS NULL
// 	`, cfg.OrgSlug).Scan(&orgID)

// 	if err == sql.ErrNoRows {
// 		res, err := tx.ExecContext(ctx, `
// 			INSERT INTO organizations (public_id, name, slug)
// 			VALUES (?, ?, ?)
// 		`, "org_"+uuid.NewString(), cfg.OrgName, cfg.OrgSlug)
// 		if err != nil {
// 			return fmt.Errorf("create default organization: %w", err)
// 		}

// 		orgID, err = res.LastInsertId()
// 		if err != nil {
// 			return fmt.Errorf("get organization id: %w", err)
// 		}
// 	} else if err != nil {
// 		return fmt.Errorf("lookup organization: %w", err)
// 	}

// 	var projectExists bool
// 	if err := tx.QueryRowContext(ctx, `
// 		SELECT EXISTS(
// 			SELECT 1
// 			FROM projects
// 			WHERE org_id = ?
// 			  AND name = ?
// 			  AND deleted_at IS NULL
// 		)
// 	`, orgID, cfg.ProjectName).Scan(&projectExists); err != nil {
// 		return fmt.Errorf("check default project: %w", err)
// 	}

// 	if !projectExists {
// 		_, err = tx.ExecContext(ctx, `
// 			INSERT INTO projects (public_id, org_id, name)
// 			VALUES (?, ?, ?)
// 		`, "proj_"+uuid.NewString(), orgID, cfg.ProjectName)
// 		if err != nil {
// 			return fmt.Errorf("create default project: %w", err)
// 		}
// 	}

// 	return tx.Commit()
// }

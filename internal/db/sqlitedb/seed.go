package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/squeakycheese75/open-incentives/configs"
	"github.com/squeakycheese75/open-incentives/internal/services"
	"golang.org/x/crypto/bcrypt"
)

func Seed(ctx context.Context, db *sql.DB, cfg configs.BootstrapConfig) error {
	publicIDGenerator := services.NanoIDGenerator{}

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
	`, cfg.OrgPublicID).Scan(&orgID)

	if err == sql.ErrNoRows {
		res, err := tx.ExecContext(ctx, `
			INSERT INTO organizations (public_id, name)
			VALUES (?, ?)
		`, cfg.OrgPublicID, cfg.OrgName)
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
		projectPublicID, err := publicIDGenerator.New("proj")
		if err != nil {
			return fmt.Errorf("generate user public id: %w", err)
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO projects (public_id, org_id, name)
			VALUES (?, ?, ?)
		`, projectPublicID, orgID, cfg.ProjectName)
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

		userPublicID, err := publicIDGenerator.New("user")
		if err != nil {
			return fmt.Errorf("generate user public id: %w", err)
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
		`, userPublicID, orgID, cfg.AdminEmail, string(passwordHash), "admin")
		if err != nil {
			return fmt.Errorf("create admin user: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

package store

import "context"

func (s *store) WithTx(ctx context.Context, fn func(Store) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	txStore := &store{
		db:      s.db,
		queries: s.queries.WithTx(tx),
	}

	if err := fn(txStore); err != nil {
		return err
	}

	return tx.Commit()
}

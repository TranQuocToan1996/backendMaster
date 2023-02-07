package db

import (
	"context"
	"database/sql"
	"fmt"
)

var (
	txKey = struct{}{}
)

type Store interface {
	Querier

	TransferTx(ctx context.Context,
		arg TransferTxParams) (TransferResult, error)

	CreateUserTx(ctx context.Context,
		arg CreateUserTxParams) (CreateUserTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SQLStore) execTx(ctx context.Context,
	fn func(*Queries) error) error {

	var (
		// If nil, Isolation level of the driver or database's default level is used
		defaultOtps *sql.TxOptions = nil
	)

	tx, err := s.db.BeginTx(ctx, defaultOtps)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("tx rollback err: %v and %v", rollbackErr, err)
		}
		return err
	}

	return tx.Commit()
}

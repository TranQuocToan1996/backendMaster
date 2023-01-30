package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

var (
	txKey = struct{}{}
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (s *Store) execTx(ctx context.Context,
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

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferResult, error) {
	var result TransferResult

	err := s.execTx(ctx, func(q *Queries) error {
		txName := ctx.Value(txKey)

		log.Println(txName, "create transfer")
		transfer, err := q.CreateTransfer(ctx, CreateTransferParams(arg))
		result.Transfer = transfer
		if err != nil {
			return err
		}

		log.Println(txName, "create from entry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // Send money
		})
		if err != nil {
			return err
		}

		log.Println(txName, "create to entry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // Receive money
		})
		if err != nil {
			return err
		}

		log.Println(txName, "add to money")
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q,
				arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q,
				arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, id1, amount1 int64,
	id2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     id1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     id2,
		Amount: amount2,
	})

	return
}

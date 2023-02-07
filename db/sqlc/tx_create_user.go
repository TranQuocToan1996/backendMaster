package db

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams

	// callback runs after user is inserted
	// IE: Rollback, commit transaction
	AfterCreateTasks func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (s *SQLStore) CreateUserTx(ctx context.Context,
	arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		user, err := q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		result.User = user

		return arg.AfterCreateTasks(result.User)
	})

	return result, err
}

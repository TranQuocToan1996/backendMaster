package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	const (
		n      = 5
		amount = 10
	)

	store := NewStore(testDB)

	sendAcc := createRandomAccount(t)
	receiverAcc := createRandomAccount(t)

	var (
		ctx     = context.Background()
		errs    = make(chan error)
		results = make(chan TransferResult)
		existed = make(map[int]bool)
	)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			res, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: sendAcc.ID,
				ToAccountID:   receiverAcc.ID,
				Amount:        amount,
			})
			errs <- err
			results <- res
		}()
	}

	for err := range errs {
		require.NoError(t, err)

		result := <-results

		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, result.Transfer.FromAccountID, sendAcc.ID)
		require.Equal(t, result.Transfer.ToAccountID, receiverAcc.ID)
		require.Equal(t, result.Transfer.Amount, amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, errGet := store.GetTransfer(ctx, result.Transfer.ID)
		require.NoError(t, errGet)

		require.NotEmpty(t, result.FromEntry)
		require.NotEmpty(t, result.FromEntry.ID)
		require.Equal(t, result.FromEntry.AccountID, sendAcc.ID)
		require.Equal(t, result.FromEntry.Amount, -amount)

		_, errGet = store.GetEntry(ctx, result.FromEntry.ID)
		require.NoError(t, errGet)

		require.NotEmpty(t, result.ToEntry)
		require.NotEmpty(t, result.ToEntry.ID)
		require.Equal(t, result.ToEntry.AccountID, receiverAcc.ID)
		require.Equal(t, result.ToEntry.Amount, amount)

		_, errGet = store.GetEntry(ctx, result.ToEntry.ID)
		require.NoError(t, errGet)

		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, result.FromAccount.ID, sendAcc.ID)

		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, result.ToAccount.ID, sendAcc.ID)

		// check balances
		fmt.Println(">> tx:", result.FromAccount.Balance, receiverAcc.Balance)

		diff1 := sendAcc.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - receiverAcc.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), sendAcc.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), receiverAcc.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, sendAcc.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, receiverAcc.Balance+int64(n)*amount, updatedAccount2.Balance)

}

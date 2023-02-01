package db

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/stretchr/testify/require"
)

//TODO: fuzz testing instead of random
//TODO: write tests for entry and transfer

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	log.Println("test create account with: ", account)

	return account
}

func TestGetAccount(t *testing.T) {
	setAcc := createRandomAccount(t)
	t.Parallel()

	account, err := testQueries.GetAccount(context.Background(), setAcc.ID)
	require.NoError(t, err)

	require.True(t, account.CreatedAt.Equal(setAcc.CreatedAt))
	require.Equal(t, account.ID, setAcc.ID)
	require.Equal(t, account.Owner, setAcc.Owner)
	require.Equal(t, account.Balance, setAcc.Balance)
	require.Equal(t, account.Currency, setAcc.Currency)

}

func TestUpdateAccount(t *testing.T) {
	setAcc := createRandomAccount(t)
	t.Parallel()
	updateMoney := util.RandomMoney()

	arg := UpdateAccountParams{
		ID:      setAcc.ID,
		Balance: updateMoney,
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	require.True(t, account.CreatedAt.Equal(setAcc.CreatedAt))
	require.Equal(t, account.ID, setAcc.ID)
	require.Equal(t, account.Owner, setAcc.Owner)
	require.Equal(t, account.Balance, updateMoney)
	require.Equal(t, account.Currency, setAcc.Currency)

}

func TestDeleteAccount(t *testing.T) {
	setAcc := createRandomAccount(t)
	t.Parallel()

	err := testQueries.DeleteAccount(context.Background(), setAcc.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), setAcc.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _, account := range accounts {
		require.NotEmpty(t, account.ID)
		require.NotEmpty(t, account.Balance)
		require.NotEmpty(t, account.CreatedAt)
		require.NotEmpty(t, account.Currency)
		require.NotEmpty(t, account.Owner)
	}

}

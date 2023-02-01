package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/stretchr/testify/require"
)

//TODO: fuzz testing instead of random
//TODO: write tests for entry and transfer

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestGetAccount(t *testing.T) {
	setAcc := createRandomAccount(t)

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

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, result.Transfer.FromAccountID, account1.ID)
		require.Equal(t, result.Transfer.ToAccountID, account2.ID)
		require.Equal(t, result.Transfer.Amount, amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		//check account entries
		require.NotEmpty(t, result.FromEntry)
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, result.FromEntry.AccountID, account1.ID)
		require.Equal(t, result.ToEntry.AccountID, account2.ID)
		require.Equal(t, result.FromEntry.Amount, -amount)
		require.Equal(t, result.ToEntry.Amount, amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)
		require.NotZero(t, result.ToEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)
		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		//TODO: check account ballance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		diff1 := account1.Ballance - fromAccount.Ballance
		diff2 := toAccount.Ballance - account2.Ballance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
		// require.Equal(t, result.FromAccount.Ballance, account1.Ballance-amount)
		// require.Equal(t, result.ToAccount.Ballance, account2.Ballance+amount)
		// require.NotZero(t, result.FromEntry.ID)
	}
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, updateAccount1.Ballance, account1.Ballance-int64(n)*amount)
	require.Equal(t, updateAccount2.Ballance, account2.Ballance+int64(n)*amount)

}

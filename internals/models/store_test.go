package models

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTxHappyCase(t *testing.T) {
	store := NewStore(dbPoolTest)

	_, insertedAccount1, _ := createUser()
	_, insertedAccount2, _ := createUser()

	n := 10
	amount := int64(10)
	errorCh := make(chan error)
	resultTxCh := make(chan TransferTxResult)

	fmt.Println(">>Start transfer")
	fmt.Println(">>Number of transfer:", n)
	fmt.Println(">>Amount:", amount)
	fmt.Println(">>FromAccountID and balance:", insertedAccount1.ID, insertedAccount1.Balance)
	fmt.Println(">>ToAccountID and balance:", insertedAccount2.ID, insertedAccount2.Balance)
	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: insertedAccount1.ID,
				ToAccountID:   insertedAccount2.ID,
				Amount:        amount,
			})
			resultTxCh <- res
			errorCh <- err
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		result := <-resultTxCh
		err := <-errorCh
		assert.NoError(t, err, "TransferTx should not return error")

		fromEntry, getFromEntryErr := store.GetByIDEntry(context.Background(), result.FromEntry.ID)
		toEntry, getToEntryErr := store.GetByIDEntry(context.Background(), result.ToEntry.ID)

		assert.NoError(t, getFromEntryErr, "GetByIDEntry should not return error")
		assert.NoError(t, getToEntryErr, "GetByIDEntry should not return error")
		assert.NotEmpty(t, fromEntry, "FromEntry should not be empty")
		assert.Equal(t, insertedAccount1.ID, fromEntry.AccountID, "FromEntry.AccountID should be equal to insertedAccount1.ID")
		assert.Equal(t, -amount, fromEntry.Amount, "FromEntry.Amount should be equal to -amount")
		assert.NotZero(t, fromEntry.ID, "FromEntry.ID should not be zero")
		assert.NotZero(t, fromEntry.CreatedAt, "FromEntry.CreatedAt should not be zero")

		assert.NotEmpty(t, toEntry, "ToEntry should not be empty")
		assert.Equal(t, insertedAccount2.ID, toEntry.AccountID, "ToEntry.AccountID should be equal to insertedAccount2.ID")
		assert.Equal(t, amount, toEntry.Amount, "ToEntry.Amount should be equal to amount")
		assert.NotZero(t, toEntry.ID, "ToEntry.ID should not be zero")
		assert.NotZero(t, toEntry.CreatedAt, "ToEntry.CreatedAt should not be zero")

		// test balance account1
		resultFromAccount := result.FromAccount

		assert.NotEmpty(t, resultFromAccount, "FromAccount should not be empty")
		assert.Equal(t, insertedAccount1.ID, resultFromAccount.ID, "FromAccount.ID should be equal to insertedAccount1.ID")

		resultToAccount := result.ToAccount
		assert.NotEmpty(t, resultToAccount, "ToAccount should not be empty")
		assert.Equal(t, insertedAccount2.ID, resultToAccount.ID, "ToAccount.ID should be equal to insertedAccount2.ID")

		fmt.Println(">>Transaction:", resultFromAccount.Balance, resultToAccount.Balance)
		diff1 := insertedAccount1.Balance - resultFromAccount.Balance
		diff2 := resultToAccount.Balance - insertedAccount2.Balance
		assert.True(t, diff1 > 0, "Balance should be greater than 0")
		assert.Equal(t, diff1, diff2, "Balance should be equal")
		assert.True(t, diff1%amount == 0, "Balance should be multiple of amount")

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n, "k should be between 1 and n")
		assert.NotContains(t, existed, k, "k should not exist")
		existed[k] = true
	}

	updatedAccount1, err := store.GetByID(context.Background(), insertedAccount1.ID)
	assert.NoError(t, err, "GetByID should not return error")

	updatedAccount2, err := store.GetByID(context.Background(), insertedAccount2.ID)
	assert.NoError(t, err, "GetByID should not return error")

	assert.Equal(t, insertedAccount1.Balance-int64(n)*amount, updatedAccount1.Balance, "Balance should be equal")
	assert.Equal(t, insertedAccount2.Balance+int64(n)*amount, updatedAccount2.Balance, "Balance should be equal")

}

// the idea is
// account 1 transfer to account 2
// account 2 transfer to account 1
func TestTransferTxDeadlockCase(t *testing.T) {
	store := NewStore(dbPoolTest)

	_, insertedAccount1, _ := createUser()
	_, insertedAccount2, _ := createUser()

	n := 10
	amount := int64(10)
	errorCh := make(chan error)

	fmt.Println(">>Start transfer")
	fmt.Println(">>Number of transfer:", n)
	fmt.Println(">>Amount:", amount)
	account1ID := insertedAccount1.ID
	account2ID := insertedAccount2.ID
	fmt.Println(">>FromAccountID and balance:", account1ID, insertedAccount1.Balance)
	fmt.Println(">>ToAccountID and balance:", account2ID, insertedAccount2.Balance)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			go func() {
				_, err := store.TransferTx(context.Background(), TransferTxParams{
					FromAccountID: account2ID,
					ToAccountID:   account1ID,
					Amount:        amount,
				})
				errorCh <- err
			}()
		} else {
			go func() {
				_, err := store.TransferTx(context.Background(), TransferTxParams{
					FromAccountID: account1ID,
					ToAccountID:   account2ID,
					Amount:        amount,
				})
				errorCh <- err
			}()
		}

	}

	for i := 0; i < n; i++ {
		err := <-errorCh
		assert.NoError(t, err, "TransferTx should not return error")
	}

	updatedAccount1, err := store.GetByID(context.Background(), insertedAccount1.ID)
	assert.NoError(t, err, "GetByID should not return error")

	updatedAccount2, err := store.GetByID(context.Background(), insertedAccount2.ID)
	assert.NoError(t, err, "GetByID should not return error")

	assert.Equal(t, insertedAccount1.Balance, updatedAccount1.Balance, "Balance should be equal")
	assert.Equal(t, insertedAccount2.Balance, updatedAccount2.Balance, "Balance should be equal")

}

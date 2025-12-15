package models

import (
	"context"
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
	}

}

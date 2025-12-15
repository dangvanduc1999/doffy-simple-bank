package models

import (
	"context"
	"testing"

	"github.com/doffy/simple-bank/internals/utils"
	"github.com/stretchr/testify/assert"
)

func clear(accountID int64) {
	AccountQueries.DeleteByID(context.Background(), accountID)
}

func createUser() (CreateParams, Account, error) {
	var user CreateParams
	user = CreateParams{
		Owner:    utils.RandomString(6),
		Balance:  utils.RandomInt(1, 1000),
		Currency: utils.RandomCurrency(),
	}
	insertedAccount, err := AccountQueries.Create(context.Background(), user)
	return user, insertedAccount, err
}

func TestCreateAccount(t *testing.T) {
	user := CreateParams{
		Owner:    utils.RandomString(6),
		Balance:  utils.RandomInt(1, 1000),
		Currency: utils.RandomCurrency(),
	}
	insertedAccount, err := AccountQueries.Create(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, user.Owner, insertedAccount.Owner)
	assert.Equal(t, user.Balance, insertedAccount.Balance)
	assert.Equal(t, user.Currency, insertedAccount.Currency)
	assert.NotZero(t, insertedAccount.ID)
	assert.NotZero(t, insertedAccount.CreatedAt)

	clear(insertedAccount.ID)
}

func TestGetAccountByID(t *testing.T) {
	user, insertedAccount, err := createUser()

	assert.NoError(t, err, "CreateUser should not return error")
	account, err := AccountQueries.GetByID(context.Background(), insertedAccount.ID)

	assert.NoError(t, err, "GetByID should not return error")
	assert.Equal(t, user.Owner, account.Owner)
	assert.Equal(t, user.Balance, account.Balance)
	assert.Equal(t, user.Currency, account.Currency)
	assert.Equal(t, insertedAccount.ID, account.ID)
	assert.Equal(t, insertedAccount.CreatedAt, account.CreatedAt)

	clear(insertedAccount.ID)
}

func TestDeleteAccountByID(t *testing.T) {
	_, insertedAccount, err := createUser()

	assert.NoError(t, err, "CreateUser should not return error")
	err = AccountQueries.DeleteByID(context.Background(), insertedAccount.ID)

	assert.NoError(t, err, "DeleteByID should not return error")
	_, err = AccountQueries.GetByID(context.Background(), insertedAccount.ID)
	assert.Error(t, err, "GetByID should return error")

	clear(insertedAccount.ID)
}

func TestUpdateAccountOwner(t *testing.T) {
	_, insertedAccount, err := createUser()

	assert.NoError(t, err, "CreateUser should not return error")
	owner := utils.RandomString(6)
	err = AccountQueries.UpdateOwner(context.Background(), UpdateOwnerParams{
		ID:    insertedAccount.ID,
		Owner: owner,
	})

	assert.NoError(t, err, "UpdateOwner should not return error")
	account, err := AccountQueries.GetByID(context.Background(), insertedAccount.ID)
	assert.NoError(t, err, "GetByID should not return error")
	assert.Equal(t, insertedAccount.ID, account.ID)
	assert.NotEqual(t, insertedAccount.Owner, account.Owner)
	assert.Equal(t, owner, account.Owner)
	assert.Equal(t, insertedAccount.Balance, account.Balance)
	assert.Equal(t, insertedAccount.Currency, account.Currency)
	assert.Equal(t, insertedAccount.CreatedAt, account.CreatedAt)

	clear(insertedAccount.ID)
}

func TestGetAccountList(t *testing.T) {
	_, insertedAccount1, _ := createUser()
	_, insertedAccount2, _ := createUser()

	list := []Account{
		insertedAccount1,
		insertedAccount2,
	}

	accountList, _ := AccountQueries.GetList(context.Background())

	for i, account := range accountList {
		assert.Equal(t, list[i].Balance, account.Balance)
		assert.Equal(t, list[i].Currency, account.Currency)
		assert.Equal(t, list[i].ID, account.ID)
		assert.Equal(t, list[i].Owner, account.Owner)
		assert.Equal(t, list[i].CreatedAt, account.CreatedAt)
	}

	clear(insertedAccount1.ID)
	clear(insertedAccount2.ID)
}

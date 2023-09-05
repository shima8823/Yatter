package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"yatter-backend-go/app/domain/object"
)

func TestAccountCreate(t *testing.T) {
	ctx := context.Background()
	cleanupDB()

	t.Run("Success", func(t *testing.T) {
		err := accountRepo.Create(ctx, &object.Account{Username: "testuser", PasswordHash: "hashedpassword"})
		assert.NoError(t, err)

		createdAccount, err := accountRepo.Retrieve(ctx, "testuser")
		assert.NoError(t, err)
		assert.NotNil(t, createdAccount)
		assert.Equal(t, "testuser", createdAccount.Username)
	})

	t.Run("duplicate username", func(t *testing.T) {
		err := accountRepo.Create(ctx, &object.Account{Username: "testuser", PasswordHash: "hashedpassword"})
		assert.Error(t, err)
	})
}

func TestAccountRetrieve(t *testing.T) {
	ctx := context.Background()
	cleanupDB()

	err := accountRepo.Create(ctx, &object.Account{Username: "testuser", PasswordHash: "hashedpassword"})
	assert.NoError(t, err)

	t.Run("Found", func(t *testing.T) {
		account, err := accountRepo.Retrieve(ctx, "testuser")
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, "testuser", account.Username)
	})

	t.Run("Not found", func(t *testing.T) {
		account, err := accountRepo.Retrieve(ctx, "nonexistentuser")
		assert.Error(t, err)
		assert.Nil(t, account)
	})
}

package dao_test

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/stretchr/testify/assert"
)

func setupAccountDAO(t *testing.T) (repository.Account, func()) {
	dao, err := setupDAO()
	if err != nil {
		t.Fatal(err)
	}
	dao.InitAll()
	accountRepo := dao.Account()

	cleanup := func() {
		dao.InitAll()
		dao.Close()
	}

	return accountRepo, cleanup
}

func TestAccountCreate(t *testing.T) {
	repo, cleanup := setupAccountDAO(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		err := repo.Create(ctx, &object.Account{Username: "testuser", PasswordHash: "hashedpassword"})
		assert.NoError(t, err)

		createdAccount, err := repo.Retrieve(ctx, "testuser")
		assert.NoError(t, err)
		assert.NotNil(t, createdAccount)
		assert.Equal(t, "testuser", createdAccount.Username)
	})

	t.Run("duplicate username", func(t *testing.T) {
		err := repo.Create(ctx, &object.Account{Username: "testuser", PasswordHash: "hashedpassword"})
		assert.Error(t, err)
	})
}

func TestAccountRetrieve(t *testing.T) {
	repo, cleanup := setupAccountDAO(t)
	defer cleanup()
	ctx := context.Background()

	err := repo.Create(ctx, &object.Account{Username: "testuser", PasswordHash: "hashedpassword"})
	assert.NoError(t, err)

	t.Run("Found", func(t *testing.T) {
		account, err := repo.Retrieve(ctx, "testuser")
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, "testuser", account.Username)
	})

	t.Run("Not found", func(t *testing.T) {
		account, err := repo.Retrieve(ctx, "nonexistentuser")
		assert.Error(t, err)
		assert.Nil(t, account)
	})
}

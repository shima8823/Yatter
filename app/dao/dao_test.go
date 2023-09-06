package dao_test

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/stretchr/testify/assert"
)

var accountRepo repository.Account
var statusRepo repository.Status
var relationshipRepo repository.Relationship
var cleanupDB func()

func TestMain(m *testing.M) {
	if dao, err := setupDAO(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		cleanupDB = func() {
			dao.InitAll()
		}
		defer dao.Close()

		accountRepo = dao.Account()
		statusRepo = dao.Status()
		relationshipRepo = dao.Relationship()
	}

	os.Exit(m.Run())
}

func insertAccountDB(t *testing.T, ctx context.Context, accounts []object.Account) {
	for _, account := range accounts {
		err := accountRepo.Create(ctx, &account)
		assert.NoError(t, err)
	}
}

func insertRelationshipDB(t *testing.T, ctx context.Context, relationships []object.Relationship) {
	for _, relationship := range relationships {
		err := relationshipRepo.Create(ctx, relationship.FollowingId, relationship.FollowerId)
		assert.NoError(t, err)
	}
}

func createAccountObject(num int) []object.Account {
	accounts := make([]object.Account, num)
	for i := 0; i < num; i++ {
		accounts[i] = object.Account{
			Username:     "test" + strconv.Itoa(i),
			PasswordHash: "password",
		}
	}
	return accounts
}

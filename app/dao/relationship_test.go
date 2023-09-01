package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
)

var accountRepo repository.Account
var relationshipRepo repository.Relationship
var cleanupDB func()

func TestMain(m *testing.M) {
	if dao, err := setupDAO(); err != nil {
		os.Exit(1)
	} else {
		cleanupDB = func() {
			dao.InitAll()
		}
		defer dao.Close()

		accountRepo = dao.Account()
		relationshipRepo = dao.Relationship()
	}

	os.Exit(m.Run())
}

func insertAccountDB(t *testing.T, ctx context.Context, accounts []object.Account) {
	for _, account := range accounts {
		err := accountRepo.CreateUser(ctx, &account)
		assert.NoError(t, err)
	}
}

func insertRelationshipDB(t *testing.T, ctx context.Context, relationships []object.Relationship) {
	for _, relationship := range relationships {
		err := relationshipRepo.CreateFollowing(ctx, relationship.FollowingId, relationship.FollowerId)
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

func TestCreateFollowing(t *testing.T) {
	cleanupDB()
	ctx := context.Background()
	insertAccountDB(t, ctx, createAccountObject(2))
	relationship := &object.Relationship{
		FollowingId: 1,
		FollowerId:  2,
	}
	err := relationshipRepo.CreateFollowing(ctx, relationship.FollowingId, relationship.FollowerId)
	assert.NoError(t, err)
}

func TestDeleteFollowing(t *testing.T) {
	cleanupDB()
	ctx := context.Background()
	insertAccountDB(t, ctx, createAccountObject(2))
	relationship := &object.Relationship{
		FollowingId: 1,
		FollowerId:  2,
	}
	insertRelationshipDB(t, ctx, []object.Relationship{
		*relationship,
	})

	err := relationshipRepo.DeleteFollowing(ctx, relationship.FollowingId, relationship.FollowerId)
	assert.NoError(t, err)
}

func TestFeatchRelationships(t *testing.T) {
	cleanupDB()
	ctx := context.Background()
	insertAccountDB(t, ctx, createAccountObject(2))
	relationships := []object.Relationship{
		{
			FollowingId: 1,
			FollowerId:  2,
		},
		{
			FollowingId: 2,
			FollowerId:  1,
		},
	}
	insertRelationshipDB(t, ctx, relationships)

	relationships, err := relationshipRepo.FeatchRelationships(ctx, relationships[0].FollowingId)
	assert.NoError(t, err)
	assert.NotNil(t, relationships)
	assert.Equal(t, 2, len(relationships))
}

func TestFeatchFollowing(t *testing.T) {
	cleanupDB()
	ctx := context.Background()
	insertAccountDB(t, ctx, createAccountObject(3))
	relationships := []object.Relationship{
		{
			FollowingId: 1,
			FollowerId:  2,
		},
		{
			FollowingId: 1,
			FollowerId:  3,
		},
	}
	limit := uint64(2)
	insertRelationshipDB(t, ctx, relationships)

	accounts, err := relationshipRepo.FeatchFollowing(ctx, relationships[0].FollowingId, &limit)
	assert.NoError(t, err)
	assert.NotNil(t, accounts)
	assert.Equal(t, 2, len(accounts))
}

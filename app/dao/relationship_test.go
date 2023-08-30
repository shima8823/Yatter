package dao_test

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/stretchr/testify/assert"
)

func setupRelationshipDAO(t *testing.T) (repository.Relationship, func()) {
	dao := setupDAO(t)
	dao.InitAll()
	relationshipRepo := dao.Relationship()

	cleanup := func() {
		dao.InitAll()
	}

	return relationshipRepo, cleanup
}

func setupAccountDB(t *testing.T, ctx context.Context) {
	dao := setupDAO(t)
	dao.InitAll()
	accountRepo := dao.Account()

	accounts := []object.Account{
		{
			Username:     "test1",
			PasswordHash: "password",
		},
		{
			Username:     "test2",
			PasswordHash: "password",
		},
	}

	for _, account := range accounts {
		err := accountRepo.CreateUser(ctx, &account)
		assert.NoError(t, err)
	}
}

func TestCreateFollowing(t *testing.T) {
	repo, cleanup := setupRelationshipDAO(t)
	defer cleanup()

	ctx := context.Background()
	setupAccountDB(t, ctx)

	relationship := &object.Relationship{
		FollowingId: 1,
		FollowerId:  2,
	}
	var err error
	err = repo.CreateFollowing(ctx, relationship.FollowingId, relationship.FollowerId)
	assert.NoError(t, err)
}

func TestDeleteFollowing(t *testing.T) {
	repo, cleanup := setupRelationshipDAO(t)
	defer cleanup()

	ctx := context.Background()
	setupAccountDB(t, ctx)

	relationship := &object.Relationship{
		FollowingId: 1,
		FollowerId:  2,
	}
	var err error
	err = repo.CreateFollowing(ctx, relationship.FollowingId, relationship.FollowerId)
	assert.NoError(t, err)

	err = repo.DeleteFollowing(ctx, relationship.FollowingId, relationship.FollowerId)
	assert.NoError(t, err)
}

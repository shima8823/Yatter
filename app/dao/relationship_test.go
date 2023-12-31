package dao_test

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"

	"github.com/stretchr/testify/assert"
)

func TestRelationshipCreate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		accounts      []object.Account
		relationships []object.Relationship
		wantErr       bool
	}{
		{
			name:     "success",
			accounts: createAccountObject(2),
			relationships: []object.Relationship{
				{
					FollowingId: 1,
					FollowerId:  2,
				},
			},
			wantErr: false,
		},
		{
			name:     "duplicate",
			accounts: createAccountObject(2),
			relationships: []object.Relationship{
				{
					FollowingId: 1,
					FollowerId:  2,
				},
				{
					FollowingId: 1,
					FollowerId:  2,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		cleanupDB()
		insertAccountDB(t, ctx, tt.accounts)
		for i, relationship := range tt.relationships {
			err := relationshipRepo.Create(ctx, relationship.FollowingId, relationship.FollowerId)
			if i == len(tt.relationships)-1 { // 重複エラーは最後のみ
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			} else {
				assert.NoError(t, err)
			}
		}

	}
}

func TestRelationshipDelete(t *testing.T) {
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

	err := relationshipRepo.Delete(ctx, relationship.FollowingId, relationship.FollowerId)
	assert.NoError(t, err)
}

func TestRelationshipRetrieve(t *testing.T) {
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

	relationships, err := relationshipRepo.Retrieve(ctx, relationships[0].FollowingId)
	assert.NoError(t, err)
	assert.NotNil(t, relationships)
	assert.Equal(t, 2, len(relationships))
}

func TestRetrieveFollowing(t *testing.T) {
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

	accounts, err := relationshipRepo.RetrieveFollowing(ctx, relationships[0].FollowingId, &limit)
	assert.NoError(t, err)
	assert.NotNil(t, accounts)
	assert.Equal(t, 2, len(accounts))
}

func TestRetrieveFollowers(t *testing.T) {
	cleanupDB()
	ctx := context.Background()
	insertAccountDB(t, ctx, createAccountObject(3))
	relationships := []object.Relationship{
		{
			FollowingId: 1,
			FollowerId:  2,
		},
		{
			FollowingId: 3,
			FollowerId:  2,
		},
	}
	max_id := uint64(3)
	since_id := uint64(1)
	limit := uint64(2)
	insertRelationshipDB(t, ctx, relationships)

	accounts, err := relationshipRepo.RetrieveFollowers(ctx, relationships[0].FollowerId, nil, &max_id, &since_id, &limit)
	assert.NoError(t, err)
	assert.NotNil(t, accounts)
	assert.Equal(t, 2, len(accounts))
}

func TestCountFollowing(t *testing.T) {
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
	insertRelationshipDB(t, ctx, relationships)

	count, err := relationshipRepo.CountFollowing(ctx, relationships[0].FollowingId)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
}

func TestCountFollowers(t *testing.T) {
	cleanupDB()
	ctx := context.Background()
	insertAccountDB(t, ctx, createAccountObject(3))
	relationships := []object.Relationship{
		{
			FollowingId: 1,
			FollowerId:  2,
		},
		{
			FollowingId: 3,
			FollowerId:  2,
		},
	}
	insertRelationshipDB(t, ctx, relationships)

	count, err := relationshipRepo.CountFollowers(ctx, relationships[0].FollowerId)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), count)
}

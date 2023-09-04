package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Relationship interface {
	Create(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error
	Delete(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error
	Retrieve(ctx context.Context, accountID object.AccountID) ([]object.Relationship, error)
	RetrieveFollowing(ctx context.Context, accountID object.AccountID, limit *uint64) ([]object.Account, error)
	RetrieveFollowers(ctx context.Context, accountID object.AccountID, only_media, max_id, since_id, limit *uint64) ([]object.Account, error)
	CountFollowing(ctx context.Context, accountID object.AccountID) (uint64, error)
	CountFollowers(ctx context.Context, accountID object.AccountID) (uint64, error)
}

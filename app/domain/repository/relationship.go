package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Relationship interface {
	// アカウントのfollowを作成
	CreateFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error
	// アカウントのunfollowを作成
	DeleteFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error

	// アカウントとの関係を取得
	FeatchRelationships(ctx context.Context, accountID object.AccountID) ([]object.Relationship, error)

	// following一覧を取得
	FeatchFollowing(ctx context.Context, accountID object.AccountID, limit *uint64) ([]object.Account, error)
	// follower一覧を取得
	FeatchFollowers(ctx context.Context, accountID object.AccountID, only_media, max_id, since_id, limit *uint64) ([]object.Account, error)

	// アカウントのfollowing数を取得
	CountFollowing(ctx context.Context, accountID object.AccountID) (uint64, error)
	// アカウントのfollower数を取得
	CountFollower(ctx context.Context, accountID object.AccountID) (uint64, error)

	// utils
	FindAccountByUsername(ctx context.Context, username string) (*object.Account, error)
}

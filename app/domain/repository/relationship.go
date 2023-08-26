package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Relationship interface {
	// アカウントのfollowを作成
	CreateFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error
	// // アカウントのunfollowを作成
	// DeleteFollowing(ctx context.Context, accountID object.AccountID, followingID object.AccountID) error

	// // アカウントとの関係を取得
	// FindRelationship(ctx context.Context, accountID object.AccountID, followingID object.AccountID) (*object.Relationship, error)

	// // following一覧を取得
	// FindFollowing(ctx context.Context, accountID object.AccountID) ([]*object.Account, error)
	// // follower一覧を取得
	// FindFollower(ctx context.Context, accountID object.AccountID) ([]*object.Account, error)

	// // followしてるtimelineを取得
	// FindFollowingTimeline(ctx context.Context, accountID object.AccountID) ([]*object.Status, error)

	// // アカウントのfollowing数を取得
	// CountFollowing(ctx context.Context, accountID object.AccountID) (int, error)
	// // アカウントのfollower数を取得
	// CountFollower(ctx context.Context, accountID object.AccountID) (int, error)
}

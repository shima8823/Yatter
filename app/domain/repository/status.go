package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {
	Create(ctx context.Context, status *object.Status) error
	Retrieve(ctx context.Context, id uint64) (*object.Status, error)
	Delete(ctx context.Context, id uint64) error

	PublicTimeline(ctx context.Context, only_media, max_id, since_id, limit *uint64) ([]*object.Status, error)
	HomeTimeline(ctx context.Context, accountID object.AccountID, only_media, max_id, since_id, limit *uint64) ([]object.Status, error)
}

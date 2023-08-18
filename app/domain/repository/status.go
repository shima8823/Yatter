package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Create status
	CreateStatus(ctx context.Context, status *object.Status) error

	// Find status by id
	FindByID(ctx context.Context, id uint64) (*object.Status, error)
}

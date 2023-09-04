package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Account interface {
	Retrieve(ctx context.Context, username string) (*object.Account, error)
	Create(ctx context.Context, account *object.Account) error
}

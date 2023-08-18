package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Status
	status struct {
		db *sqlx.DB
	}
)

// Create status repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

func (r *status) CreateStatus(ctx context.Context, status *object.Status) error {
	_, err := r.db.ExecContext(ctx, "insert into status (account_id, content) values (?, ?)", status.AccountId, status.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *status) FindByID(ctx context.Context, id uint64) (*object.Status, error) {
	entity := new(object.Status)
	println(id)
	err := r.db.QueryRowxContext(ctx, "select * from status where id = ?", id).StructScan(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

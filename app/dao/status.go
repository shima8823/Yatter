package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	status struct {
		db *sqlx.DB
	}
)

func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

func (r *status) Create(ctx context.Context, status *object.Status) error {
	_, err := r.db.ExecContext(ctx, "insert into status (account_id, content) values (?, ?)", status.AccountId, status.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *status) Retrieve(ctx context.Context, id uint64) (*object.Status, error) {
	entity := new(object.Status)
	err := r.db.QueryRowxContext(ctx, "select * from status where id = ?", id).StructScan(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *status) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, "delete from status where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *status) PublicTimeline(ctx context.Context, only_media, max_id, since_id, limit *uint64) ([]*object.Status, error) {
	var entities []*object.Status

	query, args := buildQuery("status", "id", since_id, max_id, limit)
	err := r.db.SelectContext(ctx, &entities, query, args...)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *status) HomeTimeline(ctx context.Context, accountID object.AccountID, only_media, max_id, since_id, limit *uint64) ([]object.Status, error) {
	var entities []object.Status

	query := `select status.* from status join relationship on status.account_id = relationship.follower_id where relationship.following_id = ?`

	args := []interface{}{accountID}

	// TODO only_media

	if max_id != nil {
		query += " and status.id <= ?"
		args = append(args, *max_id)
	}

	if since_id != nil {
		query += " and status.id >= ?"
		args = append(args, *since_id)
	}

	query += " order by status.create_at desc"
	if limit != nil {
		query += " limit ?"
		args = append(args, *limit)
	}
	err := r.db.SelectContext(ctx, &entities, query, args...)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

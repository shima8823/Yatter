package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Status
	relationship struct {
		db *sqlx.DB
	}
)

// Create status repository
func NewRelationship(db *sqlx.DB) repository.Relationship {
	return &relationship{db: db}
}

func (r *relationship) CreateFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error {
	if _, err := r.db.ExecContext(ctx, "insert into relationship (following_id, follower_id) values (?, ?)", followingID, followerID); err != nil {
		return err
	}
	return nil
}

func (r *relationship) FindAccountByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

func (r *relationship) DeleteFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error {
	if _, err := r.db.ExecContext(ctx, "delete from relationship where following_id = ? and follower_id = ?", followingID, followerID); err != nil {
		return err
	}
	return nil
}

func (r *relationship) FeatchRelationships(ctx context.Context, accountID object.AccountID) ([]object.Relationship, error) {
	var entities []object.Relationship

	rows, err := r.db.QueryxContext(ctx, "select * from relationship where following_id = ? or follower_id = ?", accountID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity object.Relationship
		if err := rows.StructScan(&entity); err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *relationship) FeatchFollowing(ctx context.Context, accountID object.AccountID, limit *uint64) ([]object.Account, error) {
	var entities []object.Account

	query := `
		SELECT account.*
		FROM account
		JOIN relationship ON account.id = relationship.follower_id
		WHERE relationship.following_id = ?`

	args := []interface{}{accountID}

	if limit != nil {
		query += " LIMIT ?"
		args = append(args, *limit)
	}

	err := r.db.SelectContext(ctx, &entities, query, args...)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *relationship) FeatchFollowers(ctx context.Context, accountID object.AccountID, only_media, max_id, since_id, limit *uint64) ([]object.Account, error) {
	var entities []object.Account

	query := `
		SELECT account.*
		FROM account
		JOIN relationship ON account.id = relationship.following_id
		WHERE relationship.follower_id = ?`

	args := []interface{}{accountID}

	// TODO only_media

	if max_id != nil {
		query += " AND account.id <= ?"
		args = append(args, *max_id)
	}

	if since_id != nil {
		query += " AND account.id >= ?"
		args = append(args, *since_id)
	}

	query += " ORDER BY relationship.create_at DESC"
	if limit != nil {
		query += " LIMIT ?"
		args = append(args, *limit)
	}

	err := r.db.SelectContext(ctx, &entities, query, args...)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *relationship) CountFollowing(ctx context.Context, accountID object.AccountID) (uint64, error) {
	var count uint64
	err := r.db.QueryRowxContext(ctx, "select count(*) from relationship where following_id = ?", accountID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *relationship) CountFollower(ctx context.Context, accountID object.AccountID) (uint64, error) {
	var count uint64
	err := r.db.QueryRowxContext(ctx, "select count(*) from relationship where follower_id = ?", accountID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

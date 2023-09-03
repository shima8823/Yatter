package dao

import (
	"context"
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

// following フォローする人
// follower フォローされる人

func (r *relationship) CreateFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, "insert into relationship (following_id, follower_id) values (?, ?)", followingID, followerID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := r.db.ExecContext(ctx, "update account set following_count = following_count + 1 where id = ?", followingID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := r.db.ExecContext(ctx, "update account set followers_count = followers_count + 1 where id = ?", followerID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *relationship) FindAccountByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *relationship) DeleteFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if res, err := r.db.ExecContext(ctx, "delete from relationship where following_id = ? and follower_id = ?", followingID, followerID); err != nil {
		tx.Rollback()
		return err
	} else {
		if res, _ := res.RowsAffected(); res == 0 {
			return fmt.Errorf("not found")
		}
	}
	if _, err := r.db.ExecContext(ctx, "update account set following_count = following_count - 1 where id = ?", followingID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := r.db.ExecContext(ctx, "update account set followers_count = followers_count - 1 where id = ?", followerID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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

	query := `select account.* from account join relationship on account.id = relationship.follower_id where relationship.following_id = ?`

	args := []interface{}{accountID}

	query += " order by relationship.create_at desc"

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

func (r *relationship) FeatchFollowers(ctx context.Context, accountID object.AccountID, only_media, max_id, since_id, limit *uint64) ([]object.Account, error) {
	var entities []object.Account

	query := `select account.* from account join relationship on account.id = relationship.following_id where relationship.follower_id = ?`

	args := []interface{}{accountID}

	// TODO only_media

	if max_id != nil {
		query += " and account.id <= ?"
		args = append(args, *max_id)
	}

	if since_id != nil {
		query += " and account.id >= ?"
		args = append(args, *since_id)
	}

	query += " order by relationship.create_at desc"
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

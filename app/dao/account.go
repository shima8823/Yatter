package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	account struct {
		db *sqlx.DB
	}
)

func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

func (r *account) Create(ctx context.Context, account *object.Account) error {
	_, err := r.db.ExecContext(ctx, "insert into account (username, password_hash) values (?, ?)", account.Username, account.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (r *account) Retrieve(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

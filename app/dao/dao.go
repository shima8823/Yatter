package dao

import (
	"database/sql"
	"fmt"
	"log"
	"yatter-backend-go/app/domain/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type (
	Dao interface {
		Account() repository.Account
		Status() repository.Status
		Relationship() repository.Relationship

		// Clear all data in DB
		// This function is "only" used for testing
		// Because it disables foreign key constraints after clearing all data
		InitAll() error

		// Close DB connection for testing
		Close() error
	}

	dao struct {
		db *sqlx.DB
	}
)

func New(config DBConfig) (Dao, error) {
	db, err := initDb(config)
	if err != nil {
		return nil, err
	}

	return &dao{db: db}, nil
}

func NewMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db, mock
}

func NewWithDB(db *sqlx.DB) Dao {
	return &dao{db: db}
}

func (d *dao) Account() repository.Account {
	return NewAccount(d.db)
}

func (d *dao) Status() repository.Status {
	return NewStatus(d.db)
}

func (d *dao) Relationship() repository.Relationship {
	return NewRelationship(d.db)
}

// 外部キー制約を無効化して全テーブルをクリアする
// 外部キー制約を無効化した場合、参照先のテーブルのデータを削除する必要がなくなる
func (d *dao) InitAll() error {
	if err := d.exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return fmt.Errorf("Can't disable FOREIGN_KEY_CHECKS: %w", err)
	}

	// テストのために無効化する。
	// ただし、DAOのインタフェースとして提供されている以上、外部キー制約を有効化するべき。
	// should: SET FOREIGN_KEY_CHECKS=0 -> SET FOREIGN_KEY_CHECKS=1
	defer func() {
		err := d.exec("SET FOREIGN_KEY_CHECKS=0")
		if err != nil {
			log.Printf("Can't restore FOREIGN_KEY_CHECKS: %+v", err)
		}
	}()

	for _, table := range []string{"account", "status", "relationship"} {
		if err := d.exec("TRUNCATE TABLE " + table); err != nil {
			return fmt.Errorf("Can't truncate table "+table+": %w", err)
		}
	}

	return nil
}

func (d *dao) exec(query string, args ...interface{}) error {
	_, err := d.db.Exec(query, args...)
	return err
}

func (d *dao) Close() error {
	return d.db.Close()
}

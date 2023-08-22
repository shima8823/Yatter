package dao_test

import (
	"context"
	"database/sql"
	"testing"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/object"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// account dao test
func TestFindByUsername(t *testing.T) {
	// データベースとモックを作成
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub database connection: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := dao.NewAccount(sqlxDB)

	// Mock setup
	expectedUser := &object.Account{Username: "testuser", PasswordHash: "hashedpassword"}
	mock.ExpectQuery("select \\* from account where username = ?").
		WithArgs(expectedUser.Username).
		WillReturnRows(sqlmock.NewRows([]string{"username", "password_hash"}).
			AddRow(expectedUser.Username, expectedUser.PasswordHash))

	ctx := context.Background()
	result, err := r.FindByUsername(ctx, expectedUser.Username)
	if err != nil {
		t.Fatal(err)
	}

	// mockのqueryが実行されたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	// 結果の確認
	if result.Username != expectedUser.Username {
		t.Fatalf("expected: %v, got: %v", expectedUser.Username, result.Username)
	}
}

func TestFindByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub database connection: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := dao.NewAccount(sqlxDB)

	mock.ExpectQuery("select \\* from account where username = ?").
		WithArgs("nonexistuser").
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	result, err := r.FindByUsername(ctx, "nonexistuser")
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected: %v, got: %v", nil, result)
	}
}

func TestFindByUsername_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub database connection: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := dao.NewAccount(sqlxDB)

	expectedUser := &object.Account{Username: "testuser", PasswordHash: "hashedpassword"}
	mock.ExpectQuery("select \\* from account where username = ?").
		WithArgs(expectedUser.Username).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	result, err := r.FindByUsername(ctx, expectedUser.Username)
	if err == nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected: %v, got: %v", nil, result)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub database connection: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := dao.NewAccount(sqlxDB)

	newUser := &object.Account{Username: "newuser", PasswordHash: "newhashedpassword"}
	mock.ExpectExec("insert into account \\(username, password_hash\\) values \\(\\?, \\?\\)").
		WithArgs(newUser.Username, newUser.PasswordHash).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	if err := r.CreateUser(ctx, newUser); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub database connection: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := dao.NewAccount(sqlxDB)

	newUser := &object.Account{Username: "newuser", PasswordHash: "newhashedpassword"}
	mock.ExpectExec("insert into account \\(username, password_hash\\) values \\(\\?, \\?\\)").
		WithArgs(newUser.Username, newUser.PasswordHash).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	if err := r.CreateUser(ctx, newUser); err == nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

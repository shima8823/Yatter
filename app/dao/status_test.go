package dao_test

import (
	"context"
	"database/sql"
	"testing"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type TestDBHelper struct {
	db   *sqlx.DB
	mock sqlmock.Sqlmock
	repo repository.Status
}

func setupTestDB(t *testing.T) *TestDBHelper {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub database connection: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := dao.NewStatus(sqlxDB)

	return &TestDBHelper{
		db:   sqlxDB,
		mock: mock,
		repo: r,
	}
}

func TestCreateStatus(t *testing.T) {
	helper := setupTestDB(t)
	defer helper.db.Close()

	newStatus := &object.Status{AccountId: 1, Content: "test content"}
	helper.mock.ExpectExec("insert into status \\(account_id, content\\) values \\(\\?, \\?\\)").
		WithArgs(newStatus.AccountId, newStatus.Content).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	if err := helper.repo.CreateStatus(ctx, newStatus); err != nil {
		t.Fatalf("failed to create status: %v", err)
	}
	if err := helper.mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateStatus_Error(t *testing.T) {
	helper := setupTestDB(t)

	newStatus := &object.Status{AccountId: 1, Content: "test content"}
	helper.mock.ExpectExec("insert into status \\(account_id, content\\) values \\(\\?, \\?\\)").
		WithArgs(newStatus.AccountId, newStatus.Content).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	if err := helper.repo.CreateStatus(ctx, newStatus); err == nil {
		t.Fatalf("expected error, but got nil")
	}
	if err := helper.mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindByID(t *testing.T) {
	helper := setupTestDB(t)

	expectedStatus := &object.Status{ID: 1, AccountId: 1, Content: "test content"}
	helper.mock.ExpectQuery("select \\* from status where id = ?").
		WithArgs(expectedStatus.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "content"}).
			AddRow(expectedStatus.ID, expectedStatus.AccountId, expectedStatus.Content))

	ctx := context.Background()
	result, err := helper.repo.FindByID(ctx, expectedStatus.ID)
	if err != nil {
		t.Fatalf("failed to find status: %v", err)
	}
	if result == nil {
		t.Fatalf("expected: %v, got: %v", expectedStatus, result)
	}
}

func TestFindByID_Error(t *testing.T) {
	helper := setupTestDB(t)

	expectedStatus := &object.Status{ID: 1, AccountId: 1, Content: "test content"}
	helper.mock.ExpectQuery("select \\* from status where id = ?").
		WithArgs(expectedStatus.ID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	result, err := helper.repo.FindByID(ctx, expectedStatus.ID)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
	if result != nil {
		t.Fatalf("expected: %v, got: %v", nil, result)
	}
}

func TestDeleteByID(t *testing.T) {
	helper := setupTestDB(t)

	expectedStatus := &object.Status{ID: 1, AccountId: 1, Content: "test content"}
	helper.mock.ExpectExec("delete from status where id = ?").
		WithArgs(expectedStatus.ID).
		WillReturnResult(sqlmock.NewResult(-1, 1))

	ctx := context.Background()
	if err := helper.repo.DeleteByID(ctx, expectedStatus.ID); err != nil {
		t.Fatalf("failed to delete status: %v", err)
	}
	if err := helper.mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteByID_Error(t *testing.T) {
	helper := setupTestDB(t)

	expectedStatus := &object.Status{ID: 1, AccountId: 1, Content: "test content"}
	helper.mock.ExpectExec("delete from status where id = ?").
		WithArgs(expectedStatus.ID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	if err := helper.repo.DeleteByID(ctx, expectedStatus.ID); err == nil {
		t.Fatalf("expected error, but got nil")
	}
	if err := helper.mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

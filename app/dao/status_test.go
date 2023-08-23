package dao_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
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

func TestPublicTimeline(t *testing.T) {
	helper := setupTestDB(t)
	ctx := context.Background()

	tests := []struct {
		name     string
		sinceID  *uint64
		maxID    *uint64
		limit    *uint64
		expected string
	}{
		{"Simple Timeline", nil, nil, nil, "SELECT \\* FROM status ORDER BY create_at DESC"},
		{"Since ID", pointerToUint64(5), nil, nil, "SELECT \\* FROM status WHERE id >= \\? ORDER BY create_at DESC"},
		{"Max ID", nil, pointerToUint64(10), nil, "SELECT \\* FROM status WHERE id <= \\? ORDER BY create_at DESC"},
		{"Limited", nil, nil, pointerToUint64(3), "SELECT \\* FROM status ORDER BY create_at DESC LIMIT \\?"},
		// Add more scenarios as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedArgs := extractArguments(tt.sinceID, tt.maxID, tt.limit)
			rows := sqlmock.NewRows([]string{"id", "account_id", "content"})

			args := make([]driver.Value, len(expectedArgs))
			for i, v := range expectedArgs {
				args[i] = v.(driver.Value)
			}
			helper.mock.ExpectQuery(tt.expected).WithArgs(args...).WillReturnRows(rows)

			res, err := helper.repo.PublicTimeline(ctx, nil, tt.maxID, tt.sinceID, tt.limit)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := helper.mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
			if res != nil {
				t.Fatalf("expected: %v, got: %v", nil, res)
			}
		})
	}
}

func TestPublicTimeline_Error(t *testing.T) {
	helper := setupTestDB(t)
	ctx := context.Background()

	tests := []struct {
		name     string
		sinceID  *uint64
		maxID    *uint64
		limit    *uint64
		expected string
	}{
		{"Simple Timeline", nil, nil, nil, "SELECT \\* FROM status ORDER BY create_at DESC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedArgs := extractArguments(tt.sinceID, tt.maxID, tt.limit)

			args := make([]driver.Value, len(expectedArgs))
			for i, v := range expectedArgs {
				args[i] = v.(driver.Value)
			}
			helper.mock.ExpectQuery(tt.expected).WithArgs(args...).WillReturnError(sql.ErrNoRows)

			if _, err := helper.repo.PublicTimeline(ctx, nil, tt.maxID, tt.sinceID, tt.limit); err == nil {
				t.Fatalf("expected error, but got nil")
			}
			if err := helper.mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func pointerToUint64(v uint64) *uint64 {
	return &v
}

func extractArguments(sinceID, maxID, limit *uint64) []interface{} {
	var args []interface{}
	if sinceID != nil {
		args = append(args, *sinceID)
	}
	if maxID != nil {
		args = append(args, *maxID)
	}
	if limit != nil {
		args = append(args, *limit)
	}
	return args
}

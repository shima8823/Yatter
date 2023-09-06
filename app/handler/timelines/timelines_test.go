package timelines

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetPublic(t *testing.T) {
	db, mock := dao.NewMockDB()
	h := newMockHandler(db)
	defer db.Close()

	tests := []struct {
		name     string
		username string
		mockFunc func()
		wantCode int
	}{
		{
			name:     "Success",
			username: "testuser",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from status order by create_at desc limit \\?").
					WithArgs(40).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "content"}).
						AddRow(1, 1, "test content").
						AddRow(2, 1, "test content2"))
			},
			wantCode: http.StatusOK,
		},
		{
			name: "no timeline",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from status order by create_at desc limit \\?").
					WithArgs(40).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "content"}))
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/v1/timelines", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.mockFunc != nil {
				tt.mockFunc()
			}
			h.GetPublic(w, r)
			assert.Equal(t, tt.wantCode, w.Code, "status code should be equal")
		})
	}
}

func newMockHandler(db *sql.DB) *handler {
	return &handler{
		app: &app.App{
			Dao: dao.NewWithDB(sqlx.NewDb(db, "sqlmock")),
		},
	}
}

func setChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

package relationships

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/handler/auth"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock := dao.NewMockDB()
	h := newMockHandler(db)
	defer db.Close()

	tests := []struct {
		name         string
		mockFunc     func()
		isAuth       bool
		urlParamFunc func(r *http.Request) *http.Request
		wantCode     int
	}{
		{
			name: "successfully create account",
			mockFunc: func() {
				// auth.Middlewareで実行されるクエリ
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
				// followする人のクエリ
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(2, "testuser2"))

				mock.ExpectBegin()

				mock.ExpectExec("insert into relationship \\(following_id, follower_id\\) values \\(\\?, \\?\\)").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("update account set following_count = following_count \\+ 1 where id = \\?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("update account set followers_count = followers_count \\+ 1 where id = \\?").
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			isAuth:       true,
			urlParamFunc: func(r *http.Request) *http.Request { return setChiURLParam(r, "username", "testuser") },
			wantCode:     http.StatusOK,
		},
		{
			name:     "Unauthorized",
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "same following and follower",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
			},
			isAuth:       true,
			urlParamFunc: func(r *http.Request) *http.Request { return setChiURLParam(r, "username", "testuser") },
			wantCode:     http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPost, "/v1/accounts/testuser/follow", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.urlParamFunc != nil {
				r = tt.urlParamFunc(r)
			}
			if tt.isAuth {
				r.Header.Set("Authentication", "username testuser")
			}
			if tt.mockFunc != nil {
				tt.mockFunc()
			}

			middleware := auth.Middleware(h.app)
			handlerMiddleware := middleware(http.HandlerFunc(h.Create))
			handlerMiddleware.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}

}

func TestDelete(t *testing.T) {
	db, mock := dao.NewMockDB()
	h := newMockHandler(db)
	defer db.Close()

	tests := []struct {
		name         string
		mockFunc     func()
		isAuth       bool
		urlParamFunc func(r *http.Request) *http.Request
		wantCode     int
	}{
		{
			name: "successfully delete account",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("followingUser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "followingUser"))
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("unfollowUser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(2, "unfollowUser"))

				mock.ExpectBegin()

				mock.ExpectExec("delete from relationship where following_id = \\? and follower_id = \\?").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("update account set following_count = following_count - 1 where id = \\?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("update account set followers_count = followers_count - 1 where id = \\?").
					WithArgs(2).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			isAuth:       true,
			urlParamFunc: func(r *http.Request) *http.Request { return setChiURLParam(r, "username", "unfollowUser") },
			wantCode:     http.StatusOK,
		},
		{
			name:     "Unauthorized",
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "user not found",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("followingUser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "followingUser"))
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("unfollowUser").
					WillReturnError(sql.ErrNoRows)
			},
			isAuth:       true,
			urlParamFunc: func(r *http.Request) *http.Request { return setChiURLParam(r, "username", "unfollowUser") },
			wantCode:     http.StatusNotFound,
		},
		{
			name: "relationship not found",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("followingUser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "followingUser"))
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("unfollowUser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(2, "unfollowUser"))

				mock.ExpectBegin()

				mock.ExpectExec("delete from relationship where following_id = \\? and follower_id = \\?").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			isAuth:       true,
			urlParamFunc: func(r *http.Request) *http.Request { return setChiURLParam(r, "username", "unfollowUser") },
			wantCode:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodDelete, "/v1/accounts/unfollowUser/follow", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.urlParamFunc != nil {
				r = tt.urlParamFunc(r)
			}
			if tt.isAuth {
				r.Header.Set("Authentication", "username followingUser")
			}
			if tt.mockFunc != nil {
				tt.mockFunc()
			}

			middleware := auth.Middleware(h.app)
			handlerMiddleware := middleware(http.HandlerFunc(h.Delete))
			handlerMiddleware.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)
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

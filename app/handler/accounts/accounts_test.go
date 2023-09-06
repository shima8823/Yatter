package accounts

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/object"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	db, mock := dao.NewMockDB()
	h := newMockHandler(db)
	defer db.Close()

	tests := []struct {
		name      string
		body      *AddRequest
		bodyBytes []byte
		mockFunc  func()
		wantCode  int
	}{
		{
			name: "successfully create account",
			body: &AddRequest{
				Username: "testuser",
				Password: "securepassword",
			},
			mockFunc: func() {
				mock.ExpectExec("insert into account \\(username, password_hash\\) values \\(\\?, \\?\\)").
					WithArgs("testuser", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantCode: http.StatusOK,
		},
		{
			name:      "bad request on malformed JSON",
			bodyBytes: []byte("{malformed}"),
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *http.Request
			var err error

			if tt.body != nil {
				tt.bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatal(err)
				}
			}

			r, err = http.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewReader(tt.bodyBytes))
			if err != nil {
				t.Fatal(err)
			}

			if tt.mockFunc != nil {
				tt.mockFunc()
			}

			w := httptest.NewRecorder()
			h.Create(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode == http.StatusOK {
				var resp object.Account
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatal(err)
				}
				assert.NotEmpty(t, resp.Username)
			}
		})
	}
}

func TestFindUserHandler(t *testing.T) {
	db, mock := dao.NewMockDB()
	h := newMockHandler(db)
	defer db.Close()

	tests := []struct {
		name         string
		username     string
		mockFunc     func()
		urlParamFunc func(r *http.Request) *http.Request
		wantCode     int
		wantUsername string
	}{
		{
			name:     "successfully find account",
			username: "testuser",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
			},
			urlParamFunc: func(r *http.Request) *http.Request { return setChiURLParam(r, "username", "testuser") },
			wantCode:     http.StatusOK,
			wantUsername: "testuser",
		},
		{
			name:     "user not found",
			username: "testuser",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnError(sql.ErrNoRows)
			},
			urlParamFunc: func(r *http.Request) *http.Request { return setChiURLParam(r, "username", "testuser") },
			wantCode:     http.StatusNotFound,
		},
		{
			name:     "bad request on invalid URL parameter",
			username: "testuser",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnError(sql.ErrNoRows)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/v1/accounts/"+tt.username, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.urlParamFunc != nil {
				r = tt.urlParamFunc(r)
			}
			if tt.mockFunc != nil {
				tt.mockFunc()
			}

			h.Get(w, r)

			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantCode == http.StatusOK {
				var resp object.Account
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantUsername, resp.Username)
			}
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

package statuses

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	db, mock := dao.NewMockDB()
	h := newMockHandler(db)
	defer db.Close()

	tests := []struct {
		name     string
		body     *AddRequest
		username string
		mockFunc func()
		wantCode int
	}{
		{
			name: "successfully create account",
			body: &AddRequest{
				Status: "test post",
			},
			username: "testuser",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
				mock.ExpectExec("insert into status \\(account_id, content\\) values \\(\\?, \\?\\)").
					WithArgs(1, "test post").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "unauthorized",
			body:     &AddRequest{Status: "test post"},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "bad request on malformed JSON",
			username: "testuser",
			mockFunc: func() {
				mock.ExpectQuery("select \\* from account where username = \\?").
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if tt.body != nil {
				var err error
				bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatal(err)
				}
			}

			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPost, "/v1/statuses", bytes.NewReader(bodyBytes))
			if err != nil {
				t.Fatal(err)
			}
			if tt.username != "" {
				r.Header.Set("Authentication", "username "+tt.username)
			}
			if tt.mockFunc != nil {
				tt.mockFunc()
			}

			middleware := auth.Middleware(h.app)
			handlerMiddleware := middleware(http.HandlerFunc(h.Create))
			handlerMiddleware.ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode == http.StatusOK {
				var resp object.Status
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.body.Status, resp.Content)
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

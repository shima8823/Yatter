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

	t.Run("successfully create account", func(t *testing.T) {
		body := &AddRequest{
			Status: "test post",
		}
		username := "testuser"

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, "/v1/statuses", bytes.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Authentication", "username "+username)

		// middlewareのFindbyUsernameのmock
		mock.ExpectQuery("select \\* from account where username = \\?").
			WithArgs("testuser").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
		// createStatusのmock
		mock.ExpectExec("insert into status \\(account_id, content\\) values \\(\\?, \\?\\)").
			WithArgs(1, body.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))

		middleware := auth.Middleware(h.app)
		handlerMiddleware := middleware(http.HandlerFunc(h.Create))
		handlerMiddleware.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code, "status code should be 200")
		var resp object.Status
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, body.Status, resp.Content, "status should be equal")
	})
	t.Run("unauthorized", func(t *testing.T) {
		body := &AddRequest{
			Status: "test post",
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/v1/statuses", bytes.NewReader(bodyBytes))
		respRecorder := httptest.NewRecorder()

		h.Create(respRecorder, req)

		assert.Equal(t, http.StatusUnauthorized, respRecorder.Code)
	})
	t.Run("bad request on malformed JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, "/v1/statuses", bytes.NewReader([]byte("{malformed")))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Authentication", "username "+"testuser")
		mock.ExpectQuery("select \\* from account where username = \\?").
			WithArgs("testuser").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
		middleware := auth.Middleware(h.app)
		handlerMiddleware := middleware(http.HandlerFunc(h.Create))
		handlerMiddleware.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func newMockHandler(db *sql.DB) *handler {
	return &handler{
		app: &app.App{
			Dao: dao.NewWithDB(sqlx.NewDb(db, "sqlmock")),
		},
	}
}

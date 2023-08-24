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
	h := &handler{
		app: &app.App{
			Dao: dao.NewWithDB(sqlx.NewDb(db, "sqlmock")),
		},
	}
	defer db.Close()

	t.Run("successfully create account", func(t *testing.T) {
		body := &AddRequest{
			Username: "testuser",
			Password: "securepassword",
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}

		// mockのsetup
		mock.ExpectExec("insert into account \\(username, password_hash\\) values \\(\\?, \\?\\)").
			WithArgs(body.Username, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		h.Create(w, r)

		assert.Equal(t, http.StatusOK, w.Code, "status code should be 200")
		var resp object.Account
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, body.Username, resp.Username, "username should be equal")
		assert.NotNil(t, resp.PasswordHash, "password hash should not be nil")
	})
	t.Run("bad request on malformed JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewReader([]byte("{malformed")))
		respRecorder := httptest.NewRecorder()

		h.Create(respRecorder, req)

		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
	})
}

func TestFindUserHandler(t *testing.T) {
	// mockのsetup
	db, mock := dao.NewMockDB()
	h := &handler{
		app: &app.App{
			Dao: dao.NewWithDB(sqlx.NewDb(db, "sqlmock")),
		},
	}
	defer db.Close()

	prepareRequest := func(username string) (*httptest.ResponseRecorder, *http.Request) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/v1/accounts/"+username, nil)
		if err != nil {
			t.Fatal(err)
		}
		r = setChiURLParam(r, "username", "testuser")
		return w, r
	}

	// 成功するテスト
	t.Run("successfully find account", func(t *testing.T) {
		w, r := prepareRequest("testuser")

		mock.ExpectQuery("select \\* from account where username = \\?").
			WithArgs("testuser").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))

		h.FindUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code, "status code should be 200")
		var resp object.Account
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "testuser", resp.Username, "username should be equal")
	})
	// 失敗するテスト
	// ユーザーが存在しない場合
	t.Run("user not found", func(t *testing.T) {
		w, r := prepareRequest("testuser")
		mock.ExpectQuery("select \\* from account where username = \\?").
			WithArgs("testuser").
			WillReturnError(sql.ErrNoRows)

		h.FindUser(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code, "status code should be 404")
	})
	// URLパラメータが不正な場合
	t.Run("bad request on invalid URL parameter", func(t *testing.T) {
		w, r := prepareRequest("testuser")
		r = setChiURLParam(r, "undefined", "testuser")
		mock.ExpectQuery("select \\* from account where username = \\?").
			WithArgs("testuser").
			WillReturnError(sql.ErrNoRows)

		h.FindUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "status code should be 404")
	})

}

func setChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

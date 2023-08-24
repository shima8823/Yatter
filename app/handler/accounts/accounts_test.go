package accounts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/object"

	"github.com/DATA-DOG/go-sqlmock"
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

package dao_test

import (
	"log"
	"os"
	"testing"
	"yatter-backend-go/app/domain/repository"
)

var accountRepo repository.Account
var statusRepo repository.Status
var relationshipRepo repository.Relationship
var cleanupDB func()

func TestMain(m *testing.M) {
	if dao, err := setupDAO(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		cleanupDB = func() {
			dao.InitAll()
		}
		defer dao.Close()

		accountRepo = dao.Account()
		statusRepo = dao.Status()
		relationshipRepo = dao.Relationship()
	}

	os.Exit(m.Run())
}

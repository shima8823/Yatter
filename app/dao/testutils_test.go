package dao_test

import (
	"fmt"
	"log"
	"os"
	"time"
	"yatter-backend-go/app/dao"

	"github.com/go-sql-driver/mysql"
)

// accessor namespace
var MySQL _mysql

type _mysql struct{}

// Read MySQL host
func (_mysql) host() string {
	v, err := getString("TEST_MYSQL_HOST")
	if err != nil {
		log.Fatal(err)
	}
	return v
}

// Read MySQL user
func (_mysql) user() string {
	v, err := getString("TEST_MYSQL_USER")
	if err != nil {
		log.Fatal(err)
	}
	return v
}

// Read MySQL password
func (_mysql) password() string {
	v, err := getString("TEST_MYSQL_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}
	return v
}

// Read MySQL database name
func (_mysql) database() string {
	v, err := getString("TEST_MYSQL_DATABASE")
	if err != nil {
		log.Fatal(err)
	}
	return v
}

// Read Timezone for MySQL
func (_mysql) location() *time.Location {
	tz, err := getString("TEST_MYSQL_TZ")
	if err != nil {
		return time.FixedZone("Asia/Tokyo", 9*60*60)
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatal(fmt.Errorf("Invalid timezone %+v", tz))
	}
	return loc
}

// Build mysql.Config
func testMySQLConfig() *mysql.Config {
	cfg := mysql.NewConfig()

	cfg.ParseTime = true
	cfg.Loc = MySQL.location()
	if host := MySQL.host(); host != "" {
		cfg.Net = "tcp"
		cfg.Addr = host
		log.Printf("Connecting to host: %s", host)
	}
	cfg.User = MySQL.user()
	cfg.Passwd = MySQL.password()
	cfg.DBName = MySQL.database()

	return cfg
}

func getString(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("config:[%s] not found", key)
	}
	return v, nil
}

func setupDAO() (dao.Dao, error) {
	testConf := testMySQLConfig()

	dao, err := dao.New(testConf)
	if err != nil {
		return nil, err
	}
	dao.InitAll()
	return dao, nil
}

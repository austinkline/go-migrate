package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	
	"github.com/golang/glog"

	// import the driver we're doing to be using.
	_ "github.com/go-sql-driver/mysql"
)

var (
	// DBMultiStatement - global database object which can issue multiple statements at once
	DBMultiStatement *sql.DB

	// ErrEmptyEnv - used for telling if an environment variable failed to be loaded.
	ErrEmptyEnv = errors.New("environment variable cannot be empty")
)

const (
	driverName  = "mysql"
	userEnv     = "MIGRATE_DB_USER"
	hostEnv     = "MIGRATE_DB_HOST"
	passwordEnv = "MIGRATE_DB_PASSWORD"
	portEnv     = "MIGRATE_DB_PORT"
)

// SetupMultistatementDBWithEnv - setup to allow multiple statements per query
func SetupMultistatementDBWithEnv() (err error){
	user := os.Getenv(userEnv)
	host := os.Getenv(hostEnv)
	password := os.Getenv(passwordEnv)
	port := os.Getenv(portEnv)

	if user == "" || host == "" || password == "" || port == "" {
		err = ErrEmptyEnv
		return
	}

	err = SetupDB(user, host, password, port, true /* multiStatements */)
	return
}

// SetupDB - Setups up the DB with a tcp connection\][']
func SetupDB(user, host, password, port string, multiStatements bool) (err error) {
	glog.Infof("setting up db connection...")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/sys?multiStatements=%v", user, password, host, port, multiStatements)
	DBMultiStatement, err = sql.Open(driverName, connection)

	err = DBMultiStatement.Ping()
	if err != nil {
		return
	}

	glog.Info("connected sucessfully!")
	return
}

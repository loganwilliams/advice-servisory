package types

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/loganwilliams/adviceservisory/server/env"
	"github.com/stretchr/testify/suite"
)

var (
	testDB *sql.DB
)

// package level test entry point
func TestMain(m *testing.M) {
	var (
		retCode int
		err     error
	)

	err = setup()
	if err != nil {
		log.Panic("Error setting up test db", err)
	}

	retCode = m.Run()

	err = teardown()
	if err != nil {
		log.Panic("Error tearing down test db", err)
	}

	os.Exit(retCode)
}

func setup() error {
	var (
		conf *env.TestConfig
		err  error
	)

	// get test config
	conf = env.NewTestConfig()

	// construct database info string required for connection
	dbInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.DB.Host,
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Name,
		"disable",
	)

	// open a connection to the database
	testDB, err = sql.Open(env.DB_DRIVER, dbInfo)
	if err != nil {
		log.Panic("Error connecting to test db", err)
	}

	// ping open db to verify the connection has been established.
	err = testDB.Ping()
	if err != nil {
		log.Panic("Error pinging test db", err)
	}

	return nil
}

func teardown() error {
	// destroy test databas tables
	_, err := testDB.Exec(`DROP TABLE IF EXISTS routes CASCADE;
						   DROP TABLE IF EXISTS import.routes CASCADE;
						   DROP SCHEMA IF EXISTS import;`)

	if err != nil {
		log.Panic("Error tearing down test db", err)
	}

	return nil
}

func TestRouteSuite(t *testing.T) {
	suite.Run(t, &RouteTestSuite{db: testDB})
}

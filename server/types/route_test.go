package types

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"log"
	_ "github.com/lib/pq"
)

type RouteTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *RouteTestSuite) SetupSuite() {
	// remove whatever existing tables we have

	err := DropRoutesTable(s.db)
	if err != nil {
		log.Fatal("Error tearing down routes test suite", err)
	}
}

func (s *RouteTestSuite) TearDownSuite() {
	// remove whatever tables we created

	err := DropRoutesTable(s.db)
	if err != nil {
		log.Fatal("Error tearing down routes test suite", err)
	}
}

func (s *RouteTestSuite) TestCreateAndPopulateRoutesTable() {
	err := CreateAndPopulateRoutesTable(s.db)
	s.NoError(err)
}

func (s *RouteTestSuite) TestReadRoutes() {
	routes, err := ReadRoutes(s.db)
	s.NoError(err)

	s.Equal(len(routes), 29)
}
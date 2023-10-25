package testutil

import (
	"database/sql"

	"github.com/stretchr/testify/suite"
)

type PostgresSQLDockerSQLSuite struct {
	suite.Suite

	container *PostgreSQLContainer
	rootDB    *sql.DB
	db        *sql.Tx
}

func (s *PostgresSQLDockerSQLSuite) SetupSuite() {
	container, err := NewPostgreSQLContainer(
		PostgreSQLContainerConfig{ExpireSeconds: 60},
	)
	s.Require().Nil(err, "failed to start container: %v", err)

	db, err := container.ConnectSQL()
	s.Require().Nil(err, "failed to connect with sql: %v", err)

	s.container = container
	s.rootDB = db
}

func (s *PostgresSQLDockerSQLSuite) TearDownSuite() {
	err := s.container.Close()
	s.Require().Nil(err, "failed to shutdown container: %v", err)
}

func (s *PostgresSQLDockerSQLSuite) SetupTest() {
	s.Require().Nil(s.db)

	var err error
	s.db, err = s.rootDB.Begin()
	s.Require().Nil(err)
}

func (s *PostgresSQLDockerSQLSuite) TearDownTest() {
	err := s.db.Rollback()
	s.Assert().Nil(err)
	s.db = nil
}

func (s *PostgresSQLDockerSQLSuite) DB() *sql.Tx {
	return s.db
}

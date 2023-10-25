package testutil

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresSQLDockerGORMSuite struct {
	suite.Suite

	container *PostgreSQLContainer
	rootDB    *gorm.DB
	db        *gorm.DB
}

func (s *PostgresSQLDockerGORMSuite) SetupSuite() {
	container, err := NewPostgreSQLContainer(
		PostgreSQLContainerConfig{ExpireSeconds: 60},
	)
	s.Require().Nil(err, "failed to start container: %v", err)

	db, err := container.ConnectGORM(
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)
	s.Require().Nil(err, "failed to connect with gorm: %v", err)

	s.container = container
	s.rootDB = db
}

func (s *PostgresSQLDockerGORMSuite) TearDownSuite() {
	err := s.container.Close()
	s.Require().Nil(err, "failed to shutdown container: %v", err)
}

func (s *PostgresSQLDockerGORMSuite) SetupTest() {
	s.Require().Nil(s.db)
	s.db = s.rootDB.Begin()
	s.Require().Nil(s.db.Error)
}

func (s *PostgresSQLDockerGORMSuite) TearDownTest() {
	err := s.db.Rollback().Error
	s.Assert().Nil(err)
	s.db = nil
}

func (s *PostgresSQLDockerGORMSuite) DB() *gorm.DB {
	return s.db
}

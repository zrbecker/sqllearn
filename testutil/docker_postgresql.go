package testutil

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLContainer struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

type PostgreSQLContainerConfig struct {
	ExpireSeconds uint
}

func NewPostgreSQLContainer(config PostgreSQLContainerConfig) (*PostgreSQLContainer, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	if err = pool.Client.Ping(); err != nil {
		return nil, err
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15.2",
			Env:        []string{"POSTGRES_PASSWORD=postgres"},
		},
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return nil, err
	}

	if err := resource.Expire(config.ExpireSeconds); err != nil {
		if purgeErr := pool.Purge(resource); purgeErr != nil {
			return nil, errors.Join(err, purgeErr)
		}
		return nil, err
	}

	return &PostgreSQLContainer{pool: pool, resource: resource}, nil
}

func (c *PostgreSQLContainer) ConnectSQL() (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"postgres",
		"postgres",
		"localhost",
		c.resource.GetPort("5432/tcp"),
		"postgres",
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := c.pool.Retry(func() error {
		return db.Ping()
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func (c *PostgreSQLContainer) ConnectGORM(opts ...gorm.Option) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		"localhost",
		"postgres",
		"postgres",
		"postgres",
		c.resource.GetPort("5432/tcp"),
	)

	var db *gorm.DB
	if err := c.pool.Retry(func() error {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), opts...)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func (c *PostgreSQLContainer) Close() error {
	if err := c.pool.Purge(c.resource); err != nil {
		return err
	}
	return nil
}

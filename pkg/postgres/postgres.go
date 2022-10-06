package postgres

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(ctx context.Context, conf Config) (*gorm.DB, error) {
	postgresDialector := postgres.Open(conf.DSN())
	conn, err := gorm.Open(postgresDialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

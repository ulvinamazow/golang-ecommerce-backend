package config

import (
	"context"
	"fmt"
)

func GetConnectionPool(ctx context.Context, config Config) /**pgxpool.Pool*/ string /*, error*/ {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
	)

	/*
		pool, err := pgxpool.New(ctx, connString)
		if err != nil {
			return nil, fmt.Errorf("Connecting pool can not created: %w", err)
		}
	*/

	return connString //, nil
}

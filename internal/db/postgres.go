package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/narvdeshwar/IOT-rag/internal/logger"
)

var Pool *pgxpool.Pool

func Init(ctx context.Context, dsn string) {
	var err error
	Pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		logger.L.Error("Postgres Connection Failed", "err", err)
		panic(err)
	}
	logger.L.Info("Postgres connected successfully!")
}

package repositories

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/radon1/pg-stat-test-task/internal/repositories/pgstat"
)

type Container struct {
	PGStat pgstat.Repository
}

func New(pool *pgxpool.Pool) *Container {
	return &Container{
		PGStat: pgstat.NewRepository(pool),
	}
}

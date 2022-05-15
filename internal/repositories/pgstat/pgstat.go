package pgstat

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/radon1/pg-stat-test-task/internal/models"
)

type Repository interface {
	FindQueriesStat(ctx context.Context, limit, offset int, queryFilters []string) ([]models.QueryStat, error)
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) FindQueriesStat(ctx context.Context, limit, offset int, queryFilters []string) ([]models.QueryStat, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	queryArgs := make([]interface{}, 0, 3)
	queryArgs = append(queryArgs, limit)
	queryArgs = append(queryArgs, offset)
	var queryFilter string
	if len(queryFilters) > 0 {
		queryFilter = "and lower(query) like any ($3)"
		queryFilterArgs, err := buildLikeArgsArray(queryFilters)
		if err != nil {
			return nil, err
		}
		queryArgs = append(queryArgs, queryFilterArgs)
	}

	rows, err := conn.Query(ctx, fmt.Sprintf(findQueriesStatQuery, queryFilter), queryArgs...)
	if err != nil {
		return nil, err
	}

	var queriesStatModels []models.QueryStat
	for rows.Next() {
		var model models.QueryStat
		err := rows.Scan(
			&model.Query,
			&model.Calls,
			&model.TotalTime,
			&model.MeanTime,
			&model.Percentage,
		)
		if err != nil {
			return nil, err
		}
		queriesStatModels = append(queriesStatModels, model)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return queriesStatModels, nil
}

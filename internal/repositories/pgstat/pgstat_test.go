package pgstat

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/radon1/pg-stat-test-task/internal/config"
	"github.com/radon1/pg-stat-test-task/internal/system/connections/pg"
	"github.com/radon1/pg-stat-test-task/test/fixtures"
)

func getTestRepoAndPool(t *testing.T, ctx context.Context) (Repository, *pgxpool.Pool) {
	cfg, err := config.New()
	assert.NoError(t, err)

	pool, err := pg.NewPool(ctx, cfg.PGDSN)
	assert.NoError(t, err)

	return NewRepository(pool), pool
}

func execQuery(t *testing.T, ctx context.Context, pool *pgxpool.Pool, query string) {
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	_, err = conn.Exec(context.Background(), query)
	assert.NoError(t, err)
}

func TestRepository_FindQueriesStat(t *testing.T) {
	if testing.Short() {
		t.Skipf("long test, skip")
	}

	cases := []struct {
		name          string
		prepareData   func(t *testing.T, ctx context.Context, pool *pgxpool.Pool)
		expectedCount int
		limit, offset int
		queryFilters  []string
		isError       bool
		errMessage    string
	}{
		{
			name: "check exists queries",
			prepareData: func(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
				execQuery(t, ctx, pool, "INSERT INTO orders (name) VALUES ('test1')")
				execQuery(t, ctx, pool, "UPDATE orders SET name = 'test2'")
			},
			expectedCount: 1,
			limit:         10,
			offset:        0,
			queryFilters:  []string{"update"},
		},
		{
			name: "invalid offset",
			prepareData: func(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
				execQuery(t, ctx, pool, "INSERT INTO orders (name) VALUES ('test1')")
				execQuery(t, ctx, pool, "UPDATE orders SET name = 'test2'")
			},
			expectedCount: 0,
			limit:         10,
			offset:        11,
			queryFilters:  []string{"update"},
		},
		{
			name: "a lot of queries",
			prepareData: func(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
				execQuery(t, ctx, pool, "INSERT INTO orders (name) VALUES ('test1')")
				execQuery(t, ctx, pool, "UPDATE orders SET name = 'test2'")
				execQuery(t, ctx, pool, "UPDATE orders SET name = 'test3' WHERE name = 'test2'")
				execQuery(t, ctx, pool, "UPDATE orders SET name = 'test4' WHERE created_at <> NOW()")
			},
			expectedCount: 3,
			limit:         10,
			offset:        0,
			queryFilters:  []string{"update"},
		},
		{
			name: "with different query filters",
			prepareData: func(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
				execQuery(t, ctx, pool, "INSERT INTO orders (name) VALUES ('test1')")
				execQuery(t, ctx, pool, "UPDATE orders SET name = 'test2'")
				execQuery(t, ctx, pool, "UPDATE orders SET name = 'test3' WHERE name = 'test2'")
				execQuery(t, ctx, pool, "DELETE FROM orders")
			},
			expectedCount: 3,
			limit:         10,
			offset:        0,
			queryFilters:  []string{"update", "delete"},
		},
		{
			name: "with not found queries",
			prepareData: func(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
				execQuery(t, ctx, pool, "INSERT INTO orders (name) VALUES ('test1')")
				execQuery(t, ctx, pool, "DELETE FROM orders")
			},
			expectedCount: 0,
			limit:         10,
			offset:        0,
			queryFilters:  []string{"update"},
		},
		{
			name:          "malformed filters",
			prepareData:   func(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {},
			expectedCount: 0,
			limit:         10,
			offset:        0,
			queryFilters:  []string{"update", ""},
			isError:       true,
			errMessage:    malformedArgsErr.Error(),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				ctx              = context.Background()
				repository, pool = getTestRepoAndPool(t, ctx)
			)
			defer pool.Close()

			fixtures.ExecuteFixture(t, ctx, pool, fixtures.CleanupFixture{})
			fixtures.ExecuteFixture(t, ctx, pool, statementsResetFixture{})
			testCase.prepareData(t, ctx, pool)

			stats, err := repository.FindQueriesStat(ctx, testCase.limit, testCase.offset, testCase.queryFilters)
			if testCase.isError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.expectedCount, len(stats))
		})
	}
}

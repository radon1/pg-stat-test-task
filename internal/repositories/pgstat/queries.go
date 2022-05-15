package pgstat

const findQueriesStatQuery = `
SELECT
	query,
	calls,
	round(total_exec_time::numeric, 2) AS total_time,
	round(mean_exec_time::numeric, 2) AS mean_time,
	round((100 * total_exec_time / sum(total_exec_time) OVER ())::numeric, 2) AS percentage
FROM pg_stat_statements
WHERE dbid = (select dbid from pg_database where datname = 'stat_test_task') %s
ORDER BY mean_exec_time DESC
LIMIT $1
OFFSET $2;
`

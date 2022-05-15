package pgstat

type statementsResetFixture struct{}

func (f statementsResetFixture) GetSql() []string {
	return []string{
		`SELECT pg_stat_statements_reset();`,
	}
}

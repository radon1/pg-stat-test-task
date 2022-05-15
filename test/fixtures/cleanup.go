package fixtures

type CleanupFixture struct{}

func (f CleanupFixture) GetSql() []string {
	return []string{
		`TRUNCATE TABLE orders RESTART IDENTITY CASCADE;`,
	}
}

CREATE TABLE IF NOT EXISTS orders (
    name text not null,
    created_at timestamp not null default now()
);

CREATE EXTENSION pg_stat_statements;

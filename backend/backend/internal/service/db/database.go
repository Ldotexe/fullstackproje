//go:generate mockgen -source ./database.go -destination=./mocks/database.go -package=mock_db
package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBops interface {
	GetPool(_ context.Context) *pgxpool.Pool
	Get(ctx context.Context, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type Database struct {
	cluster *pgxpool.Pool
}

func newDatabase(cluster *pgxpool.Pool) *Database {
	return &Database{cluster: cluster}
}

func (db Database) GetPool(_ context.Context) *pgxpool.Pool {
	return db.cluster
}

func (db Database) Get(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Get(ctx, db.cluster, dest, query, args...)
}

func (db Database) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, db.cluster, dest, query, args...)
}

func (db Database) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}

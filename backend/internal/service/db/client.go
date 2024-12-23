package db

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = os.Getenv("POSTGRES_PORT")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PW")
	dbname   = os.Getenv("POSTGRES_DB")
)

func NewDB(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn())
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}

func (db Database) CloseDB(ctx context.Context) {
	db.GetPool(ctx).Close()
}

func generateDsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname,
	)
}

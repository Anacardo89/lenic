package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbHandler struct {
	pool *pgxpool.Pool
}

func NewDBRepo(pool *pgxpool.Pool) DBRepo {
	return &dbHandler{
		pool: pool,
	}
}

func (db *dbHandler) Close() {
	db.pool.Close()
}

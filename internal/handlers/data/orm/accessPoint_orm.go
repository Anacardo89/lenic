package orm

import "database/sql"

var (
	Da DataAccess
)

type DataAccess struct {
	db *sql.DB
}

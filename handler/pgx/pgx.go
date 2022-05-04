package pgx

import (
	"context"

	"github.com/jackc/pgconn"
)

//go:generate mockgen -destination=../../mocks/mockDbConn.go -package=mocks example/web-service-gin/handler/pgx DbConn
//go:generate mockgen -destination=../../mocks/mockRows.go -package=mocks example/web-service-gin/handler/pgx Rows

type DbConn interface {
	Query(
		ctx context.Context,
		sql string,
		args ...interface{},
	) (Rows, error)
	Exec(
		ctx context.Context,
		sql string,
		arguments ...interface{},
	) (pgconn.CommandTag, error)
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

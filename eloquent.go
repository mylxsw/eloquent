package eloquent

import (
	"context"

	"github.com/mylxsw/eloquent/query"
)

type Database interface {
	Query(ctx context.Context, builder QueryBuilder, cb func(row Scanner) (any, error)) ([]any, error)
	Insert(ctx context.Context, tableName string, kv query.KV) (int64, error)
	Delete(ctx context.Context, builder query.SQLBuilder) (int64, error)
	Update(ctx context.Context, builder query.SQLBuilder, kv query.KV) (int64, error)
	Statement(ctx context.Context, raw string, args ...any) error
}

type QueryBuilder interface {
	ResolveQuery() (sqlStr string, args []any)
}

// Scanner is an interface which wraps sql.Rows's Scan method
type Scanner interface {
	Scan(dest ...any) error
}

package eloquent

import (
	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent/query"
)

type Database interface {
	Query(builder QueryBuilder, cb func(row Scanner) (interface{}, error)) (*coll.Collection, error)
	Insert(tableName string, kv query.KV) (int64, error)
	Delete(builder query.SQLBuilder) (int64, error)
	Update(builder query.SQLBuilder, kv query.KV) (int64, error)
	Statement(raw string, args ...interface{}) error
}

type QueryBuilder interface {
	ResolveQuery() (sqlStr string, args []interface{})
}

// Scanner is an interface which wraps sql.Rows's Scan method
type Scanner interface {
	Scan(dest ...interface{}) error
}

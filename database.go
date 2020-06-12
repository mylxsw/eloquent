package eloquent

import (
	"context"
	"database/sql"

	"github.com/mylxsw/coll"
	"github.com/mylxsw/eloquent/query"
)

// Build create a SQLBuilder with table name
func Build(tableName string) query.SQLBuilder {
	return query.Builder().Table(tableName)
}

// databaseImpl is a basic database query handler
type databaseImpl struct {
	db *query.DatabaseWrap
}

// DB create a databaseImpl
func DB(db query.Database) Database {
	return &databaseImpl{
		db: query.NewDatabaseWrap(db),
	}
}

type rawQueryBuilder struct {
	sql  string
	args []interface{}
}

func Raw(sqlStr string, args ...interface{}) QueryBuilder {
	return &rawQueryBuilder{sql: sqlStr, args: args}
}

func (r *rawQueryBuilder) ResolveQuery() (sqlStr string, args []interface{}) {
	return r.sql, r.args
}

// Query run a basic query
func (db *databaseImpl) Query(builder QueryBuilder, cb func(row Scanner) (interface{}, error)) (*coll.Collection, error) {
	results := make([]interface{}, 0)

	sqlStr, args := builder.ResolveQuery()
	rows, err := db.db.QueryContext(context.TODO(), sqlStr, args...)
	if err != nil {
		return coll.MustNew(results), err
	}

	defer rows.Close()

	for rows.Next() {
		r, err := cb(rows)
		if err != nil {
			return coll.MustNew(results), err
		}

		results = append(results, r)
	}

	return coll.MustNew(results), nil
}

// Insert to execute an insert statement
func (db *databaseImpl) Insert(tableName string, kv query.KV) (int64, error) {
	sqlStr, args := query.Builder().Table(tableName).ResolveInsert(kv)
	res, err := db.db.ExecContext(context.TODO(), sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Delete to execute an delete statement
func (db *databaseImpl) Delete(builder query.SQLBuilder) (int64, error) {
	sqlStr, args := builder.ResolveDelete()
	res, err := db.db.ExecContext(context.TODO(), sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Update to execute an update statement
func (db *databaseImpl) Update(builder query.SQLBuilder, kv query.KV) (int64, error) {
	sqlStr, args := builder.ResolveUpdate(kv)
	res, err := db.db.ExecContext(context.TODO(), sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Statement running a general statement which return no value
func (db *databaseImpl) Statement(raw string, args ...interface{}) error {
	_, err := db.db.ExecContext(context.TODO(), raw, args...)
	return err
}

// Transaction start a transaction
func Transaction(db *sql.DB, cb func(tx query.Database) error) (err error) {
	return query.Transaction(db, cb)
}

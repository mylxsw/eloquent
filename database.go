package eloquent

import (
	"context"
	"database/sql"

	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/go-toolkit/collection"
)

// Build create a SQLBuilder with table name
func Build(tableName string) query.SQLBuilder {
	return query.Builder().Table(tableName)
}

// Database is a basic database query handler
type Database struct {
	db *query.DatabaseWrap
}

// DB create a Database
func DB(db query.Database) *Database {
	return &Database{
		db: query.NewDatabaseWrap(db),
	}
}

// Query run a basic query
func (db *Database) Query(builder query.SQLBuilder, cb func(row *sql.Rows) (interface{}, error)) (*collection.Collection, error) {
	results := make([]interface{}, 0)

	sqlStr, args := builder.ResolveQuery()
	rows, err := db.db.QueryContext(context.TODO(), sqlStr, args...)
	if err != nil {
		return collection.MustNew(results), err
	}

	defer rows.Close()

	for rows.Next() {
		r, err := cb(rows)
		if err != nil {
			return collection.MustNew(results), err
		}

		results = append(results, r)
	}

	return collection.MustNew(results), nil
}

// Insert to execute an insert statement
func (db *Database) Insert(tableName string, kv query.KV) (int64, error) {
	sqlStr, args := query.Builder().Table(tableName).ResolveInsert(kv)
	res, err := db.db.ExecContext(context.TODO(), sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Delete to execute an delete statement
func (db *Database) Delete(builder query.SQLBuilder) (int64, error) {
	sqlStr, args := builder.ResolveDelete()
	res, err := db.db.ExecContext(context.TODO(), sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Update to execute an update statement
func (db *Database) Update(builder query.SQLBuilder, kv query.KV) (int64, error) {
	sqlStr, args := builder.ResolveUpdate(kv)
	res, err := db.db.ExecContext(context.TODO(), sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Statement running a general statement which return no value
func (db *Database) Statement(raw string, args ...interface{}) error {
	_, err := db.db.ExecContext(context.TODO(), raw, args...)
	return err
}

// Transaction start a transaction
func Transaction(db *sql.DB, cb func(tx query.Database) error) (err error) {
	return query.Transaction(db, cb)
}

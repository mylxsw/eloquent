package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mylxsw/eloquent/event"
)

type Database interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type DatabaseWrap struct {
	db Database
}

func NewDatabaseWrap(db Database) *DatabaseWrap {
	return &DatabaseWrap{db: db}
}

func (d *DatabaseWrap) GetDB() Database {
	return d.db
}

func (d *DatabaseWrap) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	startTs := time.Now()
	defer func() {
		event.Dispatch(event.QueryExecutedEvent{
			SQL:      query,
			Bindings: args,
			Time:     time.Now().Sub(startTs),
		})
	}()

	return d.db.ExecContext(ctx, query, args...)
}

func (d *DatabaseWrap) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	startTs := time.Now()
	defer func() {
		event.Dispatch(event.QueryExecutedEvent{
			SQL:      query,
			Bindings: args,
			Time:     time.Now().Sub(startTs),
		})
	}()

	return d.db.QueryContext(ctx, query, args...)
}

// Transaction create a transaction with auto commit support
func Transaction(db *sql.DB, cb func(tx Database) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	event.Dispatch(event.TransactionBeginningEvent{})

	defer func() {
		if err2 := recover(); err2 != nil {
			if err3 := tx.Rollback(); err3 != nil {
				err = fmt.Errorf("rollback (%s) failed: %s", err2, err3)
			} else {
				err = fmt.Errorf("rollback (%s)", err2)
				event.Dispatch(event.TransactionRolledBackEvent{})
			}
		}
	}()

	if err := cb(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf("rollback (%s) failed: %s", err, err2)
		}

		event.Dispatch(event.TransactionRolledBackEvent{})

		return fmt.Errorf("rollback (%s)", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %s", err)
	}

	event.Dispatch(event.TransactionCommittedEvent{})

	return nil
}

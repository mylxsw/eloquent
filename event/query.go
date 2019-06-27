package event

import (
	"time"
)

// QueryExecutedEvent each SQL query executed by your application will publish a QueryExecutedEvent
type QueryExecutedEvent struct {
	SQL      string
	Bindings []interface{}
	Time     time.Duration
}

// TransactionBeginningEvent fired when a new transaction started
type TransactionBeginningEvent struct{}

// TransactionCommittedEvent fired when a transaction has been committed
type TransactionCommittedEvent struct{}

// TransactionRolledBackEvent fired when a transaction has been rollback
type TransactionRolledBackEvent struct{}

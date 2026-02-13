package query

import (
	"fmt"
	"strings"
)

// Dialect defines database-specific SQL syntax differences.
type Dialect interface {
	// QuoteIdentifier wraps an identifier (column, table) with the appropriate quote character.
	QuoteIdentifier(name string) string
	// Placeholder returns the parameter placeholder for the given 1-based index.
	Placeholder(index int) string
	// SupportsLastInsertId indicates if the driver supports LastInsertId.
	// PostgreSQL does not; it requires RETURNING clause instead.
	SupportsLastInsertId() bool
}

// MySQLDialect implements MySQL-specific SQL syntax.
type MySQLDialect struct{}

func (d MySQLDialect) QuoteIdentifier(name string) string {
	return "`" + name + "`"
}

func (d MySQLDialect) Placeholder(_ int) string {
	return "?"
}

func (d MySQLDialect) SupportsLastInsertId() bool {
	return true
}

// PostgreSQLDialect implements PostgreSQL-specific SQL syntax.
type PostgreSQLDialect struct{}

func (d PostgreSQLDialect) QuoteIdentifier(name string) string {
	return `"` + name + `"`
}

func (d PostgreSQLDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

func (d PostgreSQLDialect) SupportsLastInsertId() bool {
	return false
}

var globalDialect Dialect = MySQLDialect{}

// replacePlaceholders replaces ? placeholders with dialect-specific ones.
// For MySQL this is a no-op; for PostgreSQL it converts ? to $1, $2, etc.
func replacePlaceholders(sql string) string {
	if _, ok := globalDialect.(MySQLDialect); ok {
		return sql
	}

	var result strings.Builder
	idx := 1
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' {
			result.WriteString(globalDialect.Placeholder(idx))
			idx++
		} else {
			result.WriteByte(sql[i])
		}
	}
	return result.String()
}

// makePlaceholders generates a comma-separated list of placeholders for the given count.
func makePlaceholders(count int) string {
	if count == 0 {
		return ""
	}
	parts := make([]string, count)
	for i := 0; i < count; i++ {
		parts[i] = "?"
	}
	return strings.Join(parts, ", ")
}

// SetDialect sets the global dialect used by all query builders.
func SetDialect(d Dialect) {
	globalDialect = d
}

// GetDialect returns the current global dialect.
func GetDialect() Dialect {
	return globalDialect
}

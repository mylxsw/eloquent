package query

import (
	"testing"
)

func TestMySQLDialect_QuoteIdentifier(t *testing.T) {
	d := MySQLDialect{}
	if d.QuoteIdentifier("name") != "`name`" {
		t.Error("MySQL QuoteIdentifier failed")
	}
}

func TestMySQLDialect_Placeholder(t *testing.T) {
	d := MySQLDialect{}
	if d.Placeholder(1) != "?" {
		t.Error("MySQL Placeholder failed")
	}
	if d.Placeholder(5) != "?" {
		t.Error("MySQL Placeholder should always return ?")
	}
}

func TestMySQLDialect_SupportsLastInsertId(t *testing.T) {
	d := MySQLDialect{}
	if !d.SupportsLastInsertId() {
		t.Error("MySQL should support LastInsertId")
	}
}

func TestPostgreSQLDialect_QuoteIdentifier(t *testing.T) {
	d := PostgreSQLDialect{}
	if d.QuoteIdentifier("name") != `"name"` {
		t.Error("PostgreSQL QuoteIdentifier failed")
	}
}

func TestPostgreSQLDialect_Placeholder(t *testing.T) {
	d := PostgreSQLDialect{}
	if d.Placeholder(1) != "$1" {
		t.Error("PostgreSQL Placeholder(1) failed")
	}
	if d.Placeholder(3) != "$3" {
		t.Error("PostgreSQL Placeholder(3) failed")
	}
}

func TestPostgreSQLDialect_SupportsLastInsertId(t *testing.T) {
	d := PostgreSQLDialect{}
	if d.SupportsLastInsertId() {
		t.Error("PostgreSQL should not support LastInsertId")
	}
}

func TestReplacePlaceholders_MySQL(t *testing.T) {
	SetDialect(MySQLDialect{})
	defer SetDialect(MySQLDialect{})

	sql := "SELECT * FROM users WHERE id = ? AND name = ?"
	result := replacePlaceholders(sql)
	if result != sql {
		t.Errorf("MySQL replacePlaceholders should be no-op, got %s", result)
	}
}

func TestReplacePlaceholders_PostgreSQL(t *testing.T) {
	SetDialect(PostgreSQLDialect{})
	defer SetDialect(MySQLDialect{})

	sql := "SELECT * FROM users WHERE id = ? AND name = ?"
	expected := "SELECT * FROM users WHERE id = $1 AND name = $2"
	result := replacePlaceholders(sql)
	if result != expected {
		t.Errorf("PostgreSQL replacePlaceholders failed, expected %s, got %s", expected, result)
	}
}

func TestMakePlaceholders(t *testing.T) {
	if makePlaceholders(0) != "" {
		t.Error("makePlaceholders(0) should be empty")
	}
	if makePlaceholders(1) != "?" {
		t.Error("makePlaceholders(1) failed")
	}
	if makePlaceholders(3) != "?, ?, ?" {
		t.Errorf("makePlaceholders(3) failed, got %s", makePlaceholders(3))
	}
}

func TestSetDialect(t *testing.T) {
	original := GetDialect()
	defer SetDialect(original)

	SetDialect(PostgreSQLDialect{})
	d := GetDialect()
	if _, ok := d.(PostgreSQLDialect); !ok {
		t.Error("SetDialect failed for PostgreSQL")
	}

	SetDialect(MySQLDialect{})
	d = GetDialect()
	if _, ok := d.(MySQLDialect); !ok {
		t.Error("SetDialect failed for MySQL")
	}
}

package migrate

import (
	"fmt"
	"strings"
)

// MigrateDialect defines database-specific DDL syntax for migrations.
type MigrateDialect interface {
	QuoteIdentifier(name string) string
	WrapTable(prefix, name string) string
	CompileCreateTable(temporary bool, ifNotExists bool, tableName string, columns []string) string
	CompileCreateEncoding(charset, collation string) string
	CompileCreateEngine(engine string) string
	CompileAdd(tableName string, columns []string) string
	CompileChange(tableName string, columns []string) string
	CompileDrop(tableName string) string
	CompileDropIfExists(tableName string) string
	CompileDropColumn(tableName string, columns []string) string
	CompileDropPrimary(tableName string) string
	CompileDropIndex(tableName string, indexName string) string
	CompileRename(tableName string, to string) string
	CompileRenameColumn(tableName string, from, to string) string
	CompileKey(tableName string, c *Command, indexType string) string
	CompileForeign(tableName string, c *Command) string
	CompileDropForeign(tableName string, name string) string
	CompileColumnType(col *ColumnDefinition) string
	Columnize(columns []string) string
}

// MySQLMigrateDialect implements MySQL-specific DDL.
type MySQLMigrateDialect struct{}

func (d MySQLMigrateDialect) QuoteIdentifier(name string) string {
	return "`" + name + "`"
}

func (d MySQLMigrateDialect) WrapTable(prefix, name string) string {
	return "`" + prefix + name + "`"
}

func (d MySQLMigrateDialect) CompileCreateTable(temporary bool, ifNotExists bool, tableName string, columns []string) string {
	createStatement := "CREATE"
	if temporary {
		createStatement += " TEMPORARY"
	}

	ifNotExistsStr := ""
	if ifNotExists {
		ifNotExistsStr = "IF NOT EXISTS "
	}

	return fmt.Sprintf(
		"%s TABLE %s%s (%s)",
		createStatement,
		ifNotExistsStr,
		tableName,
		strings.Join(columns, ", "),
	)
}

func (d MySQLMigrateDialect) CompileCreateEncoding(charset, collation string) string {
	sqlStr := ""
	if charset != "" {
		sqlStr += " DEFAULT CHARACTER SET " + charset
	}
	if collation != "" {
		sqlStr += " COLLATE " + collation
	}
	return sqlStr
}

func (d MySQLMigrateDialect) CompileCreateEngine(engine string) string {
	if engine != "" {
		return " ENGINE = " + engine
	}
	return ""
}

func (d MySQLMigrateDialect) CompileAdd(tableName string, columns []string) string {
	var cols = make([]string, len(columns))
	for i, c := range columns {
		cols[i] = "ADD " + c
	}
	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(cols, ", "))
}

func (d MySQLMigrateDialect) CompileChange(tableName string, columns []string) string {
	var cols = make([]string, len(columns))
	for i, c := range columns {
		cols[i] = "MODIFY " + c
	}
	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(cols, ", "))
}

func (d MySQLMigrateDialect) CompileDrop(tableName string) string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}

func (d MySQLMigrateDialect) CompileDropIfExists(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
}

func (d MySQLMigrateDialect) CompileDropColumn(tableName string, columns []string) string {
	dropColumns := make([]string, len(columns))
	for i, cc := range columns {
		dropColumns[i] = "DROP " + cc
	}
	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(dropColumns, ", "))
}

func (d MySQLMigrateDialect) CompileDropPrimary(tableName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY", tableName)
}

func (d MySQLMigrateDialect) CompileDropIndex(tableName string, indexName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP INDEX `%s`", tableName, indexName)
}

func (d MySQLMigrateDialect) CompileRename(tableName string, to string) string {
	return fmt.Sprintf("RENAME TABLE %s TO `%s`", tableName, to)
}

func (d MySQLMigrateDialect) CompileRenameColumn(tableName string, from, to string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN `%s` TO `%s`", tableName, from, to)
}

func (d MySQLMigrateDialect) CompileKey(tableName string, c *Command, indexType string) string {
	alg := ""
	if c.CommandAlgorithm != "" {
		alg = " USING " + c.CommandAlgorithm
	}
	return fmt.Sprintf(
		"ALTER TABLE %s ADD %s %s%s(%s)",
		tableName,
		indexType,
		c.CommandIndex,
		alg,
		d.Columnize(c.CommandParameters),
	)
}

func (d MySQLMigrateDialect) CompileForeign(tableName string, c *Command) string {
	sqlStr := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s ", tableName, c.CommandIndex)
	sqlStr += fmt.Sprintf(
		"FOREIGN KEY (%s) REFERENCES `%s` (%s)",
		d.Columnize(c.CommandParameters),
		c.CommandOnTable,
		d.Columnize(c.CommandReferences),
	)
	if c.CommandOnDelete != "" {
		sqlStr += " ON DELETE " + c.CommandOnDelete
	}
	if c.CommandOnUpdate != "" {
		sqlStr += " ON UPDATE " + c.CommandOnUpdate
	}
	return sqlStr
}

func (d MySQLMigrateDialect) CompileDropForeign(tableName string, name string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY `%s`", tableName, name)
}

func (d MySQLMigrateDialect) CompileColumnType(col *ColumnDefinition) string {
	return col.buildMySQL()
}

func (d MySQLMigrateDialect) Columnize(columns []string) string {
	return "`" + strings.Join(columns, "`, `") + "`"
}

// PostgreSQLMigrateDialect implements PostgreSQL-specific DDL.
type PostgreSQLMigrateDialect struct{}

func (d PostgreSQLMigrateDialect) QuoteIdentifier(name string) string {
	return `"` + name + `"`
}

func (d PostgreSQLMigrateDialect) WrapTable(prefix, name string) string {
	return `"` + prefix + name + `"`
}

func (d PostgreSQLMigrateDialect) CompileCreateTable(temporary bool, ifNotExists bool, tableName string, columns []string) string {
	createStatement := "CREATE"
	if temporary {
		createStatement += " TEMPORARY"
	}

	ifNotExistsStr := ""
	if ifNotExists {
		ifNotExistsStr = "IF NOT EXISTS "
	}

	return fmt.Sprintf(
		"%s TABLE %s%s (%s)",
		createStatement,
		ifNotExistsStr,
		tableName,
		strings.Join(columns, ", "),
	)
}

func (d PostgreSQLMigrateDialect) CompileCreateEncoding(_ string, _ string) string {
	return ""
}

func (d PostgreSQLMigrateDialect) CompileCreateEngine(_ string) string {
	return ""
}

func (d PostgreSQLMigrateDialect) CompileAdd(tableName string, columns []string) string {
	var cols = make([]string, len(columns))
	for i, c := range columns {
		cols[i] = "ADD COLUMN " + c
	}
	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(cols, ", "))
}

func (d PostgreSQLMigrateDialect) CompileChange(tableName string, columns []string) string {
	// PostgreSQL uses ALTER COLUMN ... TYPE ... for column modifications
	// Each column change produces a separate clause
	var clauses []string
	for _, c := range columns {
		clauses = append(clauses, "ALTER COLUMN "+c)
	}
	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(clauses, ", "))
}

func (d PostgreSQLMigrateDialect) CompileDrop(tableName string) string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}

func (d PostgreSQLMigrateDialect) CompileDropIfExists(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
}

func (d PostgreSQLMigrateDialect) CompileDropColumn(tableName string, columns []string) string {
	dropColumns := make([]string, len(columns))
	for i, cc := range columns {
		dropColumns[i] = "DROP COLUMN " + cc
	}
	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(dropColumns, ", "))
}

func (d PostgreSQLMigrateDialect) CompileDropPrimary(tableName string) string {
	// PostgreSQL uses constraint name for dropping primary key
	return fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s_pkey", tableName, strings.Trim(tableName, `"`))
}

func (d PostgreSQLMigrateDialect) CompileDropIndex(_ string, indexName string) string {
	return fmt.Sprintf(`DROP INDEX "%s"`, indexName)
}

func (d PostgreSQLMigrateDialect) CompileRename(tableName string, to string) string {
	return fmt.Sprintf(`ALTER TABLE %s RENAME TO "%s"`, tableName, to)
}

func (d PostgreSQLMigrateDialect) CompileRenameColumn(tableName string, from, to string) string {
	return fmt.Sprintf(`ALTER TABLE %s RENAME COLUMN "%s" TO "%s"`, tableName, from, to)
}

func (d PostgreSQLMigrateDialect) CompileKey(tableName string, c *Command, indexType string) string {
	alg := ""
	if c.CommandAlgorithm != "" {
		alg = " USING " + c.CommandAlgorithm
	}

	switch indexType {
	case "primary key":
		return fmt.Sprintf(
			"ALTER TABLE %s ADD PRIMARY KEY (%s)",
			tableName,
			d.Columnize(c.CommandParameters),
		)
	case "unique":
		return fmt.Sprintf(
			`CREATE UNIQUE INDEX "%s" ON %s%s (%s)`,
			c.CommandIndex,
			tableName,
			alg,
			d.Columnize(c.CommandParameters),
		)
	case "spatial index":
		return fmt.Sprintf(
			`CREATE INDEX "%s" ON %s USING GIST (%s)`,
			c.CommandIndex,
			tableName,
			d.Columnize(c.CommandParameters),
		)
	default:
		return fmt.Sprintf(
			`CREATE INDEX "%s" ON %s%s (%s)`,
			c.CommandIndex,
			tableName,
			alg,
			d.Columnize(c.CommandParameters),
		)
	}
}

func (d PostgreSQLMigrateDialect) CompileForeign(tableName string, c *Command) string {
	sqlStr := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s ", tableName, c.CommandIndex)
	sqlStr += fmt.Sprintf(
		`FOREIGN KEY (%s) REFERENCES "%s" (%s)`,
		d.Columnize(c.CommandParameters),
		c.CommandOnTable,
		d.Columnize(c.CommandReferences),
	)
	if c.CommandOnDelete != "" {
		sqlStr += " ON DELETE " + c.CommandOnDelete
	}
	if c.CommandOnUpdate != "" {
		sqlStr += " ON UPDATE " + c.CommandOnUpdate
	}
	return sqlStr
}

func (d PostgreSQLMigrateDialect) CompileDropForeign(tableName string, name string) string {
	return fmt.Sprintf(`ALTER TABLE %s DROP CONSTRAINT "%s"`, tableName, name)
}

func (d PostgreSQLMigrateDialect) CompileColumnType(col *ColumnDefinition) string {
	return col.buildPostgreSQL()
}

func (d PostgreSQLMigrateDialect) Columnize(columns []string) string {
	return `"` + strings.Join(columns, `", "`) + `"`
}

var globalMigrateDialect MigrateDialect = MySQLMigrateDialect{}

// SetMigrateDialect sets the global migration dialect.
func SetMigrateDialect(d MigrateDialect) {
	globalMigrateDialect = d
}

// GetMigrateDialect returns the current global migration dialect.
func GetMigrateDialect() MigrateDialect {
	return globalMigrateDialect
}

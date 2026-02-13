package migrate

import (
	"strings"
	"testing"
)

func withPgMigrate(t *testing.T, fn func()) {
	t.Helper()
	original := GetMigrateDialect()
	SetMigrateDialect(PostgreSQLMigrateDialect{})
	defer SetMigrateDialect(original)
	fn()
}

func TestMySQLMigrateDialect_QuoteIdentifier(t *testing.T) {
	d := MySQLMigrateDialect{}
	if d.QuoteIdentifier("name") != "`name`" {
		t.Error("MySQL QuoteIdentifier failed")
	}
}

func TestPostgreSQLMigrateDialect_QuoteIdentifier(t *testing.T) {
	d := PostgreSQLMigrateDialect{}
	if d.QuoteIdentifier("name") != `"name"` {
		t.Error("PostgreSQL QuoteIdentifier failed")
	}
}

func TestMySQLMigrateDialect_WrapTable(t *testing.T) {
	d := MySQLMigrateDialect{}
	if d.WrapTable("wz_", "users") != "`wz_users`" {
		t.Error("MySQL WrapTable failed")
	}
}

func TestPostgreSQLMigrateDialect_WrapTable(t *testing.T) {
	d := PostgreSQLMigrateDialect{}
	if d.WrapTable("wz_", "users") != `"wz_users"` {
		t.Error("PostgreSQL WrapTable failed")
	}
}

func TestMySQLMigrateDialect_Columnize(t *testing.T) {
	d := MySQLMigrateDialect{}
	result := d.Columnize([]string{"name", "email"})
	if result != "`name`, `email`" {
		t.Errorf("MySQL Columnize failed: %s", result)
	}
}

func TestPostgreSQLMigrateDialect_Columnize(t *testing.T) {
	d := PostgreSQLMigrateDialect{}
	result := d.Columnize([]string{"name", "email"})
	if result != `"name", "email"` {
		t.Errorf("PostgreSQL Columnize failed: %s", result)
	}
}

func TestMySQL_CreateTable(t *testing.T) {
	SetMigrateDialect(MySQLMigrateDialect{})

	builder := NewBuilder("users", "").DefaultStringLength(255)
	builder.Engine("InnoDB")
	builder.Charset("utf8mb4")
	builder.Collation("utf8mb4_unicode_ci")
	builder.Increments("id")
	builder.String("name", 100)
	builder.Create()

	sqls := builder.Build()
	if len(sqls) != 1 {
		t.Fatalf("expected 1 SQL, got %d", len(sqls))
	}

	sql := sqls[0]
	if !strings.Contains(sql, "CREATE TABLE `users`") {
		t.Errorf("expected CREATE TABLE `users`, got: %s", sql)
	}
	if !strings.Contains(sql, "AUTO_INCREMENT PRIMARY KEY") {
		t.Errorf("expected AUTO_INCREMENT PRIMARY KEY, got: %s", sql)
	}
	if !strings.Contains(sql, "ENGINE = InnoDB") {
		t.Errorf("expected ENGINE = InnoDB, got: %s", sql)
	}
	if !strings.Contains(sql, "DEFAULT CHARACTER SET utf8mb4") {
		t.Errorf("expected DEFAULT CHARACTER SET, got: %s", sql)
	}
}

func TestPostgreSQL_CreateTable(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "").DefaultStringLength(255)
		builder.Increments("id")
		builder.String("name", 100)
		builder.Create()

		sqls := builder.Build()
		if len(sqls) != 1 {
			t.Fatalf("expected 1 SQL, got %d", len(sqls))
		}

		sql := sqls[0]
		if !strings.Contains(sql, `CREATE TABLE "users"`) {
			t.Errorf("expected CREATE TABLE with double quotes, got: %s", sql)
		}
		if !strings.Contains(sql, "SERIAL PRIMARY KEY") {
			t.Errorf("expected SERIAL PRIMARY KEY, got: %s", sql)
		}
		if strings.Contains(sql, "ENGINE") {
			t.Errorf("PostgreSQL should not have ENGINE clause, got: %s", sql)
		}
		if strings.Contains(sql, "CHARACTER SET") {
			t.Errorf("PostgreSQL should not have CHARACTER SET clause, got: %s", sql)
		}
		if strings.Contains(sql, "AUTO_INCREMENT") {
			t.Errorf("PostgreSQL should not have AUTO_INCREMENT, got: %s", sql)
		}
	})
}

func TestPostgreSQL_CreateTableIfNotExists(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.Increments("id")
		builder.CreateIfNotExists()

		sqls := builder.Build()
		sql := sqls[0]
		if !strings.Contains(sql, "IF NOT EXISTS") {
			t.Errorf("expected IF NOT EXISTS, got: %s", sql)
		}
	})
}

func TestPostgreSQL_AddColumn(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "").DefaultStringLength(255)
		builder.String("email", 255).Nullable(true)

		sqls := builder.Build()
		if len(sqls) != 1 {
			t.Fatalf("expected 1 SQL, got %d", len(sqls))
		}

		sql := sqls[0]
		if !strings.Contains(sql, `ALTER TABLE "users"`) {
			t.Errorf("expected ALTER TABLE with double quotes, got: %s", sql)
		}
		if !strings.Contains(sql, "ADD COLUMN") {
			t.Errorf("expected ADD COLUMN, got: %s", sql)
		}
		if !strings.Contains(sql, `"email"`) {
			t.Errorf("expected double-quoted column name, got: %s", sql)
		}
	})
}

func TestPostgreSQL_DropTable(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.Drop()

		sqls := builder.Build()
		if len(sqls) != 1 {
			t.Fatalf("expected 1 SQL, got %d", len(sqls))
		}
		if sqls[0] != `DROP TABLE "users"` {
			t.Errorf("expected DROP TABLE with double quotes, got: %s", sqls[0])
		}
	})
}

func TestPostgreSQL_DropIfExists(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.DropIfExists()

		sqls := builder.Build()
		if sqls[0] != `DROP TABLE IF EXISTS "users"` {
			t.Errorf("expected DROP TABLE IF EXISTS, got: %s", sqls[0])
		}
	})
}

func TestPostgreSQL_DropColumn(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.DropColumn("email", "phone")

		sqls := builder.Build()
		sql := sqls[0]
		if !strings.Contains(sql, "DROP COLUMN email") {
			t.Errorf("expected DROP COLUMN, got: %s", sql)
		}
	})
}

func TestPostgreSQL_Index(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.Increments("id")
		builder.String("name", 100)
		builder.String("email", 255)
		builder.Index("", "name", "email")
		builder.Create()

		sqls := builder.Build()
		found := false
		for _, sql := range sqls {
			if strings.Contains(sql, "CREATE INDEX") && strings.Contains(sql, `ON "users"`) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected CREATE INDEX ON statement, got: %v", sqls)
		}
	})
}

func TestPostgreSQL_UniqueIndex(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.Increments("id")
		builder.String("email", 255)
		builder.Unique("", "email")
		builder.Create()

		sqls := builder.Build()
		found := false
		for _, sql := range sqls {
			if strings.Contains(sql, "CREATE UNIQUE INDEX") && strings.Contains(sql, `ON "users"`) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected CREATE UNIQUE INDEX statement, got: %v", sqls)
		}
	})
}

func TestPostgreSQL_DropIndex(t *testing.T) {
	d := PostgreSQLMigrateDialect{}
	sql := d.CompileDropIndex(`"users"`, "idx_name")
	if sql != `DROP INDEX "idx_name"` {
		t.Errorf("expected DROP INDEX, got: %s", sql)
	}
}

func TestPostgreSQL_Rename(t *testing.T) {
	d := PostgreSQLMigrateDialect{}
	sql := d.CompileRename(`"users"`, "accounts")
	if sql != `ALTER TABLE "users" RENAME TO "accounts"` {
		t.Errorf("expected RENAME TO, got: %s", sql)
	}
}

func TestPostgreSQL_RenameColumn(t *testing.T) {
	d := PostgreSQLMigrateDialect{}
	sql := d.CompileRenameColumn(`"users"`, "name", "full_name")
	if sql != `ALTER TABLE "users" RENAME COLUMN "name" TO "full_name"` {
		t.Errorf("expected RENAME COLUMN, got: %s", sql)
	}
}

func TestPostgreSQL_Foreign(t *testing.T) {
	d := PostgreSQLMigrateDialect{}
	cmd := &Command{
		CommandIndex:      "fk_user_role",
		CommandParameters: []string{"role_id"},
		CommandOnTable:    "roles",
		CommandReferences: []string{"id"},
		CommandOnDelete:   "CASCADE",
		CommandOnUpdate:   "CASCADE",
	}
	sql := d.CompileForeign(`"users"`, cmd)
	if !strings.Contains(sql, `REFERENCES "roles"`) {
		t.Errorf("expected REFERENCES with double quotes, got: %s", sql)
	}
	if !strings.Contains(sql, "ON DELETE CASCADE") {
		t.Errorf("expected ON DELETE CASCADE, got: %s", sql)
	}
}

func TestPostgreSQL_DropForeign(t *testing.T) {
	d := PostgreSQLMigrateDialect{}
	sql := d.CompileDropForeign(`"users"`, "fk_user_role")
	if sql != `ALTER TABLE "users" DROP CONSTRAINT "fk_user_role"` {
		t.Errorf("expected DROP CONSTRAINT, got: %s", sql)
	}
}

func TestPostgreSQL_ColumnTypes(t *testing.T) {
	withPgMigrate(t, func() {
		tests := []struct {
			name     string
			setup    func(b *Builder)
			contains string
			excludes string
		}{
			{
				name: "Boolean",
				setup: func(b *Builder) {
					b.Boolean("is_active")
				},
				contains: "BOOLEAN",
			},
			{
				name: "BigIncrements",
				setup: func(b *Builder) {
					b.BigIncrements("id")
				},
				contains: "BIGSERIAL PRIMARY KEY",
			},
			{
				name: "Text",
				setup: func(b *Builder) {
					b.Text("description")
				},
				contains: "TEXT",
			},
			{
				name: "LongText",
				setup: func(b *Builder) {
					b.LongText("content")
				},
				contains: "TEXT",
			},
			{
				name: "Json",
				setup: func(b *Builder) {
					b.Json("data")
				},
				contains: "JSONB",
			},
			{
				name: "Binary",
				setup: func(b *Builder) {
					b.Binary("data")
				},
				contains: "BYTEA",
			},
			{
				name: "Timestamp",
				setup: func(b *Builder) {
					b.Timestamp("created_at", 0).Nullable(true)
				},
				contains: "TIMESTAMP",
			},
			{
				name: "Decimal",
				setup: func(b *Builder) {
					b.Decimal("price", 10, 2)
				},
				contains: "DECIMAL(10, 2)",
			},
			{
				name: "Integer no unsigned",
				setup: func(b *Builder) {
					b.Integer("count", false, true)
				},
				contains: "INTEGER",
				excludes: "UNSIGNED",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				builder := NewBuilder("test_table", "")
				tt.setup(builder)
				builder.Create()

				sqls := builder.Build()
				sql := sqls[0]
				if !strings.Contains(sql, tt.contains) {
					t.Errorf("expected %s in SQL, got: %s", tt.contains, sql)
				}
				if tt.excludes != "" && strings.Contains(sql, tt.excludes) {
					t.Errorf("should not contain %s in SQL, got: %s", tt.excludes, sql)
				}
			})
		}
	})
}

func TestPostgreSQL_DefaultValues(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.TinyInteger("status", false, true).
			Default(StringExpr("1")).
			Nullable(true)
		builder.Create()

		sqls := builder.Build()
		sql := sqls[0]
		if !strings.Contains(sql, "DEFAULT '1'") {
			t.Errorf("expected DEFAULT '1', got: %s", sql)
		}
	})
}

func TestPostgreSQL_SoftDeletes(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.Increments("id")
		builder.SoftDeletes("deleted_at", 0)
		builder.Create()

		sqls := builder.Build()
		sql := sqls[0]
		if !strings.Contains(sql, `"deleted_at"`) {
			t.Errorf("expected quoted deleted_at column, got: %s", sql)
		}
		if !strings.Contains(sql, "NULL") {
			t.Errorf("expected NULL for soft delete column, got: %s", sql)
		}
	})
}

func TestPostgreSQL_Timestamps(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "")
		builder.Increments("id")
		builder.Timestamps(0)
		builder.Create()

		sqls := builder.Build()
		sql := sqls[0]
		if !strings.Contains(sql, `"created_at"`) {
			t.Errorf("expected created_at, got: %s", sql)
		}
		if !strings.Contains(sql, `"updated_at"`) {
			t.Errorf("expected updated_at, got: %s", sql)
		}
	})
}

func TestPostgreSQL_TablePrefix(t *testing.T) {
	withPgMigrate(t, func() {
		builder := NewBuilder("users", "wz_")
		builder.Increments("id")
		builder.Create()

		sqls := builder.Build()
		sql := sqls[0]
		if !strings.Contains(sql, `"wz_users"`) {
			t.Errorf("expected prefixed table name with double quotes, got: %s", sql)
		}
	})
}

func TestSetMigrateDialect(t *testing.T) {
	original := GetMigrateDialect()
	defer SetMigrateDialect(original)

	SetMigrateDialect(PostgreSQLMigrateDialect{})
	d := GetMigrateDialect()
	if _, ok := d.(PostgreSQLMigrateDialect); !ok {
		t.Error("SetMigrateDialect failed for PostgreSQL")
	}
}

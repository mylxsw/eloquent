package query

import (
	"testing"
)

func withPostgreSQL(t *testing.T, fn func()) {
	t.Helper()
	original := GetDialect()
	SetDialect(PostgreSQLDialect{})
	defer SetDialect(original)
	fn()
}

func TestPostgreSQL_ResolveQuery(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").Where("status", EQ, 1).Select("name", "id").ResolveQuery()
		expected := `SELECT users."name", users."id" FROM users WHERE users."status" = $1`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_ResolveCount(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").Where("status", EQ, 1).Select("name", "id").ResolveCount()
		expected := `SELECT COUNT(1) as count FROM users WHERE users."status" = $1`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_ResolveInsert(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, params := Builder().Table("users").ResolveInsert(KV{
			"name": "Tom",
		})
		if len(params) != 1 || params[0] != "Tom" {
			t.Error("params mismatch")
		}
		// Field name should use double quotes and placeholder should be $1
		if sqlStr != `INSERT INTO users ("name") VALUES ($1)` {
			t.Errorf("unexpected SQL: %s", sqlStr)
		}
	})
}

func TestPostgreSQL_ResolveInsertReturning(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").ResolveInsertReturning(KV{
			"name": "Tom",
		}, "id")
		expected := `INSERT INTO users ("name") VALUES ($1) RETURNING "id"`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_ResolveUpdate(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, params := Builder().Table("users").Where("id", 1).ResolveUpdate(KV{
			"name": "Jerry",
		})
		if len(params) != 2 {
			t.Errorf("expected 2 params, got %d", len(params))
		}
		// Check that identifier quoting uses double quotes and placeholders are $N
		expected := `UPDATE users SET "name" = $1 WHERE users."id" = $2`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_ResolveDelete(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").Where("id", 1).ResolveDelete()
		expected := `DELETE FROM users WHERE users."id" = $1`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_WhereIn(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").WhereIn("id", 1, 2, 3).ResolveQuery()
		expected := `SELECT * FROM users WHERE users."id" IN ( $1, $2, $3)`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_WhereBetween(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").WhereBetween("age", 15, 20).ResolveQuery()
		expected := `SELECT * FROM users WHERE users."age" BETWEEN $1 AND $2`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_WhereNull(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").WhereNull("deleted_at").ResolveQuery()
		expected := `SELECT * FROM users WHERE users."deleted_at" IS NULL`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_WhereGroup(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").Where("status", 1).WhereGroup(func(builder Condition) {
			builder.Where("phone", "123").OrWhere("email", "x@x.com")
		}).ResolveQuery()
		expected := `SELECT * FROM users WHERE users."status" = $1 AND (  users."phone" = $2 OR users."email" = $3)`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_OrderByLimitOffset(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").OrderBy("id", "DESC").Limit(10).Offset(20).ResolveQuery()
		expected := `SELECT * FROM users ORDER BY users."id" DESC LIMIT 10 OFFSET 20`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_Merge(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("test").Where("username", "zhangsan").Merge(Builder().Where("password", "123456")).ResolveQuery()
		expected := `SELECT * FROM test WHERE test."username" = $1 AND test."password" = $2`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_AggregateQueries(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().Table("users").Where("status", 1).ResolveAvg("age")
		expected := `SELECT AVG(age) as avg FROM users WHERE users."status" = $1`
		if sqlStr != expected {
			t.Errorf("Avg expected:\n%s\ngot:\n%s", expected, sqlStr)
		}

		sqlStr, _ = Builder().Table("users").Where("status", GT, 1).ResolveSum("age")
		expected = `SELECT SUM(age) as sum FROM users WHERE users."status" > $1`
		if sqlStr != expected {
			t.Errorf("Sum expected:\n%s\ngot:\n%s", expected, sqlStr)
		}

		sqlStr, _ = Builder().Table("users").Where("id", 1).ResolveMax("age")
		expected = `SELECT MAX(age) as max FROM users WHERE users."id" = $1`
		if sqlStr != expected {
			t.Errorf("Max expected:\n%s\ngot:\n%s", expected, sqlStr)
		}

		sqlStr, _ = Builder().Table("users").Where("id", 1).ResolveMin("age")
		expected = `SELECT MIN(age) as min FROM users WHERE users."id" = $1`
		if sqlStr != expected {
			t.Errorf("Min expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

func TestPostgreSQL_LeftJoin(t *testing.T) {
	withPostgreSQL(t, func() {
		sqlStr, _ := Builder().
			Table("users AS u").
			Select("u.name", "r.name AS role_name").
			LeftJoin("roles AS r", func(builder Condition) {
				builder.WhereColumn("u.role_id", "=", "r.id")
			}).
			ResolveQuery()
		expected := `SELECT u."name", r."name" AS role_name FROM users AS u LEFT JOIN roles AS r ON   u."role_id" = r."id"`
		if sqlStr != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, sqlStr)
		}
	})
}

// Ensure MySQL dialect still works correctly after PostgreSQL tests
func TestMySQL_StillWorksAfterDialectSwitch(t *testing.T) {
	SetDialect(MySQLDialect{})

	sqlStr, _ := Builder().Table("users").Where("status", EQ, 1).Select("name", "id").ResolveCount()
	if sqlStr != "SELECT COUNT(1) as count FROM users WHERE users.`status` = ?" {
		t.Errorf("MySQL test failed after dialect switch: %s", sqlStr)
	}
}

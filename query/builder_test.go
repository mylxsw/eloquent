package query

import (
	"testing"
)

func TestSQLBuilder_ResolveCount(t *testing.T) {
	builder := Builder().Table("users").Where("status", "=", 1).Select("name", "id")

	sqlStr, _ := builder.ResolveCount()
	if sqlStr != "SELECT COUNT(1) as count FROM users WHERE users.`status` = ?" {
		t.Error("test failed")
	}

	sqlStr, _ = builder.ResolveQuery()
	if sqlStr != "SELECT users.`name`, users.`id` FROM users WHERE users.`status` = ?" {
		t.Error("test failed")
	}
}

func TestSQLBuilder_Where(t *testing.T) {
	sqlStr, _ := Builder().Table("users").Where("status", 1).ResolveAvg("age")
	if sqlStr != "SELECT AVG(age) as avg FROM users WHERE users.`status` = ?" {
		t.Error("test failed")
	}

	sqlStr, _ = Builder().Table("users").Where("status", ">", 1).ResolveSum("age")
	if sqlStr != "SELECT SUM(age) as sum FROM users WHERE users.`status` > ?" {
		t.Error("test failed")
	}
}

func TestSQLBuilder_WhereBetween(t *testing.T) {
	sqlStr, _ := Builder().Table("users").WhereBetween("age", 15, 20).OrWhereBetween("age", 30, 45).ResolveQuery()
	if sqlStr != "SELECT * FROM users WHERE users.`age` BETWEEN ? AND ? OR users.`age` BETWEEN ? AND ?" {
		t.Error("test failed")
	}
}
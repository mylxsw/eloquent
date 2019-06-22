package migrate

import (
	"fmt"
	"strings"

	"github.com/mylxsw/eloquent/migrate/schema"
)

func StringExpr(value string) schema.Expr {
	return schema.Expr{
		Type:  schema.ExprTypeString,
		Value: value,
	}
}

func RawExpr(value string) schema.Expr {
	return schema.Expr{
		Type:  schema.ExprTypeRaw,
		Value: value,
	}
}

type Column struct {
	ColumnName          string
	ColumnType          string
	ColumnComment       string
	ColumnAutoIncrement bool
	ColumnUnsigned      bool
	ColumnNullable      bool
	ColumnDefault       schema.Expr
	ColumnCharset       string
	ColumnCollation     string
	ColumnUseCurrent    bool
	ColumnVirtualAs     string
	ColumnStoredAs      string
	ColumnAfter         string
	ColumnFirst         bool
	ColumnSrid          int64

	ColumnChange bool
}

func (c *Column) Build() string {
	sqlStr := "`" + c.ColumnName + "`"
	sqlStr += strings.ToUpper(c.Type())

	sqlStr += c.modifyUnsigned()
	sqlStr += c.modifyVirtualAs()
	sqlStr += c.modifyStoredAs()
	sqlStr += c.modifyCharset()
	sqlStr += c.modifyCollate()
	sqlStr += c.modifyNullable()
	sqlStr += c.modifyDefault()
	sqlStr += c.modifyIncrement()
	sqlStr += c.modifyComment()
	sqlStr += c.modifyAfter()
	sqlStr += c.modifyFirst()
	sqlStr += c.modifySrid()

	return sqlStr
}

func (c *Column) Change() schema.ColumnType {
	c.ColumnChange = true
	return c
}

func (c *Column) IsChange() bool {
	return c.ColumnChange
}

func (c *Column) Nullable(value bool) schema.ColumnType {
	c.ColumnNullable = value
	return c
}

func (c *Column) After(name string) schema.ColumnType {
	c.ColumnAfter = name
	return c
}

func (c *Column) AutoIncrement() schema.ColumnType {
	c.ColumnAutoIncrement = true
	return c
}

func (c *Column) Charset(charset string) schema.ColumnType {
	c.ColumnCharset = charset
	return c
}

func (c *Column) Collation(collation string) schema.ColumnType {
	c.ColumnCollation = collation
	return c
}

func (c *Column) Comment(comment string) schema.ColumnType {
	c.ColumnComment = comment
	return c
}

func (c *Column) Default(defaultVal schema.Expr) schema.ColumnType {
	c.ColumnDefault = defaultVal
	return c
}

func (c *Column) First() schema.ColumnType {
	c.ColumnFirst = true
	return c
}

func (c *Column) StoredAs(expression string) schema.ColumnType {
	c.ColumnStoredAs = expression
	return c
}

func (c *Column) Unsigned() schema.ColumnType {
	c.ColumnUnsigned = true
	return c
}

func (c *Column) UseCurrent() schema.ColumnType {
	c.ColumnUseCurrent = true
	return c
}

func (c *Column) VirtualAs(expression string) schema.ColumnType {
	c.ColumnVirtualAs = expression
	return c
}

func (c *Column) GeneratedAs(expression string) schema.ColumnType {
	return c
}

func (c *Column) Always() schema.ColumnType {
	return c
}

func (c *Column) Type() string {
	if c.ColumnUseCurrent {
		return " " + c.ColumnType + " DEFAULT CURRENT_TIMESTAMP"
	}

	return " " + c.ColumnType
}

func (c *Column) modifyUnsigned() string {
	if c.ColumnUnsigned {
		return " UNSIGNED"
	}

	return ""
}

func (c *Column) modifyVirtualAs() string {
	if c.ColumnVirtualAs != "" {
		return fmt.Sprintf(" AS (%s)", c.ColumnVirtualAs)
	}
	return ""
}

func (c *Column) modifyStoredAs() string {
	if c.ColumnStoredAs != "" {
		return fmt.Sprintf(" AS (%s) stored", c.ColumnStoredAs)
	}
	return ""
}

func (c *Column) modifyCharset() string {
	if c.ColumnCharset != "" {
		return " CHARACTER SET " + c.ColumnCharset
	}
	return ""
}

func (c *Column) modifyCollate() string {
	if c.ColumnCollation != "" {
		return fmt.Sprintf(" COLLATE '%s'", c.ColumnCollation)
	}

	return ""
}

func (c *Column) modifyNullable() string {
	if c.ColumnVirtualAs == "" && c.ColumnStoredAs == "" {
		if c.ColumnNullable {
			return " NULL"
		}

		return " NOT NULL"
	}

	return ""
}

func (c *Column) modifyDefault() string {
	if c.ColumnDefault.Value != "" {
		var value string

		if c.ColumnDefault.Type == schema.ExprTypeRaw {
			value = c.ColumnDefault.Value
		} else {
			value = "'" + c.ColumnDefault.Value + "'"
		}

		return " DEFAULT " + value
	}

	return ""
}

func (c *Column) modifyIncrement() string {
	if c.ColumnAutoIncrement {
		return " AUTO_INCREMENT PRIMARY KEY"
	}

	return ""
}

func (c *Column) modifyComment() string {
	if c.ColumnComment != "" {
		return " COMMENT '" + addSlashes(c.ColumnComment) + "'"
	}

	return ""
}

func (c *Column) modifyAfter() string {
	if c.ColumnAfter != "" {
		return " AFTER " + c.ColumnAfter
	}

	return ""
}
func (c *Column) modifyFirst() string {
	if c.ColumnFirst {
		return " FIRST"
	}

	return ""
}

func (c *Column) modifySrid() string {
	if c.ColumnSrid > 0 {
		return fmt.Sprintf(" SRID %d", c.ColumnSrid)
	}
	return ""
}

func addSlashes(str string) string {
	var tmpRune []rune
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

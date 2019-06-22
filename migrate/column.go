package migrate

import (
	"fmt"
	"strings"
)

func StringExpr(value string) Expr {
	return Expr{
		Type:  ExprTypeString,
		Value: value,
	}
}

func RawExpr(value string) Expr {
	return Expr{
		Type:  ExprTypeRaw,
		Value: value,
	}
}

type ColumnDefinition struct {
	ColumnName          string
	ColumnType          string
	ColumnComment       string
	ColumnAutoIncrement bool
	ColumnUnsigned      bool
	ColumnNullable      bool
	ColumnDefault       Expr
	ColumnCharset       string
	ColumnCollation     string
	ColumnUseCurrent    bool
	ColumnVirtualAs     string
	ColumnStoredAs      string
	ColumnAfter         string
	ColumnFirst         bool
	ColumnSrid          int64

	ColumnIndex        string
	ColumnPrimary      bool
	ColumnUnique       bool
	ColumnSpatialIndex bool

	ColumnChange bool
}

func (c *ColumnDefinition) Index(name string) *ColumnDefinition {
	c.ColumnIndex = name
	return c
}

func (c *ColumnDefinition) Primary() *ColumnDefinition {
	c.ColumnPrimary = true
	return c
}

func (c *ColumnDefinition) Unique() *ColumnDefinition {
	c.ColumnUnique = true
	return c
}

func (c *ColumnDefinition) SpatialIndex() *ColumnDefinition {
	c.ColumnSpatialIndex = true
	return c
}

func (c *ColumnDefinition) Build() string {
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

func (c *ColumnDefinition) Change() *ColumnDefinition {
	c.ColumnChange = true
	return c
}

func (c *ColumnDefinition) IsChange() bool {
	return c.ColumnChange
}

func (c *ColumnDefinition) Nullable(value bool) *ColumnDefinition {
	c.ColumnNullable = value
	return c
}

func (c *ColumnDefinition) After(name string) *ColumnDefinition {
	c.ColumnAfter = name
	return c
}

func (c *ColumnDefinition) AutoIncrement() *ColumnDefinition {
	c.ColumnAutoIncrement = true
	return c
}

func (c *ColumnDefinition) Charset(charset string) *ColumnDefinition {
	c.ColumnCharset = charset
	return c
}

func (c *ColumnDefinition) Collation(collation string) *ColumnDefinition {
	c.ColumnCollation = collation
	return c
}

func (c *ColumnDefinition) Comment(comment string) *ColumnDefinition {
	c.ColumnComment = comment
	return c
}

func (c *ColumnDefinition) Default(defaultVal Expr) *ColumnDefinition {
	c.ColumnDefault = defaultVal
	return c
}

func (c *ColumnDefinition) First() *ColumnDefinition {
	c.ColumnFirst = true
	return c
}

func (c *ColumnDefinition) StoredAs(expression string) *ColumnDefinition {
	c.ColumnStoredAs = expression
	return c
}

func (c *ColumnDefinition) Unsigned() *ColumnDefinition {
	c.ColumnUnsigned = true
	return c
}

func (c *ColumnDefinition) UseCurrent() *ColumnDefinition {
	c.ColumnUseCurrent = true
	return c
}

func (c *ColumnDefinition) VirtualAs(expression string) *ColumnDefinition {
	c.ColumnVirtualAs = expression
	return c
}

func (c *ColumnDefinition) GeneratedAs(expression string) *ColumnDefinition {
	return c
}

func (c *ColumnDefinition) Always() *ColumnDefinition {
	return c
}

func (c *ColumnDefinition) Type() string {
	if c.ColumnUseCurrent {
		return " " + c.ColumnType + " DEFAULT CURRENT_TIMESTAMP"
	}

	return " " + c.ColumnType
}

func (c *ColumnDefinition) modifyUnsigned() string {
	if c.ColumnUnsigned {
		return " UNSIGNED"
	}

	return ""
}

func (c *ColumnDefinition) modifyVirtualAs() string {
	if c.ColumnVirtualAs != "" {
		return fmt.Sprintf(" AS (%s)", c.ColumnVirtualAs)
	}
	return ""
}

func (c *ColumnDefinition) modifyStoredAs() string {
	if c.ColumnStoredAs != "" {
		return fmt.Sprintf(" AS (%s) stored", c.ColumnStoredAs)
	}
	return ""
}

func (c *ColumnDefinition) modifyCharset() string {
	if c.ColumnCharset != "" {
		return " CHARACTER SET " + c.ColumnCharset
	}
	return ""
}

func (c *ColumnDefinition) modifyCollate() string {
	if c.ColumnCollation != "" {
		return fmt.Sprintf(" COLLATE '%s'", c.ColumnCollation)
	}

	return ""
}

func (c *ColumnDefinition) modifyNullable() string {
	if c.ColumnVirtualAs == "" && c.ColumnStoredAs == "" {
		if c.ColumnNullable {
			return " NULL"
		}

		return " NOT NULL"
	}

	return ""
}

func (c *ColumnDefinition) modifyDefault() string {
	if c.ColumnDefault.Value != "" {
		var value string

		if c.ColumnDefault.Type == ExprTypeRaw {
			value = c.ColumnDefault.Value
		} else {
			value = "'" + c.ColumnDefault.Value + "'"
		}

		return " DEFAULT " + value
	}

	return ""
}

func (c *ColumnDefinition) modifyIncrement() string {
	if c.ColumnAutoIncrement {
		return " AUTO_INCREMENT PRIMARY KEY"
	}

	return ""
}

func (c *ColumnDefinition) modifyComment() string {
	if c.ColumnComment != "" {
		return " COMMENT '" + addSlashes(c.ColumnComment) + "'"
	}

	return ""
}

func (c *ColumnDefinition) modifyAfter() string {
	if c.ColumnAfter != "" {
		return " AFTER " + c.ColumnAfter
	}

	return ""
}
func (c *ColumnDefinition) modifyFirst() string {
	if c.ColumnFirst {
		return " FIRST"
	}

	return ""
}

func (c *ColumnDefinition) modifySrid() string {
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

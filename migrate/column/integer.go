package column

import (
	"fmt"
	"strconv"

	"github.com/mylxsw/eloquent/migrate/schema"
)

type IntegerColumn struct {
	ColumnName          string
	ColumnIntType       string
	ColumnAutoIncrement bool
	ColumnUnsigned      bool
	ColumnComment       string
	ColumnNullable      bool
	ColumnDefault       int
}

func (c *IntegerColumn) Nullable(value bool) schema.ColumnType {
	c.ColumnNullable = value
	return c
}

func (c *IntegerColumn) After(name string) schema.ColumnType {
	return c
}

func (c *IntegerColumn) AutoIncrement() schema.ColumnType {
	c.ColumnAutoIncrement = true
	return c
}

func (c *IntegerColumn) Charset(charset string) schema.ColumnType {
	return c
}

func (c *IntegerColumn) Collation(collation string) schema.ColumnType {
	return c
}

func (c *IntegerColumn) Comment(comment string) schema.ColumnType {
	c.ColumnComment = comment
	return c
}

func (c *IntegerColumn) Default(defaultVal string) schema.ColumnType {
	val, err := strconv.Atoi(defaultVal)
	if err != nil {
		panic(err)
	}

	c.ColumnDefault = val
	return c
}

func (c *IntegerColumn) First() schema.ColumnType {
	return c
}

func (c *IntegerColumn) StoredAs(expression string) schema.ColumnType {
	return c
}

func (c *IntegerColumn) Unsigned() schema.ColumnType {
	c.ColumnUnsigned = true
	return c
}

func (c *IntegerColumn) UseCurrent() schema.ColumnType {
	return c
}

func (c *IntegerColumn) VirtualAs(expression string) schema.ColumnType {
	return c
}

func (c *IntegerColumn) GeneratedAs(expression string) schema.ColumnType {
	return c
}

func (c *IntegerColumn) Always() schema.ColumnType {
	return c
}

func (c *IntegerColumn) Build() string {
	return fmt.Sprintf("%s %s COMMENT %s", c.ColumnName, c.ColumnIntType, c.ColumnComment)
}

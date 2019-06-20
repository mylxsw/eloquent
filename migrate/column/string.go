package column

import (
	"fmt"

	"github.com/mylxsw/eloquent/migrate/schema"
)

type StringColumn struct {
	ColumnName    string
	ColumnLength  int
	ColumnComment string
}

func (c *StringColumn) Nullable(value bool) schema.ColumnType {
	return c
}

func (c *StringColumn) After(name string) schema.ColumnType {
	return c
}

func (c *StringColumn) AutoIncrement() schema.ColumnType {
	return c
}

func (c *StringColumn) Charset(charset string) schema.ColumnType {
	return c
}

func (c *StringColumn) Collation(collation string) schema.ColumnType {
	return c
}

func (c *StringColumn) Comment(comment string) schema.ColumnType {
	c.ColumnComment = comment
	return c
}

func (c *StringColumn) Default(defaultVal string) schema.ColumnType {
	return c
}

func (c *StringColumn) First() schema.ColumnType {
	return c
}

func (c *StringColumn) StoredAs(expression string) schema.ColumnType {
	return c
}

func (c *StringColumn) Unsigned() schema.ColumnType {
	return c
}

func (c *StringColumn) UseCurrent() schema.ColumnType {
	return c
}

func (c *StringColumn) VirtualAs(expression string) schema.ColumnType {
	return c
}

func (c *StringColumn) GeneratedAs(expression string) schema.ColumnType {
	return c
}

func (c *StringColumn) Always() schema.ColumnType {
	return c
}

func (c *StringColumn) Build() string {
	return fmt.Sprintf("%s VARCHAR(%d) COMMENT %s", c.ColumnName, c.ColumnLength, c.ColumnComment)
}
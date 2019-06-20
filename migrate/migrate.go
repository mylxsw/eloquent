package migrate

import (
	"fmt"

	"github.com/mylxsw/eloquent/migrate/schema"
)

type Schema struct {
	tableCreate        string
	tableCreateBuilder schema.TableBuilder

	tableDrop string

	tableUpdate        string
	tableUpdateBuilder schema.TableBuilder
}

// Create creat a new table
func (s *Schema) Create(table string, apply func(builder schema.TableBuilder)) {
	builder := tableBuilder{}
	apply(&builder)

	s.tableCreate = table
	s.tableCreateBuilder = &builder

	fmt.Println(s.tableCreate, s.tableCreateBuilder.Build())
}

// Drop drop a existing table
func (s *Schema) Drop(table string) {
	s.tableDrop = table
}

// Table update a existing table
func (s *Schema) Table(table string, apply func(builder schema.TableBuilder)) {
	builder := tableBuilder{}
	apply(&builder)

	s.tableUpdate = table
	s.tableUpdateBuilder = &builder
}

// NewSchema create a new Schema
func NewSchema() schema.Schema {
	return &Schema{}
}

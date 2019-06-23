package migrate

type Schema struct {
	m       *Manager
	version string

	tableName string
}

// Create creat a new table
func (s *Schema) Create(table string, apply func(builder *Builder)) {
	builder := NewBuilder(table, s.m.Prefix).DefaultStringLength(s.m.DefaultStringLength)
	builder.Engine(s.m.Engine)
	builder.Charset(s.m.Charset)
	builder.Collation(s.m.Collation)

	builder.Create()
	apply(builder)

	s.tableName = table
	s.m.Execute(builder, s.version)
}

// Table update a existing table
func (s *Schema) Table(table string, apply func(builder *Builder)) {
	builder := NewBuilder(table, s.m.Prefix).DefaultStringLength(s.m.DefaultStringLength)
	builder.Engine(s.m.Engine)
	builder.Charset(s.m.Charset)
	builder.Collation(s.m.Collation)

	apply(builder)

	s.tableName = table
	s.m.Execute(builder, s.version)
}

// NewSchema create a new Schema
func NewSchema(m *Manager, v string) *Schema {
	return &Schema{version: v, m: m}
}

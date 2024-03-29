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

// Table update an existing table
func (s *Schema) Table(table string, apply func(builder *Builder)) {
	builder := NewBuilder(table, s.m.Prefix).DefaultStringLength(s.m.DefaultStringLength)
	builder.Engine(s.m.Engine)
	builder.Charset(s.m.Charset)
	builder.Collation(s.m.Collation)

	apply(builder)

	s.tableName = table
	s.m.Execute(builder, s.version)
}

// Drop a table
func (s *Schema) Drop(table string) {
	builder := NewBuilder(table, s.m.Prefix)
	builder.Drop()

	s.tableName = table
	s.m.Execute(builder, s.version)
}

// DropIfExists Drop a table if exists
func (s *Schema) DropIfExists(table string) {
	builder := NewBuilder(table, s.m.Prefix)
	builder.DropIfExists()

	s.tableName = table
	s.m.Execute(builder, s.version)
}

// Raw execute a raw sql statements
func (s *Schema) Raw(table string, apply func() []string) {
	s.m.ExecuteRaw(s.version, table, apply()...)
}

// NewSchema create a new Schema
func NewSchema(m *Manager, v string) *Schema {
	return &Schema{version: v, m: m}
}

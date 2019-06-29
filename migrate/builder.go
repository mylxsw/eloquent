package migrate

import (
	"fmt"
	"strings"
)

type Builder struct {
	name                string
	prefix              string
	engine              string
	charset             string
	collation           string
	temporary           bool
	defaultStringLength int
	columns             []*ColumnDefinition
	commands            []*Command
}

func NewBuilder(name string, prefix string) *Builder {
	return &Builder{
		name:     name,
		prefix:   prefix,
		columns:  make([]*ColumnDefinition, 0),
		commands: make([]*Command, 0),
	}
}

func (t *Builder) GetTableName() string {
	return t.prefix + t.name
}

func (t *Builder) DefaultStringLength(length int) *Builder {
	t.defaultStringLength = length
	return t
}

func (t *Builder) addImpliedCommands() {
	if !t.creating() {
		if len(t.getAddedColumns()) > 0 {
			t.addCommand("add")
		}

		if len(t.getChangedColumns()) > 0 {
			t.addCommand("change")
		}
	}

	// add column index
	for _, c := range t.columns {
		if c.ColumnIndex != "" {
			t.Index(c.ColumnIndex, c.ColumnName)
		}

		if c.ColumnUnique {
			t.Unique("", c.ColumnName)
		}

		if c.ColumnPrimary {
			t.Primary("", c.ColumnName)
		}

		if c.ColumnSpatialIndex {
			t.SpatialIndex("", c.ColumnName)
		}
	}
}

func (t *Builder) getAddedColumns() []*ColumnDefinition {
	var cols = make([]*ColumnDefinition, 0)
	for _, c := range t.columns {
		if !c.IsChange() {
			cols = append(cols, c)
		}
	}

	return cols
}

func (t *Builder) getChangedColumns() []*ColumnDefinition {
	var cols = make([]*ColumnDefinition, 0)
	for _, c := range t.columns {
		if c.IsChange() {
			cols = append(cols, c)
		}
	}

	return cols
}

func (t *Builder) Build() []string {
	t.addImpliedCommands()

	sqlStrs := make([]string, 0)
	for _, c := range t.commands {
		sqlStrs = append(sqlStrs, c.Build())
	}

	return sqlStrs
}

func (t *Builder) wrapTable() string {
	return "`" + t.GetTableName() + "`"
}

func (t *Builder) getColumns() []string {
	cols := make([]string, 0)
	for _, c := range t.getAddedColumns() {
		colStr := c.Build()
		if colStr != "" {
			cols = append(cols, colStr)
		}
	}

	return cols
}

func (t *Builder) Engine(engine string) {
	t.engine = engine
}

func (t *Builder) Charset(charset string) {
	t.charset = charset
}

func (t *Builder) Collation(collation string) {
	t.collation = collation
}

func (t *Builder) Temporary() {
	t.temporary = true
}

func (t *Builder) BigIncrements(name string) *ColumnDefinition {
	return t.UnsignedBigInteger(name, true)
}

func (t *Builder) BigInteger(name string, autoIncrement bool, unsigned bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "bigint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *Builder) Binary(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "blob",
	})
}

func (t *Builder) Boolean(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "tinyint(1)",
	})
}

func (t *Builder) Char(name string, length int) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: fmt.Sprintf("char(%d)", length),
	})
}

func (t *Builder) Date(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "date",
	})
}

func (t *Builder) DateTime(name string, precision int) *ColumnDefinition {
	dateType := "datetime"
	if precision > 0 {
		dateType += fmt.Sprintf("(%d)", precision)
	}

	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: dateType,
	})
}

func (t *Builder) DateTimeTz(name string, precision int) *ColumnDefinition {
	return t.DateTime(name, precision)
}

func (t *Builder) Decimal(name string, total int, scale int) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:     name,
		ColumnUnsigned: false,
		ColumnType:     fmt.Sprintf("decimal(%d, %d)", total, scale),
	})
}

func (t *Builder) Double(name string, total int, scale int) *ColumnDefinition {
	fieldType := "double"
	if total != 0 && scale != 0 {
		fieldType += fmt.Sprintf("(%d, %d)", total, scale)
	}

	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: fieldType,
	})
}

func (t *Builder) Enum(name string, items ...string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: fmt.Sprintf("enum('%s')", strings.Join(items, "', '")),
	})
}

func (t *Builder) Float(name string, total int, scale int) *ColumnDefinition {
	return t.Double(name, total, scale)
}

func (t *Builder) Geometry(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "geometry",
	})
}

func (t *Builder) GeometryCollection(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "geometrycollection",
	})
}

func (t *Builder) Increments(name string) *ColumnDefinition {
	return t.UnsignedInteger(name, true)
}

func (t *Builder) Integer(name string, autoIncrement bool, unsigned bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "int",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *Builder) IpAddress(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "varchar(45)",
	})
}

func (t *Builder) Json(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "json",
	})
}

func (t *Builder) Jsonb(name string) *ColumnDefinition {
	return t.Json(name)
}

func (t *Builder) LineString(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "linestring",
	})
}

func (t *Builder) LongText(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "longtext",
	})
}

func (t *Builder) MacAddress(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "varchar(17)",
	})
}

func (t *Builder) MediumIncrements(name string) *ColumnDefinition {
	return t.UnsignedMediumInteger(name, true)
}

func (t *Builder) MediumInteger(name string, autoIncrement bool, unsigned bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "mediumint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *Builder) MediumText(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "mediumtext",
	})
}

func (t *Builder) Morphs(name string, indexName string) {
	t.String(name+"_type", t.defaultStringLength)
	t.UnsignedBigInteger(name+"_id", false)
	t.Index(indexName, name+"_type", name+"_id")
}

func (t *Builder) NullableMorphs(name string, indexName string) {
	t.String(name+"_type", t.defaultStringLength).Nullable(true)
	t.UnsignedBigInteger(name+"_id", false).Nullable(true)
	t.Index(indexName, name+"_type", name+"_id")
}

func (t *Builder) DropMorphs(name string, indexName string) {
	t.DropIndex(indexName)
	t.DropColumn(name+"_type", name+"_id")
}

func (t *Builder) DropColumn(columns ...string) {
	t.addCommand("dropColumn", columns...)
}

func (t *Builder) MultiLineString(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "multilinestring",
	})
}

func (t *Builder) MultiPoint(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "multipoint",
	})
}

func (t *Builder) MultiPolygon(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "multipolygon",
	})
}

func (t *Builder) NullableTimestamps(precision int) {
	t.Timestamps(precision)
}

func (t *Builder) Point(name string, srid int) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "point",
	})
}

func (t *Builder) Polygon(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "polygon",
	})
}

func (t *Builder) RememberToken() *ColumnDefinition {
	return t.String("remember_token", 100).Nullable(true)
}

func (t *Builder) DropRememberToken() {
	t.DropColumn("remember_token")
}

func (t *Builder) SmallIncrements(name string) *ColumnDefinition {
	return t.UnsignedSmallInteger(name, true)
}

func (t *Builder) SmallInteger(name string, autoIncrement bool, unsigned bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "smallint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *Builder) SoftDeletes(column string, precision int) *ColumnDefinition {
	return t.Timestamp(column, precision).Nullable(true)
}

func (t *Builder) SoftDeletesTz(column string, precision int) *ColumnDefinition {
	return t.TimestampTz(column, precision).Nullable(true)
}

func (t *Builder) String(name string, length int) *ColumnDefinition {
	if length <= 0 {
		length = t.defaultStringLength
	}

	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: fmt.Sprintf("VARCHAR(%d)", length),
	})
}

func (t *Builder) Text(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "text",
	})
}

func (t *Builder) Time(name string, precision int) *ColumnDefinition {
	dateType := "time"
	if precision > 0 {
		dateType += fmt.Sprintf("(%d)", precision)
	}

	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: dateType,
	})
}

func (t *Builder) TimeTz(name string, precision int) *ColumnDefinition {
	return t.Time(name, precision)
}

func (t *Builder) Timestamp(name string, precision int) *ColumnDefinition {
	dateType := "timestamp"
	if precision > 0 {
		dateType += fmt.Sprintf("(%d)", precision)
	}

	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: dateType,
	})
}

func (t *Builder) TimestampTz(name string, precision int) *ColumnDefinition {
	return t.Timestamp(name, precision)
}

func (t *Builder) Timestamps(precision int) {
	t.Timestamp("created_at", precision).Nullable(true)
	t.Timestamp("updated_at", precision).Nullable(true)
}

func (t *Builder) TimestampsTz(precision int) {
	t.TimestampTz("created_at", precision).Nullable(true)
	t.TimestampTz("updated_at", precision).Nullable(true)
}

func (t *Builder) TinyIncrements(name string) *ColumnDefinition {
	return t.UnsignedTinyInteger(name, true)
}

func (t *Builder) TinyInteger(name string, autoIncrement bool, unsigned bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "tinyint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *Builder) UnsignedBigInteger(name string, autoIncrement bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "bigint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *Builder) UnsignedDecimal(name string, total int, scale int) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:     name,
		ColumnUnsigned: true,
		ColumnType:     fmt.Sprintf("decimal(%d, %d)", total, scale),
	})
}

func (t *Builder) UnsignedInteger(name string, autoIncrement bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "int",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *Builder) UnsignedMediumInteger(name string, autoIncrement bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "mediumint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *Builder) UnsignedSmallInteger(name string, autoIncrement bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "smallint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *Builder) UnsignedTinyInteger(name string, autoIncrement bool) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName:          name,
		ColumnType:          "tinyint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *Builder) Uuid(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "char(36)",
	})
}

func (t *Builder) Year(name string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: "year",
	})
}

func (t *Builder) Set(name string, items ...string) *ColumnDefinition {
	return t.addColumn(&ColumnDefinition{
		ColumnName: name,
		ColumnType: fmt.Sprintf("set('%s')", strings.Join(items, "', '")),
	})
}

func (t *Builder) addColumn(c *ColumnDefinition) *ColumnDefinition {
	t.columns = append(t.columns, c)

	return c
}

func (t *Builder) Unique(name string, columns ...string) *Command {
	return t.indexCommand("unique", name, columns...)
}

func (t *Builder) DropUnique(name string) *Command {
	return t.indexCommand("dropUnique", name)
}

func (t *Builder) Index(name string, columns ...string) *Command {
	return t.indexCommand("index", name, columns...)
}

func (t *Builder) DropIndex(name string) *Command {
	return t.indexCommand("dropIndex", name)
}

func (t *Builder) Primary(name string, columns ...string) *Command {
	return t.indexCommand("primary", name, columns...)
}

func (t *Builder) DropPrimary(name string) *Command {
	return t.indexCommand("dropPrimary", name)
}

func (t *Builder) SpatialIndex(name string, columns ...string) *Command {
	return t.indexCommand("spatialIndex", name, columns...)
}

func (t *Builder) DropSpatialIndex(name string) *Command {
	return t.indexCommand("dropSpatialIndex", name)
}

func (t *Builder) Drop() *Command {
	return t.addCommand("drop")
}

func (t *Builder) Create() *Command {
	return t.addCommand("create")
}

func (t *Builder) CreateIfNotExists() *Command {
	return t.addCommand("createIfNotExists")
}

func (t *Builder) DropIfExists() *Command {
	return t.addCommand("dropIfExists")
}

func (t *Builder) Rename(to string) *Command {
	return t.addCommand("rename", to)
}

// RenameColumn rename a column name (only support for MySQL 8.0)
func (t *Builder) RenameColumn(from string, to string) *Command {
	return t.addCommand("renameColumn", from, to)
}

func (t *Builder) Foreign(name string, columns ...string) *Command {
	return t.indexCommand("foreign", name, columns...)
}

func (t *Builder) DropForeign(name string) *Command {
	return t.addCommand("dropForeign", name)
}

func (t *Builder) indexCommand(indexType string, indexName string, columns ...string) *Command {
	if indexName == "" {
		indexName = createIndexName(t.GetTableName(), indexType, columns...)
	}

	return t.addCommand(indexType, columns...).Index(indexName)
}

func (t *Builder) addCommand(name string, parameters ...string) *Command {
	cmd := NewCommand(t)
	cmd.Name(name).Columns(parameters...)

	t.commands = append(t.commands, cmd)

	return cmd
}

func (t *Builder) creating() bool {
	for _, c := range t.commands {
		if c.Equal("create") || c.Equal("createIfNotExists") {
			return true
		}
	}

	return false
}

func (t *Builder) compileKey(c *Command, indexType string) string {
	alg := ""
	if c.CommandAlgorithm != "" {
		alg = " USING " + c.CommandAlgorithm
	}

	return fmt.Sprintf(
		"ALTER TABLE %s ADD %s %s%s(%s)",
		t.wrapTable(),
		indexType,
		c.CommandIndex,
		alg,
		columnize(c.CommandParameters),
	)
}

func (t *Builder) compileDrop(c *Command) string {
	return fmt.Sprintf("DROP TABLE %s", t.wrapTable())
}

func (t *Builder) compileDropIfExists(c *Command) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", t.wrapTable())
}

func (t *Builder) compileDropColumn(c *Command) string {
	dropColumns := make([]string, len(c.CommandParameters))
	for i, cc := range c.CommandParameters {
		dropColumns[i] = "DROP " + cc
	}

	return fmt.Sprintf("ALTER TABLE %s %s", t.wrapTable(), strings.Join(dropColumns, ", "))
}

func (t *Builder) compileDropPrimary(c *Command) string {
	return fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY", t.wrapTable())
}

func (t *Builder) compileDropIndex(c *Command) string {
	return fmt.Sprintf("ALTER TABLE %s DROP INDEX `%s`", t.wrapTable(), c.CommandIndex)
}

func (t *Builder) compileRename(c *Command) string {
	return fmt.Sprintf("RENAME TABLE %s TO `%s`", t.wrapTable(), c.CommandParameters[0])
}

func (t *Builder) compileCreateCommand(createIfNotExist bool) string {
	sqlStr := t.compileCreateTable(createIfNotExist)
	sqlStr += t.compileCreateEncoding()
	sqlStr += t.compileCreateEngine()

	return sqlStr
}

func (t *Builder) compileCreateEngine() string {
	if t.engine != "" {
		return " ENGINE = " + t.engine
	}

	return ""
}

func (t *Builder) compileCreateEncoding() string {
	sqlStr := ""
	if t.charset != "" {
		sqlStr += " DEFAULT CHARACTER SET " + t.charset
	}

	if t.collation != "" {
		sqlStr += " COLLATE " + t.collation
	}

	return sqlStr
}

func (t *Builder) compileCreateTable(createIfNotExist bool) string {
	createStatement := "CREATE"
	if t.temporary {
		createStatement += " TEMPORARY"
	}

	ifNotExists := ""
	if createIfNotExist {
		ifNotExists = "IF NOT EXISTS "
	}

	return fmt.Sprintf(
		"%s TABLE %s%s (%s)",
		createStatement,
		ifNotExists,
		t.wrapTable(),
		strings.Join(t.getColumns(), ", "),
	)
}

func (t *Builder) compileAdd() string {
	var cols = make([]string, 0)
	for _, c := range t.getColumns() {
		cols = append(cols, "ADD "+c)
	}

	return fmt.Sprintf("ALTER TABLE %s %s", t.wrapTable(), strings.Join(cols, ", "))
}

func (t *Builder) compileChange() string {
	var cols = make([]string, 0)
	for _, c := range t.getChangedColumns() {
		cols = append(cols, "MODIFY "+c.Build())
	}

	return fmt.Sprintf("ALTER TABLE %s %s", t.wrapTable(), strings.Join(cols, ", "))
}

func (t *Builder) compileRenameColumn(from string, to string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN `%s` TO `%s`", t.wrapTable(), from, to)
}

func (t *Builder) compileForeign(c *Command) string {
	sqlStr := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s ", t.wrapTable(), c.CommandIndex)
	sqlStr += fmt.Sprintf(
		"FOREIGN KEY (%s) REFERENCES `%s` (%s)",
		columnize(c.CommandParameters),
		c.CommandOnTable,
		columnize(c.CommandReferences),
	)

	if c.CommandOnDelete != "" {
		sqlStr += " ON DELETE " + c.CommandOnDelete
	}

	if c.CommandOnUpdate != "" {
		sqlStr += " ON UPDATE " + c.CommandOnUpdate
	}

	return sqlStr
}

func (t *Builder) compileDropForeign(name string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY `%s`", t.wrapTable(), name)
}

func createIndexName(tableName string, indexType string, columns ...string) string {
	index := fmt.Sprintf("%s_%s_%s", tableName, strings.Join(columns, "_"), indexType)
	index = strings.ReplaceAll(index, "-", "_")
	index = strings.ReplaceAll(index, ".", "_")

	return index
}

func columnize(columns []string) string {
	return "`" + strings.Join(columns, "`, `") + "`"
}

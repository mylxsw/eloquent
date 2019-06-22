package migrate

import (
	"fmt"
	"strings"

	"github.com/mylxsw/eloquent/migrate/schema"
)

type tableBuilder struct {
	name      string
	prefix    string
	engine    string
	charset   string
	collation string
	temporary bool
	columns   []schema.ColumnType
	commands  []schema.Command
}

func NewTableBuilder(name string, prefix string) schema.TableBuilder {
	return &tableBuilder{
		name:      name,
		prefix:    prefix,
		columns:   make([]schema.ColumnType, 0),
		commands:  make([]schema.Command, 0),
		charset:   "utf8mb4",
		collation: "utf8mb4_unicode_ci",
	}
}

func (t *tableBuilder) addImpliedCommands() {
	if !t.creating() {
		if len(t.getAddedColumns()) > 0 {
			t.addCommand("add")
		}

		if len(t.getChangedColumns()) > 0 {
			t.addCommand("change")
		}
	}

	// TODO
}

func (t *tableBuilder) getAddedColumns() []schema.ColumnType {
	var cols = make([]schema.ColumnType, 0)
	for _, c := range t.columns {
		if !c.IsChange() {
			cols = append(cols, c)
		}
	}

	return cols
}

func (t *tableBuilder) getChangedColumns() []schema.ColumnType {
	var cols = make([]schema.ColumnType, 0)
	for _, c := range t.columns {
		if c.IsChange() {
			cols = append(cols, c)
		}
	}

	return cols
}

func (t *tableBuilder) Build() string {
	t.addImpliedCommands()

	sqlStrs := make([]string, 0)
	for _, c := range t.commands {
		sqlStrs = append(sqlStrs, c.Build())
	}

	return strings.Join(sqlStrs, ";\n") + ";\n"
}

func (t *tableBuilder) wrapTable() string {
	return "`" + t.name + "`"
}

func (t *tableBuilder) getColumns() []string {
	cols := make([]string, 0)
	for _, c := range t.getAddedColumns() {
		colStr := c.Build()
		if colStr != "" {
			cols = append(cols, colStr)
		}
	}

	return cols
}

func (t *tableBuilder) Engine(engine string) {
	t.engine = engine
}

func (t *tableBuilder) Charset(charset string) {
	t.charset = charset
}

func (t *tableBuilder) Collation(collation string) {
	t.collation = collation
}

func (t *tableBuilder) Temporary() {
	t.temporary = true
}

func (t *tableBuilder) BigIncrements(name string) schema.ColumnType {
	return t.UnsignedBigInteger(name, true)
}

func (t *tableBuilder) BigInteger(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "bigint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *tableBuilder) Binary(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "blob",
	})
}

func (t *tableBuilder) Boolean(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "tinyint(1)",
	})
}

func (t *tableBuilder) Char(name string, length int) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: fmt.Sprintf("char(%d)", length),
	})
}

func (t *tableBuilder) Date(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "date",
	})
}

func (t *tableBuilder) DateTime(name string, precision int) schema.ColumnType {
	dateType := "datetime"
	if precision > 0 {
		dateType += fmt.Sprintf("(%d)", precision)
	}

	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: dateType,
	})
}

func (t *tableBuilder) DateTimeTz(name string, precision int) schema.ColumnType {
	return t.DateTime(name, precision)
}

func (t *tableBuilder) Decimal(name string, total int, scale int) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:     name,
		ColumnUnsigned: false,
		ColumnType:     fmt.Sprintf("decimal(%d, %d)", total, scale),
	})
}

func (t *tableBuilder) Double(name string, total int, scale int) schema.ColumnType {
	fieldType := "double"
	if total != 0 && scale != 0 {
		fieldType += fmt.Sprintf("(%d, %d)", total, scale)
	}

	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: fieldType,
	})
}

func (t *tableBuilder) Enum(name string, items ...string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: fmt.Sprintf("enum('%s')", strings.Join(items, "', '")),
	})
}

func (t *tableBuilder) Float(name string, total int, scale int) schema.ColumnType {
	return t.Double(name, total, scale)
}

func (t *tableBuilder) Geometry(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "geometry",
	})
}

func (t *tableBuilder) GeometryCollection(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "geometrycollection",
	})
}

func (t *tableBuilder) Increments(name string) schema.ColumnType {
	return t.UnsignedInteger(name, true)
}

func (t *tableBuilder) Integer(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "int",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *tableBuilder) IpAddress(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "varchar(45)",
	})
}

func (t *tableBuilder) Json(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "json",
	})
}

func (t *tableBuilder) Jsonb(name string) schema.ColumnType {
	return t.Json(name)
}

func (t *tableBuilder) LineString(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "linestring",
	})
}

func (t *tableBuilder) LongText(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "longtext",
	})
}

func (t *tableBuilder) MacAddress(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "varchar(17)",
	})
}

func (t *tableBuilder) MediumIncrements(name string) schema.ColumnType {
	return t.UnsignedMediumInteger(name, true)
}

func (t *tableBuilder) MediumInteger(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "mediumint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *tableBuilder) MediumText(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "mediumtext",
	})
}

func (t *tableBuilder) Morphs(name string, indexName string) {
	t.String(name+"_type", 255)
	t.UnsignedBigInteger(name+"_id", false)
	t.Index(indexName, name+"_type", name+"_id")
}

func (t *tableBuilder) NullableMorphs(name string, indexName string) {
	t.String(name+"_type", 255).Nullable(true)
	t.UnsignedBigInteger(name+"_id", false).Nullable(true)
	t.Index(indexName, name+"_type", name+"_id")
}

func (t *tableBuilder) DropMorphs(name string, indexName string) {
	t.DropIndex(indexName)
	t.DropColumn(name+"_type", name+"_id")
}

func (t *tableBuilder) DropColumn(columns ...string) {
	panic("implement me")
}

func (t *tableBuilder) MultiLineString(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "multilinestring",
	})
}

func (t *tableBuilder) MultiPoint(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "multipoint",
	})
}

func (t *tableBuilder) MultiPolygon(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "multipolygon",
	})
}

func (t *tableBuilder) NullableTimestamps(precision int) {
	t.Timestamps(precision)
}

func (t *tableBuilder) Point(name string, srid int) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "point",
	})
}

func (t *tableBuilder) Polygon(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "polygon",
	})
}

func (t *tableBuilder) RememberToken() schema.ColumnType {
	return t.String("remember_token", 100).Nullable(true)
}

func (t *tableBuilder) DropRememberToken() {
	t.DropColumn("remember_token")
}

func (t *tableBuilder) SmallIncrements(name string) schema.ColumnType {
	return t.UnsignedSmallInteger(name, true)
}

func (t *tableBuilder) SmallInteger(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "smallint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *tableBuilder) SoftDeletes(column string, precision int) schema.ColumnType {
	return t.Timestamp(column, precision).Nullable(true)
}

func (t *tableBuilder) SoftDeletesTz(column string, precision int) schema.ColumnType {
	return t.TimestampTz(column, precision).Nullable(true)
}

func (t *tableBuilder) String(name string, length int) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: fmt.Sprintf("VARCHAR(%d)", length),
	})
}

func (t *tableBuilder) Text(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "text",
	})
}

func (t *tableBuilder) Time(name string, precision int) schema.ColumnType {
	dateType := "time"
	if precision > 0 {
		dateType += fmt.Sprintf("(%d)", precision)
	}

	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: dateType,
	})
}

func (t *tableBuilder) TimeTz(name string, precision int) schema.ColumnType {
	return t.Time(name, precision)
}

func (t *tableBuilder) Timestamp(name string, precision int) schema.ColumnType {
	dateType := "timestamp"
	if precision > 0 {
		dateType += fmt.Sprintf("(%d)", precision)
	}

	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: dateType,
	})
}

func (t *tableBuilder) TimestampTz(name string, precision int) schema.ColumnType {
	return t.Timestamp(name, precision)
}

func (t *tableBuilder) Timestamps(precision int) {
	t.Timestamp("created_at", precision).Nullable(true)
	t.Timestamp("updated_at", precision).Nullable(true)
}

func (t *tableBuilder) TimestampsTz(precision int) {
	t.TimestampTz("created_at", precision).Nullable(true)
	t.TimestampTz("updated_at", precision).Nullable(true)
}

func (t *tableBuilder) TinyIncrements(name string) schema.ColumnType {
	return t.UnsignedTinyInteger(name, true)
}

func (t *tableBuilder) TinyInteger(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "tinyint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      unsigned,
	})
}

func (t *tableBuilder) UnsignedBigInteger(name string, autoIncrement bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "bigint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *tableBuilder) UnsignedDecimal(name string, total int, scale int) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:     name,
		ColumnUnsigned: true,
		ColumnType:     fmt.Sprintf("decimal(%d, %d)", total, scale),
	})
}

func (t *tableBuilder) UnsignedInteger(name string, autoIncrement bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "int",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *tableBuilder) UnsignedMediumInteger(name string, autoIncrement bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "mediumint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *tableBuilder) UnsignedSmallInteger(name string, autoIncrement bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "smallint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *tableBuilder) UnsignedTinyInteger(name string, autoIncrement bool) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName:          name,
		ColumnType:          "tinyint",
		ColumnAutoIncrement: autoIncrement,
		ColumnUnsigned:      true,
	})
}

func (t *tableBuilder) Uuid(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "char(36)",
	})
}

func (t *tableBuilder) Year(name string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: "year",
	})
}

func (t *tableBuilder) Set(name string, items ...string) schema.ColumnType {
	return t.addColumn(&Column{
		ColumnName: name,
		ColumnType: fmt.Sprintf("set('%s')", strings.Join(items, "', '")),
	})
}

func (t *tableBuilder) addColumn(c schema.ColumnType) schema.ColumnType {
	t.columns = append(t.columns, c)

	return c
}

func (t *tableBuilder) Unique(name string, columns ...string) schema.Command {
	return t.indexCommand("unique", name, columns...)
}

func (t *tableBuilder) DropUnique(name string) schema.Command {
	return t.indexCommand("dropUnique", name)
}

func (t *tableBuilder) Index(name string, columns ...string) schema.Command {
	return t.indexCommand("index", name, columns...)
}

func (t *tableBuilder) DropIndex(name string) schema.Command {
	return t.indexCommand("dropIndex", name)
}

func (t *tableBuilder) Primary(name string, columns ...string) schema.Command {
	return t.indexCommand("primary", name, columns...)
}

func (t *tableBuilder) DropPrimary(name string) schema.Command {
	return t.indexCommand("dropPrimary", name)
}

func (t *tableBuilder) SpatialIndex(name string, columns ...string) schema.Command {
	return t.indexCommand("spatialIndex", name, columns...)
}

func (t *tableBuilder) DropSpatialIndex(name string) schema.Command {
	return t.indexCommand("dropSpatialIndex", name)
}

func (t *tableBuilder) Drop() schema.Command {
	return t.addCommand("drop")
}

func (t *tableBuilder) Create() schema.Command {
	return t.addCommand("create")
}

func (t *tableBuilder) DropIfExists() schema.Command {
	return t.addCommand("dropIfExists")
}

func (t *tableBuilder) Rename(to string) schema.Command {
	return t.addCommand("rename", to)
}

func (t *tableBuilder) indexCommand(indexType string, indexName string, columns ...string) schema.Command {
	if indexName == "" {
		indexName = createIndexName(t.name, indexType, columns...)
	}

	return t.addCommand(indexType, columns...).Index(indexName)
}

func (t *tableBuilder) addCommand(name string, parameters ...string) schema.Command {
	cmd := NewCommand(t)
	cmd.Name(name).Columns(parameters...)

	t.commands = append(t.commands, cmd)

	return cmd
}

func (t *tableBuilder) creating() bool {
	for _, c := range t.commands {
		if c.Equal("create") {
			return true
		}
	}

	return false
}

func (t *tableBuilder) compileKey(c *Command, indexType string) string {
	alg := ""
	if c.CommandAlgorithm != "" {
		alg = " USING " + c.CommandAlgorithm
	}

	return fmt.Sprintf(
		"ALTER TABLE %s ADD %s %s%s(`%s`)",
		t.wrapTable(),
		indexType,
		c.CommandIndex,
		alg,
		strings.Join(c.CommandParameters, "`, `"),
	)
}

func (t *tableBuilder) compileDrop(c *Command) string {
	return fmt.Sprintf("DROP TABLE %s", t.wrapTable())
}

func (t *tableBuilder) compileDropIfExists(c *Command) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", t.wrapTable())
}

func (t *tableBuilder) compileDropColumn(c *Command) string {
	dropColumns := make([]string, len(c.CommandParameters))
	for i, cc := range c.CommandParameters {
		dropColumns[i] = "DROP " + cc
	}

	return fmt.Sprintf("ALTER TABLE %s %s", t.wrapTable(), strings.Join(dropColumns, ", "))
}

func (t *tableBuilder) compileDropPrimary(c *Command) string {
	return fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY", t.wrapTable())
}

func (t *tableBuilder) compileDropIndex(c *Command) string {
	return fmt.Sprintf("ALTER TABLE %s DROP INDEX `%s`", t.wrapTable(), c.CommandIndex)
}

func (t *tableBuilder) compileRename(c *Command) string {
	return fmt.Sprintf("RENAME TABLE %s TO `%s`", t.wrapTable(), c.CommandParameters[0])
}

func (t *tableBuilder) compileCreateCommand() string {
	sqlStr := t.compileCreateTable()
	sqlStr += t.compileCreateEncoding()
	sqlStr += t.compileCreateEngine()

	return sqlStr
}

func (t *tableBuilder) compileCreateEngine() string {
	if t.engine != "" {
		return " ENGINE = " + t.engine
	}

	return ""
}

func (t *tableBuilder) compileCreateEncoding() string {
	sqlStr := ""
	if t.charset != "" {
		sqlStr += " DEFAULT CHARACTER SET " + t.charset
	}

	if t.collation != "" {
		sqlStr += " COLLATE " + t.collation
	}

	return sqlStr
}

func (t *tableBuilder) compileCreateTable() string {
	createStatement := "CREATE"
	if t.temporary {
		createStatement += " TEMPORARY"
	}

	return fmt.Sprintf(
		"%s TABLE %s (%s)",
		createStatement,
		t.wrapTable(),
		strings.Join(t.getColumns(), ", "),
	)
}

func (t *tableBuilder) compileAdd() string {
	var cols = make([]string, 0)
	for _, c := range t.getColumns() {
		cols = append(cols, "ADD " + c)
	}

	return fmt.Sprintf("ALTER TABLE %s %s", t.wrapTable(), strings.Join(cols, ", "))
}

func (t *tableBuilder) compileChange() string {
	// TODO
	return ""
}

func (t *tableBuilder) compileRenameColumn() string {
	// TODO
	return ""
}

func createIndexName(tableName string, indexType string, columns ...string) string {
	index := fmt.Sprintf("%s_%s_%s", tableName, strings.Join(columns, "_"), indexType)
	index = strings.ReplaceAll(index, "-", "_")
	index = strings.ReplaceAll(index, ".", "_")

	return index
}

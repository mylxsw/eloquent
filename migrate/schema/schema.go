package schema

type Schema interface {
	Create(table string, apply func(builder TableBuilder))
	Drop(table string)
	Table(table string, apply func(builder TableBuilder))
}

type ExprType int

type Expr struct {
	Type  ExprType
	Value string
}

const (
	ExprTypeString ExprType = iota
	ExprTypeRaw
)

type Command interface {
	Name(name string) Command
	Columns(columns ...string) Command
	Algorithm(algorithm string) Command
	Index(name string) Command
	Equal(name string) bool
	Build() string
}

type ColumnType interface {
	// Nullable Allows (by default) NULL values to be inserted into the column
	Nullable(value bool) ColumnType
	// After Place the column "after" another column (MySQL)
	After(name string) ColumnType
	// ColumnAutoIncrement Set INTEGER columns as auto-increment (primary key)
	AutoIncrement() ColumnType
	// Charset Specify a character set for the column (MySQL)
	Charset(charset string) ColumnType
	// Collation Specify a collation for the column (MySQL/SQL Server)
	Collation(collation string) ColumnType
	// Comment Add a comment to a column (MySQL/PostgreSQL)
	Comment(comment string) ColumnType
	// Default Specify a "default" value for the column
	Default(defaultVal Expr) ColumnType
	// First Place the column "first" in the table (MySQL)
	First() ColumnType
	// StoredAs Create a stored generated column (MySQL)
	StoredAs(expression string) ColumnType
	// ColumnUnsigned Set INTEGER columns as UNSIGNED (MySQL)
	Unsigned() ColumnType
	// UseCurrent Set TIMESTAMP columns to use CURRENT_TIMESTAMP as default value
	UseCurrent() ColumnType
	// VirtualAs Create a virtual generated column (MySQL)
	VirtualAs(expression string) ColumnType
	// GeneratedAs Create an identity column with specified sequence options (PostgreSQL)
	GeneratedAs(expression string) ColumnType
	// Always Defines the precedence of sequence values over input for an identity column (PostgreSQL)
	Always() ColumnType

	// Index add an index
	// Index(name string) ColumnType
	// Primary add a primary index
	// Primary() bool

	// Change modify some existing column types to a new type or modify the column's attributes
	Change() ColumnType
	IsChange() bool

	Build() string
	Type() string
}

type TableBuilder interface {
	Engine(engine string)
	Charset(charset string)
	Collation(collation string)
	Temporary()

	BigIncrements(name string) ColumnType
	BigInteger(name string, autoIncrement bool, unsigned bool) ColumnType
	Binary(name string) ColumnType
	Boolean(name string) ColumnType
	Char(name string, length int) ColumnType
	Date(name string) ColumnType
	DateTime(name string, precision int) ColumnType
	DateTimeTz(name string, precision int) ColumnType
	Decimal(name string, total int, scale int) ColumnType
	Double(name string, total int, scale int) ColumnType
	Enum(name string, items ...string) ColumnType
	Float(name string, total int, scale int) ColumnType
	Geometry(name string) ColumnType
	GeometryCollection(name string) ColumnType
	Increments(name string) ColumnType
	Integer(name string, autoIncrement bool, unsigned bool) ColumnType
	IpAddress(name string) ColumnType
	Json(name string) ColumnType
	Jsonb(name string) ColumnType
	LineString(name string) ColumnType
	LongText(name string) ColumnType
	MacAddress(name string) ColumnType
	MediumIncrements(name string) ColumnType
	MediumInteger(name string, autoIncrement bool, unsigned bool) ColumnType
	MediumText(name string) ColumnType
	Morphs(name string, indexName string)
	MultiLineString(name string) ColumnType
	MultiPoint(name string) ColumnType
	MultiPolygon(name string) ColumnType
	NullableMorphs(name string, indexName string)
	DropMorphs(name string, indexName string)
	DropColumn(columns ...string)
	NullableTimestamps(precision int)
	Point(name string, srid int) ColumnType
	Polygon(name string) ColumnType
	RememberToken() ColumnType
	DropRememberToken()
	SmallIncrements(name string) ColumnType
	SmallInteger(name string, autoIncrement bool, unsigned bool) ColumnType
	SoftDeletes(column string, precision int) ColumnType
	SoftDeletesTz(column string, precision int) ColumnType
	String(name string, length int) ColumnType
	Text(name string) ColumnType
	Time(name string, precision int) ColumnType
	TimeTz(name string, precision int) ColumnType
	Timestamp(name string, precision int) ColumnType
	TimestampTz(name string, precision int) ColumnType
	Timestamps(precision int)
	TimestampsTz(precision int)
	TinyIncrements(name string) ColumnType
	TinyInteger(name string, autoIncrement bool, unsigned bool) ColumnType
	UnsignedBigInteger(name string, autoIncrement bool) ColumnType
	UnsignedDecimal(name string, total int, scale int) ColumnType
	UnsignedInteger(name string, autoIncrement bool) ColumnType
	UnsignedMediumInteger(name string, autoIncrement bool) ColumnType
	UnsignedSmallInteger(name string, autoIncrement bool) ColumnType
	UnsignedTinyInteger(name string, autoIncrement bool) ColumnType
	Uuid(name string) ColumnType
	Year(name string) ColumnType
	Set(name string, items ...string) ColumnType

	// Drop add drop table command
	Drop() Command
	// DropIfExists add dropIfExists command
	DropIfExists() Command
	// Rename add rename table command
	Rename(to string) Command
	// Create add create table command
	Create() Command

	// Unique specifies a column's values should be unique
	Unique(name string, columns ...string) Command
	// DropUnique drop a unique index
	DropUnique(name string) Command
	// Index create a index
	Index(name string, columns ...string) Command
	// DropIndex drop a index
	DropIndex(name string) Command
	// Primary adds a primary key
	Primary(name string, columns ...string) Command
	// DropPrimary drop a primary key
	DropPrimary(name string) Command
	// SpatialIndex adds a spatial index
	SpatialIndex(name string, columns ...string) Command
	// DropSpatialIndex drop a spatial index
	DropSpatialIndex(name string) Command

	Build() string
}

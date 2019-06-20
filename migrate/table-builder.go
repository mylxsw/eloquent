package migrate

import (
	"strings"

	"github.com/mylxsw/eloquent/migrate/column"
	"github.com/mylxsw/eloquent/migrate/schema"
)

type tableBuilder struct {
	engine    string
	charset   string
	collation string
	temporary bool
	columns   []schema.ColumnType
}

func NewTableBuilder() schema.TableBuilder {
	return &tableBuilder{
		columns: make([]schema.ColumnType, 0),
	}
}

func (t *tableBuilder) Build() string {

	cols := make([]string, 0)
	for _, c := range t.columns {
		colStr := c.Build()
		if colStr != "" {
			cols = append(cols, colStr)
		}
	}

	return strings.Join(cols, ",\n")
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
	panic("implement me")
}

func (t *tableBuilder) BigInteger(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Binary(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Boolean(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Char(name string, length int) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Date(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) DateTime(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) DateTimeTz(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Decimal(name string, total int, scale int) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Double(name string, total int, scale int) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Enum(name string, items ...string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Float(name string, total int, scale int) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Geometry(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) GeometryCollection(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Increments(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Integer(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	return t.addColumn(&column.IntegerColumn{ColumnName: name, ColumnAutoIncrement: autoIncrement, ColumnUnsigned: unsigned,})
}

func (t *tableBuilder) IpAddress(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Json(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Jsonb(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) LineString(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) LongText(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) MacAddress(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) MediumIncrements(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) MediumInteger(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) MediumText(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Morphs(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) MultiLineString(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) MultiPoint(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) MultiPolygon(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) NullableMorphs(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) NullableTimestamps() schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Point(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Polygon(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) RememberToken() schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) SmallIncrements(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) SmallInteger(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) SoftDeletes() schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) SoftDeletesTz() schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) String(name string, length int) schema.ColumnType {
	return t.addColumn(&column.StringColumn{ColumnName: name, ColumnLength: length})
}

func (t *tableBuilder) Text(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Time(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) TimeTz(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Timestamp(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) TimestampTz(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Timestamps() schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) TimestampsTz() schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) TinyIncrements(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) TinyInteger(name string, autoIncrement bool, unsigned bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) UnsignedBigInteger(name string, autoIncrement bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) UnsignedDecimal(name string, total int, scale int) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) UnsignedInteger(name string, autoIncrement bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) UnsignedMediumInteger(name string, autoIncrement bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) UnsignedSmallInteger(name string, autoIncrement bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) UnsignedTinyInteger(name string, autoIncrement bool) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Uuid(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Year(name string) schema.ColumnType {
	panic("implement me")
}

func (t *tableBuilder) Set(name string, items ...string) schema.ColumnType {
	panic("xxx")
}

func (t *tableBuilder) Unique(name ...string) {
	panic("implement me")
}

func (t *tableBuilder) DropUnique(name ...string) {
	panic("implement me")
}

func (t *tableBuilder) Index(name ...string) {
	panic("implement me")
}

func (t *tableBuilder) DropIndex(name ...string) {
	panic("implement me")
}

func (t *tableBuilder) Primary(name ...string) {
	panic("implement me")
}

func (t *tableBuilder) DropPrimary(name ...string) {
	panic("implement me")
}

func (t *tableBuilder) SpatialIndex(name ...string) {
	panic("implement me")
}

func (t *tableBuilder) DropSpatialIndex(name ...string) {
	panic("implement me")
}

func (t *tableBuilder) addColumn(c schema.ColumnType) schema.ColumnType {
	t.columns = append(t.columns, c)

	return c
}


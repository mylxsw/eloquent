package migrate

import (
	"testing"

	"github.com/mylxsw/eloquent/migrate/schema"
)

func TestSchema_Create(t *testing.T) {
	NewSchema().Create("test", func(builder schema.TableBuilder) {
		builder.Integer("id", true, true).Comment("自增主键")
		builder.String("username", 255).Nullable(true).Comment("用户名")
	})
}

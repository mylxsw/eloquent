package migrate

import (
	"testing"

	"github.com/mylxsw/eloquent/migrate/schema"
)

func TestSchema_Create(t *testing.T) {
	NewSchema().Create("wz_users", func(builder schema.TableBuilder) {
		builder.Increments("id")
		builder.String("name", 255)
		builder.String("email", 255)
		builder.String("password", 255)
		builder.Timestamps(0)
		builder.RememberToken()

		builder.Unique("", "email")
	})

	NewSchema().Create("wz_password_resets", func(builder schema.TableBuilder) {
		builder.String("email", 255)
		builder.String("token", 255)
		builder.Timestamp("created_at", 0).Nullable(true)
	})

	NewSchema().Create("wz_projects", func(builder schema.TableBuilder) {
		builder.Increments("id")
		builder.String("name", 100).Comment("项目名称")
		builder.Text("description").Nullable(true).Comment("项目描述")
		builder.UnsignedTinyInteger("visibility", false).Comment("可见性")
		builder.UnsignedInteger("user_id", false).Nullable(true).Comment("创建用户ID")
		builder.Timestamps(0)
		builder.SoftDeletes("deleted_at", 0)
	})
}

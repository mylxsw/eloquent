package migrate

import (
	"testing"
)

func TestSchema_Create(t *testing.T) {
	m := NewManager()
	m.Schema("201905141221").Create("wz_users", func(builder *Builder) {
		builder.Increments("id")
		builder.String("name", 255)
		builder.String("email", 255).Unique()
		builder.String("password", 255)
		builder.Timestamps(0)
		builder.RememberToken()

		builder.Index("", "name", "password")
	})

	m.Schema("201905160813").Create("wz_password_resets", func(builder *Builder) {
		builder.String("email", 255)
		builder.String("token", 255)
		builder.Timestamp("created_at", 0).Nullable(true)
	})

	m.Schema("201907150945").Create("wz_projects", func(builder *Builder) {
		builder.Increments("id")
		builder.String("name", 100).Comment("项目名称")
		builder.Text("description").Nullable(true).Comment("项目描述")
		builder.UnsignedTinyInteger("visibility", false).Comment("可见性")
		builder.UnsignedInteger("user_id", false).Nullable(true).Comment("创建用户ID")
		builder.Timestamps(0)
		builder.SoftDeletes("deleted_at", 0)
	})

	m.Schema("201908101010").Table("wz_users", func(builder *Builder) {
		builder.TinyInteger("status", false, true).
			Nullable(true).
			Default(StringExpr("1")).
			Comment("用户状态：0-未激活，1-已激活，2-已禁用")
	})

	m.Schema("201909190302").Table("wz_projects", func(builder *Builder) {
		builder.Drop()
		builder.DropColumn("description")
		builder.Rename("wz_project")
		builder.RenameColumn("user_id", "creator_id")
	})
}

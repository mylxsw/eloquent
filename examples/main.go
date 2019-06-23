package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mylxsw/eloquent/migrate"
)

func main() {

	connURI := "root:@tcp(127.0.0.1:3306)/eloquent_example?parseTime=true"
	db, err := sql.Open("mysql", connURI)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	createMigrate(db)

	// projectModel := models.NewProjectModel(db)
	//
	// id, err := projectModel.Save(models.Project{
	// 	Name: "test",
	// })
	// misc.AssertError(err)
	//
	// fmt.Println("insert id=", id)
	//
	// _, err = projectModel.UpdateById(id, models.Project{
	// 	Name:      "test2",
	// 	CatalogId: null.IntFrom(100),
	// })
	// misc.AssertError(err)
	//
	// _, err = projectModel.DeleteById(id)
	// misc.AssertError(err)
	//
	// projects, err := projectModel.Get(query.Builder().OrderBy("id", "desc").Limit(10))
	// misc.AssertError(err)
	//
	// for _, p := range projects {
	// 	fmt.Printf("%v\n", p)
	// }
}

func createMigrate(db *sql.DB) {
	m := migrate.NewManager(db).Init()
	m.Schema("201905141221").Create("wz_users", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("name", 255)
		builder.String("email", 255).Unique()
		builder.String("password", 255)
		builder.Timestamps(0)
		builder.RememberToken()

		builder.Index("", "name", "password")
	})

	m.Schema("201905160813").Create("wz_password_resets", func(builder *migrate.Builder) {
		builder.String("email", 255)
		builder.String("token", 255)
		builder.Timestamp("created_at", 0).Nullable(true)
	})

	m.Schema("201907150945").Create("wz_projects", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("name", 100).Comment("项目名称")
		builder.Text("description").Nullable(true).Comment("项目描述")
		builder.UnsignedTinyInteger("visibility", false).Comment("可见性")
		builder.UnsignedInteger("user_id", false).Nullable(true).Comment("创建用户ID")
		builder.Timestamps(0)
		builder.SoftDeletes("deleted_at", 0)
	})

	m.Schema("201908101010").Table("wz_users", func(builder *migrate.Builder) {
		builder.TinyInteger("status", false, true).
			Nullable(true).
			Default(migrate.StringExpr("1")).
			Comment("用户状态：0-未激活，1-已激活，2-已禁用")
	})

	m.Schema("201909190302").Table("wz_projects", func(builder *migrate.Builder) {
		builder.Text("description").Nullable(true).Comment("项目描述")
		// builder.RenameColumn("user_id", "creator_id")
	})

	m.Schema("201909040102").Table("wz_projects", func(builder *migrate.Builder) {
		builder.DropColumn("description")
	})

	m.Schema("201909040103").Drop("wz_projects")

	if err := m.Run(); err != nil {
		panic(err)
	}
}
package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mylxsw/asteria"
	"github.com/mylxsw/eloquent"
	"github.com/mylxsw/eloquent/event"
	"github.com/mylxsw/eloquent/examples/models"
	"github.com/mylxsw/eloquent/migrate"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/go-toolkit/events"
	"github.com/mylxsw/go-toolkit/misc"
)

var logger = asteria.Module("example").WithFileLine(true)

func main() {
	createEventDispatcher()

	connURI := "root:@tcp(127.0.0.1:3306)/eloquent_example?parseTime=true"
	db, err := sql.Open("mysql", connURI)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	createMigrate(db)
	databaseOperationExample(db)
	modelOperationExample(db)
}

func modelOperationExample(db *sql.DB) {
	err := eloquent.Transaction(db, func(tx query.Database) error {
		userModel := models.NewUserModel(tx)

		id, err := userModel.Save(models.User{
			Name:     "guan",
			Email:    "guan@aicode.cc",
			Password: "88959q",
		})
		misc.AssertError(err)

		logger.Infof("Insert User ID=%d", id)

		user, err := userModel.Find(id)
		misc.AssertError(err)

		logger.Infof("User id=%d, name=%s, email=%s", user.Id, user.Name, user.Email)

		roleId, err := user.Role().Create(models.Role{
			Name:        "admin",
			Description: "root user",
		})
		misc.AssertError(err)

		logger.Infof("Insert Role ID=%d", roleId)

		users, err := userModel.Get()
		misc.AssertError(err)

		for _, user := range users {
			logger.Infof("User id=%d, name=%s, email=%s, role_id=%d", user.Id, user.Name, user.Email, user.RoleId)
		}

		_, err = userModel.DeleteById(user.Id)
		misc.AssertError(err)

		c1, err := userModel.Count()
		misc.AssertError(err)

		c2, err := userModel.WithTrashed().Count()
		misc.AssertError(err)

		logger.Infof("After soft deleted count=%d/%d", c1, c2)

		_, err = userModel.ForceDeleteById(user.Id)
		misc.AssertError(err)

		c1, err = userModel.Count()
		misc.AssertError(err)

		c2, err = userModel.WithTrashed().Count()
		misc.AssertError(err)

		logger.Infof("After force deleted count=%d/%d", c1, c2)

		return nil
	})

	misc.AssertError(err)
}

func databaseOperationExample(db *sql.DB) {
	err := eloquent.Transaction(db, func(tx query.Database) error {
		id, err := eloquent.DB(tx).Insert(
			"wz_user",
			query.KV{
				"name":     "mylxsw",
				"email":    "mylxsw@aicode.cc",
				"password": "123455",
			},
		)
		misc.AssertError(err)

		logger.Infof("Insert ID=%d", id)

		id, err = eloquent.DB(tx).Insert(
			"wz_user",
			query.KV{
				"name":     "adanos",
				"email":    "adanos@aicode.cc",
				"password": "123455",
			},
		)
		misc.AssertError(err)

		logger.Infof("Insert ID=%d", id)

		res, err := eloquent.DB(tx).Query(
			eloquent.Build("wz_user").Select("id", "name", "email"),
			func(row eloquent.Scanner) (interface{}, error) {
				user := models.User{}
				err := row.Scan(&user.Id, &user.Name, &user.Email)

				return user, err
			},
		)
		misc.AssertError(err)

		res.Each(func(user models.User) {
			logger.Infof("user_id=%d, name=%s, email=%s", user.Id, user.Name, user.Email)
		})

		affected, err := eloquent.DB(tx).Delete(eloquent.Build("wz_user"))
		misc.AssertError(err)

		logger.Infof("Deleted rows %d", affected)

		return nil
	})

	misc.AssertError(err)
}

func createMigrate(db *sql.DB) {
	m := migrate.NewManager(db).Init()
	m.Schema("20190692901").Create("wz_user", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("name", 255)
		builder.String("email", 255).Unique()
		builder.String("password", 255)
		builder.Timestamps(0)
		builder.RememberToken()

		builder.Index("", "name", "password")
	})

	m.Schema("20190692902").Create("wz_password_reset", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("email", 255)
		builder.String("token", 255)
		builder.Timestamp("created_at", 0).Nullable(true)
	})

	m.Schema("20190692903").Table("wz_user", func(builder *migrate.Builder) {
		builder.TinyInteger("status", false, true).
			Nullable(true).
			Default(migrate.StringExpr("1")).
			Comment("用户状态：0-未激活，1-已激活，2-已禁用")
	})

	m.Schema("20190692904").Create("wz_role", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("name", 100).Comment("角色名")
		builder.Text("description").Comment("备注")
		builder.Timestamps(0)
	})

	m.Schema("20190692905").Table("wz_user", func(builder *migrate.Builder) {
		builder.Integer("role_id", false, false).Nullable(true).Comment("角色ID")
	})

	m.Schema("20190692906").Table("wz_user", func(builder *migrate.Builder) {
		builder.SoftDeletes("deleted_at", 0)
	})

	m.Schema("2019063001").Create("wz_enterprise", func(builder *migrate.Builder) {
		builder.Increments("id")
		builder.String("name", 255).Comment("企业名称")
		builder.String("address", 255).Comment("企业地址")
		builder.TinyInteger("status", false, false).Default(migrate.StringExpr("0")).Comment("企业状态：0-未审核 1-审核未通过 2-审核通过")

		builder.SoftDeletes("deleted_at", 0)
	})

	m.Schema("2019063002").Table("wz_user", func(builder *migrate.Builder) {
		builder.Integer("enterprise_id", false, true).Nullable(true).Comment("企业ID")
	})

	if err := m.Run(); err != nil {
		panic(err)
	}
}

func createEventDispatcher() {
	// create event listener
	eventManager := events.NewEventManager(events.NewMemoryEventStore(false))
	event.SetDispatcher(eventManager)

	// eventManager.Listen(func(evt event.MigrationStartedEvent) {
	// 	logger.Debugf("MigrationStartedEvent received: %s", evt.SQL)
	// })

	eventManager.Listen(func(evt event.QueryExecutedEvent) {
		logger.WithContext(asteria.C{
			"sql": evt.SQL,
			// "bindings": evt.Bindings,
			// "elapse":   evt.Time.Seconds(),
		}).Debugf("QueryExecutedEvent received")
	})

	// 	eventManager.Listen(func(evt event.TransactionBeginningEvent) {
	// 		logger.Debugf("Transaction starting")
	// 	})
	//
	// 	eventManager.Listen(func(evt event.TransactionCommittedEvent) {
	// 		logger.Debugf("Transaction committed")
	// 	})
	//
	// 	eventManager.Listen(func(evt event.TransactionRolledBackEvent) {
	// 		logger.Debugf("Transaction rollback")
	// 	})
}

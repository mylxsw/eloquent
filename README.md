# Eloquent ORM

Eloquent 是一款为 Golang 开发的基于代码生成的数据库 ORM 框架，它的设计灵感来源于著名的 PHP 开发框架 Laravel，支持 MySQL 等数据库。

## 模型定义

Eloquent 使用 YAML 文件来定义模型结构，通过代码生成的方式来创建模型对象。下面是模型定义的基本格式：

```yaml
package: models                   // 当前模型所在的包名
imports:                          // 要 import 的包，当 models.definition.fields 中包含外部类型时，在这里引入外部的包，可选
- github.com/mylxsw/eloquent
meta:
  table_prefix: el_               // 表前缀，默认为空，可选
models:                           // 模型定义部分
- name: user                      // 模型名称，也是默认表名，模型对象会转换为首字母大写的驼峰命名格式
  relations:                      // 模型关联定义，可选
  - model: role
    rel: n-1
    foreign_key: role_id
    owner_key: id
    local_key: ""
    table: ""
    package: ""
    method: ""
  definition:                     // 模型基本信息
    table_name: user              // 对应的数据库表名，不指定该选项则默认使用模型名
    without_create_time: false    // 不添加 created_at 字段，默认会自动添加该字段，类型为 time.Time，插入数据时自动更新
    without_update_time: false    // 不添加 updated_at 字段，默认会自动添加该字段，类型为 time.Time，更新数据时自动更新
    soft_delete: false            // 是否启用软删除支持，启用软删除后，会自动添加 deleted_at 字段，类型为 time.Time
    fields:                       // 表字段对应关系
    - name: id                    // 主键 ID
      type: int64                 // 主键 ID 类型
      tag: json:"id"              // 主键 ID 类型的 Tags
    - name: name                  // 其它字段名
      type: string                // 其它字段类型
      tag: json:"name"            // 其它字段的 Tags
    - name: age
      type: int64
      tag: json:"age"
```

你可以使用命令行工具 `eloquent create-model` 来生成模型定义文件

```bash
$ eloquent create-model -h
NAME:
    create-model - 创建表模型定义文件

USAGE:
    create-model [command options] [arguments...]

OPTIONS:
   --table value         表名
   --package value       包名
   --output value        输出目录
   --table-prefix value  表前缀
   --no-created_at       不自动添加 created_at 字段
   --no-updated_at       不自动添加 updated_at 字段
   --soft-delete         启用软删除支持
   --import value        引入包
```

比如创建一个用户模型定义文件，表名为 `users`，包名为 `models`，启用软删除支持，输出到当前目录

```bash
$ eloquent create-model --table users --package models --soft-delete --output .
$ ls -al
total 8
-rwxr-xr-x  1 mylxsw  wheel   162B  6 17 21:32 users.yml
```

模型定义创建之后，默认是只有主键 id 的，我们手动修改模型定义文件，添加几个额外的字段，

```yaml
package: models
models:
- name: users
  definition:
    table_name: users
    soft_delete: true
    fields:
    - name: id
      type: int64
      tag: json:"id"
    - name: name 
      type: string
    - name: age
      type: int
```

使用 `eloquent gen` 命令来生成模型文件

```bash
$ eloquent gen -h
NAME:
    gen - 根据模型文件定义生成模型对象

USAGE:
    gen [command options] [arguments...]

OPTIONS:
   --source value  模型定义所在文件的 Glob 表达式，比如 ./models/*.yml
```

比如上面我们创建的 `users.yml` 模型定义，执行下面的命令创建模型对象

```bash
$ eloquent gen --source './*.yml'
users.orm.go
$ ls -al
total 40
drwxr-xr-x   4 mylxsw  wheel    128  6 17 21:36 .
drwxrwxrwt  17 root    wheel    544  6 17 21:32 ..
-rwxr-xr-x   1 mylxsw  wheel  15614  6 17 21:36 users.orm.go
-rwxr-xr-x   1 mylxsw  wheel    162  6 17 21:32 users.yml
```

生成的模型对象会包含两个模型定义的结构体，一个是模型名称本身命名的结构体 Xxx，用于与数据库进行交互

```go
type User struct {
	Id            null.Int `json:"id"`
	Name          null.String
 	Age           null.Int
	CreatedAt     null.Time
	UpdatedAt     null.Time
	DeletedAt     null.Time
}
```

> 这里的结构体字段类型使用了 `gopkg.in/guregu/null.v3` 对基本类型的封装，解决 Golang 基本类型字段不支持 `null` 值的问题，但是这样也给使用者带来了不便，我们必须使用 `Name.ValueOrZero()` 这种方式来获取字段的基本类型值。

另一个结构体是 XxxPlain，该结构体将 `null` 值转换为了字段的基本类型，当数据库中为 null 时，会返回字段的默认值。通过模型对象的 `ToXxxPlain()` 方法可以将模型对象转换为该结构体。

```go
type UserPlain struct {
	Id            int64
	Name          string
 	Age           int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
```

## 模型使用

本文的讲解将基于前面创建的用户模型（`users.yml`），假定我们的模型目录为 `models`。在使用模型的 API 之前，需要先创建数据库连接对象，我们使用 Golang 标准库的 `database/sql` 包来创建

```go
db, err := sql.Open("mysql", connURI)
if err != nil {
  panic(err)
}

defer db.Close()
```

### 创建模型实例

创建模型定义文件和生成模型对象文件后，使用 `models.NewXxxModel(db query.Database)` 来创建一个模型对象。

```go
// 创建用户模型对象
userModel := models.NewUsersModel(db)
```

### 查询条件

在 Eloquent 中，查询条件使用 `query.Builder()` 方法来构建，该方法会生成一个 `SQLBuilder` 对象，使用它我们可以使用链式语法来构建查询条件

```go
builder := query.Builder()
```

`SQLBuilder` 对象提供了一系列的方法来帮助我们构建灵活的查询条件

- `WhereColumn(field, operator string, value string) Condition`
- `OrWhereColumn(field, operator string, value string) Condition`
- `OrWhereNotExist(subQuery SubQuery) Condition`
- `OrWhereExist(subQuery SubQuery) Condition`
- `WhereNotExist(subQuery SubQuery) Condition`
- `WhereExist(subQuery SubQuery) Condition`
- `OrWhereNotNull(field string) Condition`
- `OrWhereNull(field string) Condition`
- `WhereNotNull(field string) Condition`
- `WhereNull(field string) Condition`
- `OrWhereRaw(raw string, items ...interface{}) Condition`
- `WhereRaw(raw string, items ...interface{}) Condition`
- `OrWhereNotIn(field string, items ...interface{}) Condition`
- `OrWhereIn(field string, items ...interface{}) Condition`
- `WhereNotIn(field string, items ...interface{}) Condition`
- `WhereIn(field string, items ...interface{}) Condition`
- `WhereGroup(wc ConditionGroup) Condition`
- `OrWhereGroup(wc ConditionGroup) Condition`
- `Where(field string, value ...interface{}) Condition`
- `OrWhere(field string, value ...interface{}) Condition`
- `WhereBetween(field string, min, max interface{}) Condition`
- `WhereNotBetween(field string, min, max interface{}) Condition`
- `OrWhereBetween(field string, min, max interface{}) Condition`
- `OrWhereNotBetween(field string, min, max interface{}) Condition`
- `WhereCondition(cond sqlCondition) Condition`
- `When(when When, cg ConditionGroup) Condition`
- `OrWhen(when When, cg ConditionGroup) Condition`
- `Get() []sqlCondition`
- `Append(cond Condition) Condition`
- `Resolve(tableAlias string) (string, []interface{})`

比如我们要查询用户名模糊匹配 `Tom`，年龄大于 30 岁的用户

```go
query.Builder().Where(model.UserFieldName, "LIKE", "%Tom%").Where("age", ">", 30)
```

### CRUD

查询用户列表

```go
users, err := model.NewUsersModel(s.db).Get()
for _, user := range users {
  // user.Id.ValueOrZero() 
  // user.Name.ValueOrZero()

  // user.ToUserPlain().Name
}
```

查询第一个匹配的用户

```go
user, err := model.NewUserModel(s.db).First(query.Builder().Where(model.UserFieldName, username))
```

创建一个用户

```go
userID, err := model.NewUserModel(s.db).Save(model.User{
  Account:  null.StringFrom(username),
  Status:   null.IntFrom(int64(UserStatusEnabled)),
  Password: null.StringFrom(password),
})

// 也可以这样
userID, err := model.NewUserModel(s.db).Create(query.KV{
  model.UserFieldUuid:     userInfo.Uuid,
  model.UserFieldName:     userInfo.Name,
  model.UserFieldAccount:  username,
  model.UserFieldStatus:   userInfo.Status,
  model.UserFieldPassword: userInfo.Password,
})

// 还可以这样
user := models.UserPlain{
  Name: "Tom",
  Age: 32,
}

userID, err := userModel.Save(user.ToUser())
```

## 数据库迁移

Eloquent 支持与 Laravel 框架类似的数据库迁移功能，使用语法也基本一致。

```go
m := migrate.NewManager(db).Init()

m.Schema("202104222322").Create("user", func(builder *migrate.Builder) {
  builder.Increments("id")
  builder.String("uuid", 255).Comment("用户 uuid")
  builder.String("name", 100).Comment("用户名")
  builder.Timestamps(0)
})
m.Schema("202106100943").Table("user", func(builder *migrate.Builder) {
  builder.String("account", 100).Comment("账号名")
  builder.TinyInteger("status", false, true).Default(migrate.RawExpr("1")).Comment("状态：0-禁用 1-启用")
})
m.Schema("202106102309").Table("user", func(builder *migrate.Builder) {
  builder.String("password", 256).Nullable(true).Comment("密码")
})

if err := m.Run(); err != nil {
  panic(err)
}
```

## 示例项目

- [tech-share](https://github.com/mylxsw/tech-share) 这是一个简单的 web 项目，用于企业内部对技术分享的管理
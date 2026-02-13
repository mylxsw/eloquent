# Eloquent ORM for Go

> **[English](./README.md)** | 中文

Eloquent 是一款为 Go 语言开发的基于**代码生成**的数据库 ORM 框架，设计灵感来源于 PHP 框架 Laravel 的 Eloquent ORM。支持 MySQL 和 PostgreSQL 数据库。

## 目录

- [核心概念](#核心概念)
- [快速开始](#快速开始)
- [安装](#安装)
- [模型定义与代码生成](#模型定义与代码生成)
- [数据库方言 (MySQL / PostgreSQL)](#数据库方言-mysql--postgresql)
- [数据库连接](#数据库连接)
- [CRUD 操作](#crud-操作)
- [查询构建器](#查询构建器)
- [模型关联](#模型关联)
- [软删除](#软删除)
- [作用域](#作用域)
- [数据库迁移](#数据库迁移)
- [事务](#事务)
- [事件系统](#事件系统)
- [直接使用数据库操作](#直接使用数据库操作)
- [示例项目](#示例项目)

## 核心概念

在开始使用之前，理解以下核心概念将帮助你更好地掌握 Eloquent：

### 1. 代码生成驱动

Eloquent 不使用运行时反射来映射数据库表，而是通过 **YAML 文件定义模型 → 命令行工具生成 Go 代码** 的方式工作。这意味着：
- 你需要先编写一个 YAML 文件描述你的数据库表结构
- 然后运行 `eloquent gen` 命令生成对应的 Go 代码（`*.orm.go` 文件）
- 生成的代码包含类型安全的模型结构体和完整的 CRUD 方法

### 2. 两种模型结构体

对于每个模型（例如 `User`），Eloquent 会生成两个结构体：

- **`UserN`**（Nullable 版本）：所有字段使用 `null.Int`、`null.String` 等类型（来自 `gopkg.in/guregu/null.v3`），能够正确处理数据库中的 `NULL` 值。这是与数据库交互时使用的主要结构体。
- **`User`**（Plain 版本）：所有字段使用 Go 基本类型（`int64`、`string` 等），当数据库字段为 `NULL` 时返回零值。适合在业务逻辑中使用。

两者可以相互转换：

```go
// UserN → User (Nullable → Plain)
user := userN.ToUser()

// User → UserN (Plain → Nullable)
userN := user.ToUserN()
```

### 3. 查询构建器 (SQLBuilder)

所有查询条件都通过 `query.Builder()` 链式调用来构建，而不是手写 SQL。查询构建器支持 `Where`、`WhereIn`、`OrderBy`、`Limit` 等常见操作。

### 4. 上下文 (Context)

所有数据库操作都需要传入 `context.Context`，这是 Go 语言的最佳实践，便于超时控制和取消操作。

## 快速开始

下面用一个完整的例子展示从定义模型到操作数据库的完整流程。

### 第一步：安装命令行工具

```bash
go install github.com/mylxsw/eloquent/cmd/orm@latest
```

安装后你将获得一个 `orm` 命令行工具（也可以自行编译项目，生成 `eloquent` 二进制文件）。

### 第二步：创建模型定义文件

可以手动编写 YAML，也可以用命令行生成一个基础模板：

```bash
eloquent create-model --table users --package models --soft-delete --output ./models/
```

然后编辑生成的 `users.yaml` 文件，添加你需要的字段：

```yaml
package: models
models:
  - name: User
    definition:
      table_name: users
      soft_delete: true
      fields:
        - name: id
          type: int64
          tag: 'json:"id"'
        - name: name
          type: string
          tag: 'json:"name"'
        - name: email
          type: string
          tag: 'json:"email"'
        - name: age
          type: int64
          tag: 'json:"age"'
```

### 第三步：生成模型代码

```bash
eloquent gen --source './models/*.yaml'
```

这会在同目录下生成 `users.orm.go` 文件，包含 `UserN`、`User`、`UserModel` 等所有代码。

> **注意：** 生成的 `*.orm.go` 文件不应手动修改，每次执行 `gen` 命令都会重新生成。

### 第四步：添加依赖

```bash
go get github.com/mylxsw/eloquent
go get gopkg.in/guregu/null.v3

# MySQL 驱动
go get github.com/go-sql-driver/mysql

# PostgreSQL 驱动（二选一）
go get github.com/lib/pq
# 或者: go get github.com/jackc/pgx/v5
```

### 第五步：编写业务代码

```go
package main

import (
    "context"
    "database/sql"
    "fmt"

    "yourproject/models"

    _ "github.com/go-sql-driver/mysql"
    "gopkg.in/guregu/null.v3"
)

func main() {
    db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/mydb?parseTime=true")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    ctx := context.TODO()
    userModel := models.NewUserModel(db)

    // 创建用户
    id, err := userModel.Save(ctx, models.UserN{
        Name:  null.StringFrom("张三"),
        Email: null.StringFrom("zhangsan@example.com"),
        Age:   null.IntFrom(25),
    })
    fmt.Printf("新用户ID: %d\n", id)

    // 查询用户
    user, err := userModel.Find(ctx, id)
    fmt.Printf("用户名: %s, 邮箱: %s\n", user.Name.String, user.Email.String)

    // 转换为 Plain 结构体使用
    u := user.ToUser()
    fmt.Printf("用户名: %s, 年龄: %d\n", u.Name, u.Age)
}
```

## 安装

在你的项目中添加 Eloquent 依赖：

```bash
go get github.com/mylxsw/eloquent
```

编译命令行工具（用于代码生成）：

```bash
# 方式一：直接安装
go install github.com/mylxsw/eloquent/cmd/orm@latest

# 方式二：克隆项目后编译
git clone https://github.com/mylxsw/eloquent.git
cd eloquent
make build
# 生成的二进制文件在 ./bin/eloquent
```

## 模型定义与代码生成

### YAML 模型定义文件

模型定义使用 YAML 格式，以下是完整的字段说明：

```yaml
package: models                      # 生成代码的 Go 包名
imports:                             # 需要额外 import 的包（当字段使用外部类型时）
  - github.com/mylxsw/eloquent

meta:
  table_prefix: wz_                  # 表前缀，会自动拼接到表名前

models:
  - name: User                       # 模型名称（生成的结构体名称）
    definition:
      table_name: users              # 数据库表名（不指定则使用模型名的蛇形命名）
      without_create_time: false     # 设为 true 则不自动添加 created_at 字段
      without_update_time: false     # 设为 true 则不自动添加 updated_at 字段
      soft_delete: false             # 设为 true 则启用软删除（自动添加 deleted_at 字段）
      fields:
        - name: id                   # 字段名（蛇形命名）
          type: int64                # Go 类型
          tag: 'json:"id"'          # 结构体 Tag
        - name: name
          type: string
          tag: 'json:"name"'
```

### 支持的字段类型

你可以使用所有 Go 基本类型，以及通过 `imports` 引入的外部类型：

- `int`、`int64`、`int32` 等整数类型
- `float32`、`float64`
- `string`
- `time.Time`
- 自定义类型（需在 `imports` 中声明包路径）

> 不指定 `type` 时，默认为 `string`。

### 自动字段

除非显式关闭，Eloquent 会自动为每个模型添加以下字段：

| 字段 | 类型 | 说明 | 关闭方式 |
|------|------|------|----------|
| `created_at` | `time.Time` | 创建时间，插入时自动设置 | `without_create_time: true` |
| `updated_at` | `time.Time` | 更新时间，更新时自动设置 | `without_update_time: true` |
| `deleted_at` | `time.Time` | 删除时间，软删除时自动设置 | 仅在 `soft_delete: true` 时添加 |

### 代码生成命令

```bash
# 生成模型代码
eloquent gen --source './models/*.yaml'

# 创建模型定义文件模板
eloquent create-model \
  --table users \
  --package models \
  --output ./models/ \
  --table-prefix wz_ \
  --soft-delete \
  --no-created_at \
  --no-updated_at \
  --import "github.com/mylxsw/eloquent"
```

### 生成代码包含什么

对于一个 `User` 模型，生成的代码中包含：

| 内容 | 说明 |
|------|------|
| `UserN` 结构体 | Nullable 版本，字段使用 `null.Int`、`null.String` 等类型 |
| `User` 结构体 | Plain 版本，字段使用 Go 基本类型 |
| `UserModel` | 模型操作对象，提供所有 CRUD 方法 |
| `NewUserModel(db)` | 创建模型实例的工厂函数 |
| `FieldUser*` 常量 | 每个字段名的常量定义，如 `FieldUserName = "name"` |
| `UserFields()` | 返回所有字段名的切片 |
| `UserTable()` | 返回表名 |
| 关联方法 | 如果定义了 `relations`，会生成对应的关联方法 |

## 数据库方言 (MySQL / PostgreSQL)

Eloquent 通过方言（Dialect）系统同时支持 MySQL 和 PostgreSQL。默认使用 MySQL 方言。要切换到 PostgreSQL，在应用启动时设置方言：

```go
import (
    "github.com/mylxsw/eloquent/query"
    "github.com/mylxsw/eloquent/migrate"
)

// 切换到 PostgreSQL 模式（在应用启动时调用一次，必须在任何数据库操作之前）
query.SetDialect(query.PostgreSQLDialect{})
migrate.SetMigrateDialect(migrate.PostgreSQLMigrateDialect{})
```

方言系统会自动处理 MySQL 和 PostgreSQL 之间的差异：

| 特性 | MySQL | PostgreSQL |
|------|-------|------------|
| 标识符引用 | `` `name` `` | `"name"` |
| 参数占位符 | `?` | `$1, $2, ...` |
| 自增主键 | `AUTO_INCREMENT PRIMARY KEY` | `SERIAL PRIMARY KEY` |
| 插入返回 ID | `LastInsertId()` | `INSERT ... RETURNING id` |
| 引擎 / 字符集 / 排序规则 | 支持 | 忽略（不适用） |
| `TINYINT(1)` | `TINYINT(1)` | `BOOLEAN` |
| `BLOB` | `BLOB` | `BYTEA` |
| `LONGTEXT` / `MEDIUMTEXT` | 原样 | `TEXT` |
| `JSON` | `JSON` | `JSONB` |
| `DATETIME` | `DATETIME` | `TIMESTAMP` |
| `UNSIGNED` | 支持 | 忽略（不适用） |
| 创建索引 | `ALTER TABLE ... ADD INDEX` | `CREATE INDEX ... ON` |
| 删除索引 | `ALTER TABLE ... DROP INDEX` | `DROP INDEX` |
| 删除外键 | `DROP FOREIGN KEY` | `DROP CONSTRAINT` |

> **注意：** 使用原生 SQL 迁移（`Schema.Raw()`）时，你需要自行编写对应数据库的 SQL 语句。

## 数据库连接

Eloquent 直接使用 Go 标准库的 `database/sql` 来管理数据库连接。

### MySQL

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/mydb?parseTime=true")
if err != nil {
    panic(err)
}
defer db.Close()
```

> **重要：** 连接字符串中必须加上 `parseTime=true` 参数，否则时间字段无法正确解析。

### PostgreSQL

```go
import (
    "database/sql"

    _ "github.com/lib/pq"
    "github.com/mylxsw/eloquent/query"
    "github.com/mylxsw/eloquent/migrate"
)

// 在打开连接之前设置方言
query.SetDialect(query.PostgreSQLDialect{})
migrate.SetMigrateDialect(migrate.PostgreSQLMigrateDialect{})

db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=secret dbname=mydb sslmode=disable")
if err != nil {
    panic(err)
}
defer db.Close()
```

> **提示：** 你也可以使用 `pgx` 作为 PostgreSQL 驱动，导入 `github.com/jackc/pgx/v5/stdlib` 替代 `github.com/lib/pq`。

## CRUD 操作

所有操作都需要先创建模型实例：

```go
userModel := models.NewUserModel(db)
ctx := context.TODO()
```

### 创建 (Create)

```go
// 方式一：使用 Save 方法（传入 Nullable 结构体）
id, err := userModel.Save(ctx, models.UserN{
    Name:  null.StringFrom("张三"),
    Email: null.StringFrom("zhangsan@example.com"),
})

// 方式二：使用 Create 方法（传入 KV 键值对）
id, err := userModel.Create(ctx, query.KV{
    models.FieldUserName:  "张三",
    models.FieldUserEmail: "zhangsan@example.com",
})

// 方式三：从 Plain 结构体转换
user := models.User{Name: "张三", Email: "zhangsan@example.com"}
id, err := userModel.Save(ctx, user.ToUserN())
```

### 查询 (Read)

```go
// 根据 ID 查询
user, err := userModel.Find(ctx, 1)

// 查询第一条记录
user, err := userModel.First(ctx, query.Builder().Where("name", "张三"))

// 查询多条记录
users, err := userModel.Get(ctx)

// 带条件查询
users, err := userModel.Get(ctx, query.Builder().
    Where("age", ">", 18).
    OrderBy("created_at", "DESC").
    Limit(10),
)

// 查询指定字段
users, err := userModel.Get(ctx, query.Builder().Select("id", "name"))

// 分页查询
users, meta, err := userModel.Paginate(ctx, 1, 15)  // 第1页，每页15条
// meta.Total    - 总记录数
// meta.LastPage - 最后一页
// meta.Page     - 当前页
// meta.PerPage  - 每页条数

// 统计数量
count, err := userModel.Count(ctx)

// 判断是否存在
exists, err := userModel.Exists(ctx, query.Builder().Where("email", "zhangsan@example.com"))
```

### 读取字段值

```go
user, _ := userModel.Find(ctx, 1)

// 方式一：直接从 Nullable 结构体读取
name := user.Name.String           // 直接读取，如果是 NULL 则为空字符串
id := user.Id.Int64                // 直接读取 int64 值
valid := user.Name.Valid           // 判断字段是否为 NULL

// 方式二：转换为 Plain 结构体
u := user.ToUser()
fmt.Println(u.Name)    // string 类型
fmt.Println(u.Age)     // int64 类型

// 方式三：转换为自定义结构体（按字段名匹配）
type UserView struct {
    Name  string
    Email string
}
var view UserView
err := user.As(&view)
```

### 更新 (Update)

```go
// 方式一：通过模型方法更新指定字段
affected, err := userModel.UpdateFields(ctx,
    query.KV{models.FieldUserName: "李四"},
    query.Builder().Where("id", 1),
)

// 方式二：通过模型实例 Save（会更新所有被修改的字段）
user, _ := userModel.Find(ctx, 1)
user.Name = null.StringFrom("李四")
err := user.Save(ctx)  // 需要先 SetModel 或通过 Find 获取
```

### 删除 (Delete)

```go
// 根据 ID 删除（如果启用了软删除，则执行软删除）
affected, err := userModel.DeleteById(ctx, 1)

// 根据条件删除
affected, err := userModel.Delete(ctx, query.Builder().Where("age", "<", 18))

// 强制删除（即使启用了软删除，也从数据库中物理删除）
affected, err := userModel.ForceDeleteById(ctx, 1)
```

## 查询构建器

查询构建器是 Eloquent 中构建 SQL 查询条件的核心工具。所有条件都通过 `query.Builder()` 开始，支持链式调用。

### 基本条件

```go
// 等于（默认操作符）
query.Builder().Where("name", "张三")

// 指定操作符
query.Builder().Where("age", ">", 18)
query.Builder().Where("name", "LIKE", "%张%")

// 可用的操作符常量
query.EQ   // =
query.NEQ  // !=
query.GT   // >
query.GTE  // >=
query.LT   // <
query.LTE  // <=
query.LIKE // LIKE
```

### OR 条件

```go
query.Builder().Where("name", "张三").OrWhere("name", "李四")
// WHERE name = '张三' OR name = '李四'
```

### IN / NOT IN

```go
query.Builder().WhereIn("id", 1, 2, 3)
query.Builder().WhereNotIn("status", 0, -1)

// 使用切片时需要转换
ids := []int{1, 2, 3}
query.Builder().WhereIn("id", query.ToAnys(ids)...)
```

### NULL 判断

```go
query.Builder().WhereNull("deleted_at")
query.Builder().WhereNotNull("email")
```

### BETWEEN

```go
query.Builder().WhereBetween("age", 18, 30)
query.Builder().WhereNotBetween("age", 0, 17)
```

### 条件分组（括号）

```go
// WHERE (name = '张三' AND age > 18) OR (name = '李四')
query.Builder().
    WhereGroup(func(builder query.Condition) {
        builder.Where("name", "张三").Where("age", ">", 18)
    }).
    OrWhereGroup(func(builder query.Condition) {
        builder.Where("name", "李四")
    })
```

### 条件查询 (When)

根据运行时条件动态添加查询条件：

```go
keyword := "张三"
query.Builder().When(
    func() bool { return keyword != "" },
    func(builder query.Condition) {
        builder.Where("name", "LIKE", "%"+keyword+"%")
    },
)
```

### 原生表达式

```go
query.Builder().WhereRaw("YEAR(created_at) = ?", 2024)
```

### 子查询

```go
// WHERE id IN (SELECT user_id FROM orders WHERE amount > 100)
query.Builder().WhereIn("id",
    query.Builder().Table("orders").Select("user_id").Where("amount", ">", 100),
)

// WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)
query.Builder().WhereExist(
    query.Builder().Table("orders").Select(query.Raw("1")).WhereColumn("orders.user_id", "=", "users.id"),
)
```

### 排序、分页、分组

```go
query.Builder().
    OrderBy("created_at", "DESC").
    OrderBy("id", "ASC").
    Limit(10).
    Offset(20).
    GroupBy("status").
    Having(func(builder query.Condition) {
        builder.WhereRaw("COUNT(*) > ?", 5)
    })
```

### JOIN 查询

```go
query.Builder().
    Table("users AS u").
    Select("u.name", "r.name AS role_name").
    LeftJoin("roles AS r", func(builder query.Condition) {
        builder.WhereColumn("u.role_id", "=", "r.id")
    })
```

支持的 JOIN 类型：`LeftJoin`、`RightJoin`、`InnerJoin`、`CrossJoin`

### UNION

```go
query.Builder().Table("table1").Select("name").
    Union(query.Builder().Table("table2").Select("name"), false)  // false=UNION ALL, true=UNION DISTINCT
```

### SELECT 原生表达式

```go
query.Builder().Select("id", "name", query.Raw("DATE(created_at) AS create_date"))
```

## 模型关联

Eloquent 支持四种关联关系，通过在 YAML 中定义 `relations` 来使用。

### 关联类型

| 关联类型 | YAML 中的 `rel` 值 | 说明 | 示例 |
|---------|-------------------|------|------|
| 一对一 | `hasOne` | 当前模型拥有一个子模型 | 用户有一个用户扩展信息 |
| 反向一对一/多对一 | `belongsTo` | 当前模型属于另一个模型 | 用户属于一个角色 |
| 多对多 | `belongsToMany` | 通过中间表关联 | 用户属于多个组织 |

### 定义关联

```yaml
models:
  - name: User
    relations:
      # 多对一：用户属于一个角色
      - model: role
        rel: belongsTo
        foreign_key: role_id     # 当前表中的外键字段
        owner_key: id            # 关联表的主键

      # 一对一：用户有一个扩展信息
      - model: userExt
        rel: hasOne
        # 默认外键为 user_id（当前模型名_id）

      # 多对多：用户属于多个组织
      - model: organization
        rel: belongsToMany
        # 默认中间表名为 user_organization_ref
    definition:
      fields:
        - name: id
          type: int64
        - name: role_id
          type: int64
```

### 使用关联

```go
user, _ := userModel.Find(ctx, 1)

// belongsTo - 查询用户所属角色
role, err := user.Role().First(ctx)

// belongsTo - 创建关联对象并自动绑定外键
roleId, err := user.Role().Create(ctx, models.RoleN{
    Name: null.StringFrom("管理员"),
})

// belongsTo - 关联/解除关联
err = user.Role().Associate(ctx, role)     // 绑定
err = user.Role().Dissociate(ctx)          // 解绑

// hasOne - 查询用户扩展信息
ext, err := user.UserExt().First(ctx)
exists, err := user.UserExt().Exists(ctx)

// belongsToMany - 查询用户的所有组织
orgs, err := user.Organizations().Get(ctx)
count, err := user.Organizations().Count(ctx)

// belongsToMany - 附加/分离
err = user.Organizations().Attach(ctx, org)      // 在中间表插入记录
err = user.Organizations().Detach(ctx, org)      // 从中间表删除记录
err = user.Organizations().DetachAll(ctx)         // 删除所有关联
orgId, err := user.Organizations().Create(ctx, models.OrganizationN{
    Name: null.StringFrom("研发部"),
})
```

## 软删除

启用软删除后，调用 `Delete` 不会真正删除数据，而是设置 `deleted_at` 字段。查询时会自动过滤已软删除的记录。

### 启用软删除

在 YAML 中设置 `soft_delete: true`：

```yaml
models:
  - name: User
    definition:
      soft_delete: true
      fields:
        - name: id
          type: int64
```

### 使用

```go
// 普通删除（软删除，设置 deleted_at）
userModel.DeleteById(ctx, 1)

// 普通查询（自动排除已删除的记录）
users, _ := userModel.Get(ctx)

// 查询包含已删除的记录
users, _ := userModel.WithTrashed().Get(ctx)

// 强制删除（物理删除，从数据库中彻底移除）
userModel.ForceDeleteById(ctx, 1)
```

## 作用域

作用域 (Scope) 是一种将常用查询条件封装为可复用函数的机制。

### 全局作用域 (Global Scope)

全局作用域会自动应用到模型的所有查询中。软删除就是通过全局作用域实现的。

```go
// 注册全局作用域（通常在 init 函数中）
models.AddGlobalScopeForUser("active_only", func(builder query.Condition) {
    builder.Where("status", 1)
})

// 查询时自动应用
users, _ := userModel.Get(ctx) // 自动加上 WHERE status = 1

// 临时禁用某个全局作用域
users, _ := userModel.WithoutGlobalScopes("active_only").Get(ctx)
```

### 本地作用域 (Local Scope)

本地作用域不会自动应用，需要显式启用。

```go
// 注册本地作用域
models.AddLocalScopeForUser("recent", func(builder query.Condition) {
    builder.WhereRaw("created_at > DATE_SUB(NOW(), INTERVAL 7 DAY)")
})

// 显式启用
users, _ := userModel.WithLocalScopes("recent").Get(ctx)
```

## 数据库迁移

Eloquent 提供了与 Laravel 类似的数据库迁移功能，通过 Go 代码来管理数据库表结构的变更。

### 基本用法

```go
import "github.com/mylxsw/eloquent/migrate"

m := migrate.NewManager(db).Init(context.TODO())

// 创建新表
m.Schema("20240101_001").Create("users", func(builder *migrate.Builder) {
    builder.Increments("id")
    builder.String("name", 100).Comment("用户名")
    builder.String("email", 255).Unique().Comment("邮箱")
    builder.String("password", 255)
    builder.Timestamps(0)
})

// 修改现有表（添加字段）
m.Schema("20240101_002").Table("users", func(builder *migrate.Builder) {
    builder.TinyInteger("status", false, true).
        Default(migrate.StringExpr("1")).
        Comment("状态：0-禁用 1-启用")
    builder.SoftDeletes("deleted_at", 0)
})

// 执行迁移
if err := m.Run(context.TODO()); err != nil {
    panic(err)
}
```

> 每个 `Schema` 的 version 字符串是唯一标识，已经执行过的迁移不会重复执行。Eloquent 会自动创建 `migrations` 表来记录迁移历史。

### 可用的字段类型方法

| 方法 | MySQL 类型 | PostgreSQL 类型 |
|------|-----------|-----------------|
| `Increments(name)` | INT UNSIGNED AUTO_INCREMENT | SERIAL |
| `BigIncrements(name)` | BIGINT UNSIGNED AUTO_INCREMENT | BIGSERIAL |
| `Integer(name, autoIncrement, unsigned)` | INT | INTEGER |
| `BigInteger(name, autoIncrement, unsigned)` | BIGINT | BIGINT |
| `TinyInteger(name, autoIncrement, unsigned)` | TINYINT | SMALLINT |
| `SmallInteger(name, autoIncrement, unsigned)` | SMALLINT | SMALLINT |
| `MediumInteger(name, autoIncrement, unsigned)` | MEDIUMINT | INTEGER |
| `Float(name, total, scale)` | DOUBLE | DOUBLE PRECISION |
| `Double(name, total, scale)` | DOUBLE | DOUBLE PRECISION |
| `Decimal(name, total, scale)` | DECIMAL(m,n) | DECIMAL(m,n) |
| `String(name, length)` | VARCHAR(n) | VARCHAR(n) |
| `Char(name, length)` | CHAR(n) | CHAR(n) |
| `Text(name)` | TEXT | TEXT |
| `MediumText(name)` | MEDIUMTEXT | TEXT |
| `LongText(name)` | LONGTEXT | TEXT |
| `Boolean(name)` | TINYINT(1) | BOOLEAN |
| `Date(name)` | DATE | DATE |
| `DateTime(name, precision)` | DATETIME | TIMESTAMP |
| `Timestamp(name, precision)` | TIMESTAMP | TIMESTAMP |
| `Time(name, precision)` | TIME | TIME |
| `Year(name)` | YEAR | SMALLINT |
| `Binary(name)` | BLOB | BYTEA |
| `Json(name)` | JSON | JSONB |
| `Enum(name, items...)` | ENUM | TEXT |
| `Set(name, items...)` | SET | TEXT |
| `Uuid(name)` | CHAR(36) | CHAR(36) |
| `IpAddress(name)` | VARCHAR(45) | VARCHAR(45) |
| `MacAddress(name)` | VARCHAR(17) | VARCHAR(17) |
| `Timestamps(precision)` | created_at + updated_at | created_at + updated_at |
| `SoftDeletes(column, precision)` | 可空 TIMESTAMP | 可空 TIMESTAMP |
| `RememberToken()` | VARCHAR(100) | VARCHAR(100) |

### 字段修饰方法

```go
builder.String("name", 100).
    Nullable(true).          // 允许 NULL
    Default(migrate.StringExpr("default_value")). // 默认值（字符串）
    Default(migrate.RawExpr("CURRENT_TIMESTAMP")). // 默认值（原生表达式）
    Comment("用户名").        // 字段注释（仅 MySQL）
    Unique().                // 唯一索引
    After("id").             // 放在某字段之后（仅 MySQL）
    First().                 // 放在第一列（仅 MySQL）
    Charset("utf8mb4").      // 字符集（仅 MySQL）
    Collation("utf8mb4_unicode_ci"). // 排序规则（仅 MySQL）
    Unsigned().              // 无符号（仅 MySQL，PostgreSQL 下忽略）
    Change()                 // 修改已有字段（而非添加）
```

> **PostgreSQL 注意：** 标注"仅 MySQL"的修饰方法在 PostgreSQL 方言下会被静默忽略。

### 索引操作

```go
builder.Index("", "name", "email")         // 创建普通索引（名称自动生成）
builder.Index("idx_name", "name")          // 创建指定名称的索引
builder.Unique("", "email")                // 创建唯一索引
builder.Primary("", "id")                  // 创建主键索引

builder.DropIndex("idx_name")              // 删除索引
builder.DropUnique("unique_name")          // 删除唯一索引
builder.DropColumn("column1", "column2")   // 删除字段
```

### 外键

```go
builder.Foreign("", "user_id").
    References("id").
    On("users").
    OnDelete("CASCADE").
    OnUpdate("CASCADE")

builder.DropForeign("foreign_key_name")    // 删除外键
```

### 表操作

```go
// 删除表
m.Schema("20240102_001").Drop("temp_table")
m.Schema("20240102_002").DropIfExists("temp_table")

// 原生 SQL 迁移
m.Schema("20240102_003").Raw("custom_table", func() []string {
    return []string{
        `CREATE TABLE custom_table (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255))`,
    }
})
```

## 事务

```go
import "github.com/mylxsw/eloquent/query"

err := eloquent.Transaction(db, func(tx query.Database) error {
    userModel := models.NewUserModel(tx)  // 使用 tx 而不是 db

    id, err := userModel.Save(ctx, models.UserN{
        Name: null.StringFrom("张三"),
    })
    if err != nil {
        return err  // 返回 error 会自动回滚
    }

    // ... 更多操作

    return nil  // 返回 nil 会自动提交
})
```

> 事务中发生 panic 也会自动回滚。

## 事件系统

Eloquent 内置了一个事件系统，可以监听 SQL 执行、事务和迁移等事件。

### 设置事件监听

```go
import "github.com/mylxsw/eloquent/event"

em := event.NewEventManager(event.NewMemoryEventStore())
event.SetDispatcher(em)

// 监听 SQL 执行事件（适合日志记录、性能监控）
em.Listen(func(evt event.QueryExecutedEvent) {
    fmt.Printf("[SQL] %s | args: %v | time: %s\n", evt.SQL, evt.Bindings, evt.Time)
})

// 监听事务事件
em.Listen(func(evt event.TransactionBeginningEvent) {
    fmt.Println("事务开始")
})
em.Listen(func(evt event.TransactionCommittedEvent) {
    fmt.Println("事务已提交")
})
em.Listen(func(evt event.TransactionRolledBackEvent) {
    fmt.Println("事务已回滚")
})

// 监听迁移事件
em.Listen(func(evt event.MigrationStartedEvent) {
    fmt.Printf("执行迁移: %s\n", evt.SQL)
})
```

### 可用事件

| 事件 | 说明 |
|------|------|
| `QueryExecutedEvent` | 每条 SQL 执行后触发，包含 SQL、参数和耗时 |
| `TransactionBeginningEvent` | 事务开始时触发 |
| `TransactionCommittedEvent` | 事务提交时触发 |
| `TransactionRolledBackEvent` | 事务回滚时触发 |
| `MigrationsStartedEvent` | 一批迁移开始前触发 |
| `MigrationsEndedEvent` | 一批迁移完成后触发 |
| `MigrationStartedEvent` | 单条迁移 SQL 执行前触发 |
| `MigrationEndedEvent` | 单条迁移 SQL 执行后触发 |

## 直接使用数据库操作

除了通过模型操作外，Eloquent 还提供了低层级的数据库操作 API，适合执行复杂查询或不想定义模型的场景。

```go
import "github.com/mylxsw/eloquent"

dbOp := eloquent.DB(db)
ctx := context.TODO()

// 插入
id, err := dbOp.Insert(ctx, "users", query.KV{
    "name":  "张三",
    "email": "zhangsan@example.com",
})

// 查询
results, err := dbOp.Query(ctx,
    eloquent.Build("users").Select("id", "name").Where("age", ">", 18),
    func(row eloquent.Scanner) (any, error) {
        var id int64
        var name string
        err := row.Scan(&id, &name)
        return map[string]any{"id": id, "name": name}, err
    },
)

// 使用泛型查询（Go 1.18+）
type UserRow struct {
    Id   int64
    Name string
}
users, err := eloquent.Query[UserRow](ctx, db,
    eloquent.Build("users").Select("id", "name"),
    func(row eloquent.Scanner) (UserRow, error) {
        var u UserRow
        err := row.Scan(&u.Id, &u.Name)
        return u, err
    },
)

// 原生 SQL 查询
results, err := dbOp.Query(ctx,
    eloquent.Raw("SELECT * FROM users WHERE age > ?", 18),
    func(row eloquent.Scanner) (any, error) {
        // ...
    },
)

// 更新
affected, err := dbOp.Update(ctx,
    eloquent.Build("users").Where("id", 1),
    query.KV{"name": "李四"},
)

// 删除
affected, err := dbOp.Delete(ctx, eloquent.Build("users").Where("id", 1))

// 执行无返回值的语句
err := dbOp.Statement(ctx, "TRUNCATE TABLE users")
```

## 示例项目

- [_examples](./_examples) - 项目内置的完整示例代码
- [tech-share](https://github.com/mylxsw/tech-share) - 一个简单的 Web 项目，展示了 Eloquent 在实际项目中的使用

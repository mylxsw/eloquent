# Eloquent ORM for Go

> English | **[中文文档](./README_CN.md)**

Eloquent is a **code-generation based** database ORM framework for Go, inspired by the Eloquent ORM in the PHP framework Laravel. Supports MySQL and PostgreSQL.

## Table of Contents

- [Core Concepts](#core-concepts)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Model Definition & Code Generation](#model-definition--code-generation)
- [Database Dialect (MySQL / PostgreSQL)](#database-dialect-mysql--postgresql)
- [Database Connection](#database-connection)
- [CRUD Operations](#crud-operations)
- [Query Builder](#query-builder)
- [Model Relationships](#model-relationships)
- [Soft Deletes](#soft-deletes)
- [Scopes](#scopes)
- [Database Migrations](#database-migrations)
- [Transactions](#transactions)
- [Event System](#event-system)
- [Direct Database Operations](#direct-database-operations)
- [Example Projects](#example-projects)

## Core Concepts

Understanding these core concepts will help you get started with Eloquent quickly:

### 1. Code Generation Driven

Eloquent does not use runtime reflection to map database tables. Instead, it works through a **YAML model definition -> CLI tool generates Go code** workflow:
- Write a YAML file describing your database table structure
- Run the `eloquent gen` command to generate the corresponding Go code (`*.orm.go` files)
- The generated code contains type-safe model structs and complete CRUD methods

### 2. Two Model Structs

For each model (e.g., `User`), Eloquent generates two structs:

- **`UserN`** (Nullable version): All fields use `null.Int`, `null.String`, etc. (from `gopkg.in/guregu/null.v3`), correctly handling `NULL` values from the database. This is the primary struct used for database interaction.
- **`User`** (Plain version): All fields use Go primitive types (`int64`, `string`, etc.). When a database field is `NULL`, it returns the zero value. Suitable for use in business logic.

They can be converted to each other:

```go
// UserN -> User (Nullable -> Plain)
user := userN.ToUser()

// User -> UserN (Plain -> Nullable)
userN := user.ToUserN()
```

### 3. Query Builder (SQLBuilder)

All query conditions are built through `query.Builder()` chained calls, rather than writing raw SQL. The query builder supports `Where`, `WhereIn`, `OrderBy`, `Limit`, and other common operations.

### 4. Context

All database operations require a `context.Context` parameter, following Go best practices for timeout control and cancellation.

## Quick Start

Here is a complete example showing the full workflow from model definition to database operations.

### Step 1: Install the CLI Tool

```bash
go install github.com/mylxsw/eloquent/cmd/orm@latest
```

This gives you an `orm` CLI tool (you can also build the project yourself to get an `eloquent` binary).

### Step 2: Create a Model Definition File

You can write YAML manually or use the CLI to generate a basic template:

```bash
eloquent create-model --table users --package models --soft-delete --output ./models/
```

Then edit the generated `users.yaml` file to add the fields you need:

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

### Step 3: Generate Model Code

```bash
eloquent gen --source './models/*.yaml'
```

This generates a `users.orm.go` file in the same directory, containing `UserN`, `User`, `UserModel`, and all related code.

> **Note:** Generated `*.orm.go` files should not be manually modified, as they are regenerated each time the `gen` command is run.

### Step 4: Add Dependencies

```bash
go get github.com/mylxsw/eloquent
go get gopkg.in/guregu/null.v3

# For MySQL
go get github.com/go-sql-driver/mysql

# For PostgreSQL (use one of the following)
go get github.com/lib/pq
# or: go get github.com/jackc/pgx/v5
```

### Step 5: Write Business Code

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

    // Create a user
    id, err := userModel.Save(ctx, models.UserN{
        Name:  null.StringFrom("Tom"),
        Email: null.StringFrom("tom@example.com"),
        Age:   null.IntFrom(25),
    })
    fmt.Printf("New user ID: %d\n", id)

    // Query a user
    user, err := userModel.Find(ctx, id)
    fmt.Printf("Name: %s, Email: %s\n", user.Name.String, user.Email.String)

    // Convert to Plain struct
    u := user.ToUser()
    fmt.Printf("Name: %s, Age: %d\n", u.Name, u.Age)
}
```

## Installation

Add the Eloquent dependency to your project:

```bash
go get github.com/mylxsw/eloquent
```

Build the CLI tool (for code generation):

```bash
# Option 1: Direct install
go install github.com/mylxsw/eloquent/cmd/orm@latest

# Option 2: Clone and build
git clone https://github.com/mylxsw/eloquent.git
cd eloquent
make build
# Binary is generated at ./bin/eloquent
```

## Model Definition & Code Generation

### YAML Model Definition File

Model definitions use YAML format. Here is a complete field reference:

```yaml
package: models                      # Go package name for generated code
imports:                             # Additional packages to import (when fields use external types)
  - github.com/mylxsw/eloquent

meta:
  table_prefix: wz_                  # Table prefix, automatically prepended to table names

models:
  - name: User                       # Model name (generated struct name)
    definition:
      table_name: users              # Database table name (defaults to snake_case of model name)
      without_create_time: false     # Set to true to skip adding created_at field
      without_update_time: false     # Set to true to skip adding updated_at field
      soft_delete: false             # Set to true to enable soft deletes (adds deleted_at field)
      fields:
        - name: id                   # Field name (snake_case)
          type: int64                # Go type
          tag: 'json:"id"'          # Struct tag
        - name: name
          type: string
          tag: 'json:"name"'
```

### Supported Field Types

You can use all Go primitive types, as well as external types imported via `imports`:

- `int`, `int64`, `int32`, and other integer types
- `float32`, `float64`
- `string`
- `time.Time`
- Custom types (declare the package path in `imports`)

> When `type` is not specified, it defaults to `string`.

### Automatic Fields

Unless explicitly disabled, Eloquent automatically adds the following fields to each model:

| Field | Type | Description | Disable With |
|-------|------|-------------|--------------|
| `created_at` | `time.Time` | Creation time, auto-set on insert | `without_create_time: true` |
| `updated_at` | `time.Time` | Update time, auto-set on update | `without_update_time: true` |
| `deleted_at` | `time.Time` | Deletion time, auto-set on soft delete | Only added when `soft_delete: true` |

### Code Generation Commands

```bash
# Generate model code
eloquent gen --source './models/*.yaml'

# Create a model definition file template
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

### What Gets Generated

For a `User` model, the generated code includes:

| Content | Description |
|---------|-------------|
| `UserN` struct | Nullable version, fields use `null.Int`, `null.String`, etc. |
| `User` struct | Plain version, fields use Go primitive types |
| `UserModel` | Model operation object, provides all CRUD methods |
| `NewUserModel(db)` | Factory function to create a model instance |
| `FieldUser*` constants | Constant definitions for each field name, e.g., `FieldUserName = "name"` |
| `UserFields()` | Returns a slice of all field names |
| `UserTable()` | Returns the table name |
| Relationship methods | Generated if `relations` are defined |

## Database Dialect (MySQL / PostgreSQL)

Eloquent supports both MySQL and PostgreSQL through a dialect system. By default, MySQL dialect is used. To switch to PostgreSQL, set the dialects at application startup:

```go
import (
    "github.com/mylxsw/eloquent/query"
    "github.com/mylxsw/eloquent/migrate"
)

// Switch to PostgreSQL mode (call once at startup, before any database operations)
query.SetDialect(query.PostgreSQLDialect{})
migrate.SetMigrateDialect(migrate.PostgreSQLMigrateDialect{})
```

The dialect system automatically handles the differences between MySQL and PostgreSQL:

| Feature | MySQL | PostgreSQL |
|---------|-------|------------|
| Identifier quoting | `` `name` `` | `"name"` |
| Parameter placeholders | `?` | `$1, $2, ...` |
| Auto-increment primary key | `AUTO_INCREMENT PRIMARY KEY` | `SERIAL PRIMARY KEY` |
| Insert returning ID | `LastInsertId()` | `INSERT ... RETURNING id` |
| Engine / Charset / Collation | Supported | Ignored (not applicable) |
| `TINYINT(1)` | `TINYINT(1)` | `BOOLEAN` |
| `BLOB` | `BLOB` | `BYTEA` |
| `LONGTEXT` / `MEDIUMTEXT` | As-is | `TEXT` |
| `JSON` | `JSON` | `JSONB` |
| `DATETIME` | `DATETIME` | `TIMESTAMP` |
| `UNSIGNED` | Supported | Ignored (not applicable) |
| Index creation | `ALTER TABLE ... ADD INDEX` | `CREATE INDEX ... ON` |
| Drop index | `ALTER TABLE ... DROP INDEX` | `DROP INDEX` |
| Drop foreign key | `DROP FOREIGN KEY` | `DROP CONSTRAINT` |

> **Note:** When using raw SQL in migrations (`Schema.Raw()`), you are responsible for writing database-specific SQL yourself.

## Database Connection

Eloquent uses Go's standard `database/sql` library to manage database connections.

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

> **Important:** You must include `parseTime=true` in the connection string, otherwise time fields cannot be parsed correctly.

### PostgreSQL

```go
import (
    "database/sql"

    _ "github.com/lib/pq"
    "github.com/mylxsw/eloquent/query"
    "github.com/mylxsw/eloquent/migrate"
)

// Set dialect before opening the connection
query.SetDialect(query.PostgreSQLDialect{})
migrate.SetMigrateDialect(migrate.PostgreSQLMigrateDialect{})

db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=secret dbname=mydb sslmode=disable")
if err != nil {
    panic(err)
}
defer db.Close()
```

> **Tip:** You can also use `pgx` as the PostgreSQL driver by importing `github.com/jackc/pgx/v5/stdlib` instead of `github.com/lib/pq`.

## CRUD Operations

All operations require creating a model instance first:

```go
userModel := models.NewUserModel(db)
ctx := context.TODO()
```

### Create

```go
// Method 1: Using Save with a Nullable struct
id, err := userModel.Save(ctx, models.UserN{
    Name:  null.StringFrom("Tom"),
    Email: null.StringFrom("tom@example.com"),
})

// Method 2: Using Create with KV pairs
id, err := userModel.Create(ctx, query.KV{
    models.FieldUserName:  "Tom",
    models.FieldUserEmail: "tom@example.com",
})

// Method 3: Converting from a Plain struct
user := models.User{Name: "Tom", Email: "tom@example.com"}
id, err := userModel.Save(ctx, user.ToUserN())
```

### Read

```go
// Find by ID
user, err := userModel.Find(ctx, 1)

// Find the first matching record
user, err := userModel.First(ctx, query.Builder().Where("name", "Tom"))

// Get all records
users, err := userModel.Get(ctx)

// Get with conditions
users, err := userModel.Get(ctx, query.Builder().
    Where("age", ">", 18).
    OrderBy("created_at", "DESC").
    Limit(10),
)

// Select specific fields
users, err := userModel.Get(ctx, query.Builder().Select("id", "name"))

// Pagination
users, meta, err := userModel.Paginate(ctx, 1, 15)  // Page 1, 15 per page
// meta.Total    - total record count
// meta.LastPage - last page number
// meta.Page     - current page
// meta.PerPage  - records per page

// Count
count, err := userModel.Count(ctx)

// Check existence
exists, err := userModel.Exists(ctx, query.Builder().Where("email", "tom@example.com"))
```

### Reading Field Values

```go
user, _ := userModel.Find(ctx, 1)

// Method 1: Read directly from Nullable struct
name := user.Name.String           // Direct read, empty string if NULL
id := user.Id.Int64                // Direct int64 read
valid := user.Name.Valid           // Check if field is NULL

// Method 2: Convert to Plain struct
u := user.ToUser()
fmt.Println(u.Name)    // string type
fmt.Println(u.Age)     // int64 type

// Method 3: Convert to a custom struct (matched by field name)
type UserView struct {
    Name  string
    Email string
}
var view UserView
err := user.As(&view)
```

### Update

```go
// Method 1: Update specific fields via model method
affected, err := userModel.UpdateFields(ctx,
    query.KV{models.FieldUserName: "Jerry"},
    query.Builder().Where("id", 1),
)

// Method 2: Save via model instance (updates all modified fields)
user, _ := userModel.Find(ctx, 1)
user.Name = null.StringFrom("Jerry")
err := user.Save(ctx)  // Requires SetModel or obtained via Find
```

### Delete

```go
// Delete by ID (soft delete if enabled)
affected, err := userModel.DeleteById(ctx, 1)

// Delete with conditions
affected, err := userModel.Delete(ctx, query.Builder().Where("age", "<", 18))

// Force delete (physical delete even if soft delete is enabled)
affected, err := userModel.ForceDeleteById(ctx, 1)
```

## Query Builder

The query builder is the core tool for constructing SQL query conditions in Eloquent. All conditions start with `query.Builder()` and support chaining.

### Basic Conditions

```go
// Equals (default operator)
query.Builder().Where("name", "Tom")

// Specify operator
query.Builder().Where("age", ">", 18)
query.Builder().Where("name", "LIKE", "%Tom%")

// Available operator constants
query.EQ   // =
query.NEQ  // !=
query.GT   // >
query.GTE  // >=
query.LT   // <
query.LTE  // <=
query.LIKE // LIKE
```

### OR Conditions

```go
query.Builder().Where("name", "Tom").OrWhere("name", "Jerry")
// WHERE name = 'Tom' OR name = 'Jerry'
```

### IN / NOT IN

```go
query.Builder().WhereIn("id", 1, 2, 3)
query.Builder().WhereNotIn("status", 0, -1)

// When using slices, convert first
ids := []int{1, 2, 3}
query.Builder().WhereIn("id", query.ToAnys(ids)...)
```

### NULL Checks

```go
query.Builder().WhereNull("deleted_at")
query.Builder().WhereNotNull("email")
```

### BETWEEN

```go
query.Builder().WhereBetween("age", 18, 30)
query.Builder().WhereNotBetween("age", 0, 17)
```

### Condition Grouping (Parentheses)

```go
// WHERE (name = 'Tom' AND age > 18) OR (name = 'Jerry')
query.Builder().
    WhereGroup(func(builder query.Condition) {
        builder.Where("name", "Tom").Where("age", ">", 18)
    }).
    OrWhereGroup(func(builder query.Condition) {
        builder.Where("name", "Jerry")
    })
```

### Conditional Clauses (When)

Dynamically add query conditions based on runtime values:

```go
keyword := "Tom"
query.Builder().When(
    func() bool { return keyword != "" },
    func(builder query.Condition) {
        builder.Where("name", "LIKE", "%"+keyword+"%")
    },
)
```

### Raw Expressions

```go
query.Builder().WhereRaw("YEAR(created_at) = ?", 2024)
```

### Subqueries

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

### Ordering, Pagination, Grouping

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

### JOIN Queries

```go
query.Builder().
    Table("users AS u").
    Select("u.name", "r.name AS role_name").
    LeftJoin("roles AS r", func(builder query.Condition) {
        builder.WhereColumn("u.role_id", "=", "r.id")
    })
```

Supported JOIN types: `LeftJoin`, `RightJoin`, `InnerJoin`, `CrossJoin`

### UNION

```go
query.Builder().Table("table1").Select("name").
    Union(query.Builder().Table("table2").Select("name"), false)  // false=UNION ALL, true=UNION DISTINCT
```

### SELECT Raw Expressions

```go
query.Builder().Select("id", "name", query.Raw("DATE(created_at) AS create_date"))
```

## Model Relationships

Eloquent supports relationship types defined via `relations` in the YAML file.

### Relationship Types

| Relationship | `rel` Value in YAML | Description | Example |
|-------------|---------------------|-------------|---------|
| One-to-One | `hasOne` | Current model owns one child model | User has one profile |
| Belongs To | `belongsTo` | Current model belongs to another model | User belongs to a role |
| Many-to-Many | `belongsToMany` | Related through a pivot table | User belongs to many organizations |

### Defining Relationships

```yaml
models:
  - name: User
    relations:
      # Belongs To: User belongs to a role
      - model: role
        rel: belongsTo
        foreign_key: role_id     # Foreign key in current table
        owner_key: id            # Primary key in related table

      # Has One: User has one extension record
      - model: userExt
        rel: hasOne
        # Default foreign key is user_id (current_model_name_id)

      # Many-to-Many: User belongs to many organizations
      - model: organization
        rel: belongsToMany
        # Default pivot table is user_organization_ref
    definition:
      fields:
        - name: id
          type: int64
        - name: role_id
          type: int64
```

### Using Relationships

```go
user, _ := userModel.Find(ctx, 1)

// belongsTo - Query the user's role
role, err := user.Role().First(ctx)

// belongsTo - Create a related object and auto-bind the foreign key
roleId, err := user.Role().Create(ctx, models.RoleN{
    Name: null.StringFrom("admin"),
})

// belongsTo - Associate / Dissociate
err = user.Role().Associate(ctx, role)     // Bind
err = user.Role().Dissociate(ctx)          // Unbind

// hasOne - Query user extension
ext, err := user.UserExt().First(ctx)
exists, err := user.UserExt().Exists(ctx)

// belongsToMany - Query all organizations for the user
orgs, err := user.Organizations().Get(ctx)
count, err := user.Organizations().Count(ctx)

// belongsToMany - Attach / Detach
err = user.Organizations().Attach(ctx, org)      // Insert into pivot table
err = user.Organizations().Detach(ctx, org)      // Remove from pivot table
err = user.Organizations().DetachAll(ctx)         // Remove all associations
orgId, err := user.Organizations().Create(ctx, models.OrganizationN{
    Name: null.StringFrom("Engineering"),
})
```

## Soft Deletes

When soft deletes are enabled, calling `Delete` does not physically remove data -- it sets the `deleted_at` field instead. Queries automatically filter out soft-deleted records.

### Enable Soft Deletes

Set `soft_delete: true` in the YAML:

```yaml
models:
  - name: User
    definition:
      soft_delete: true
      fields:
        - name: id
          type: int64
```

### Usage

```go
// Soft delete (sets deleted_at)
userModel.DeleteById(ctx, 1)

// Normal query (automatically excludes deleted records)
users, _ := userModel.Get(ctx)

// Query including deleted records
users, _ := userModel.WithTrashed().Get(ctx)

// Force delete (physically removes from database)
userModel.ForceDeleteById(ctx, 1)
```

## Scopes

Scopes encapsulate commonly used query conditions into reusable functions.

### Global Scopes

Global scopes are automatically applied to all queries on a model. Soft delete is implemented via a global scope.

```go
// Register a global scope (typically in an init function)
models.AddGlobalScopeForUser("active_only", func(builder query.Condition) {
    builder.Where("status", 1)
})

// Automatically applied during queries
users, _ := userModel.Get(ctx) // Automatically adds WHERE status = 1

// Temporarily disable a global scope
users, _ := userModel.WithoutGlobalScopes("active_only").Get(ctx)
```

### Local Scopes

Local scopes are not automatically applied; they must be explicitly enabled.

```go
// Register a local scope
models.AddLocalScopeForUser("recent", func(builder query.Condition) {
    builder.WhereRaw("created_at > DATE_SUB(NOW(), INTERVAL 7 DAY)")
})

// Explicitly enable
users, _ := userModel.WithLocalScopes("recent").Get(ctx)
```

## Database Migrations

Eloquent provides Laravel-like database migration functionality, managing table schema changes through Go code.

### Basic Usage

```go
import "github.com/mylxsw/eloquent/migrate"

m := migrate.NewManager(db).Init(context.TODO())

// Create a new table
m.Schema("20240101_001").Create("users", func(builder *migrate.Builder) {
    builder.Increments("id")
    builder.String("name", 100).Comment("Username")
    builder.String("email", 255).Unique().Comment("Email")
    builder.String("password", 255)
    builder.Timestamps(0)
})

// Modify an existing table (add columns)
m.Schema("20240101_002").Table("users", func(builder *migrate.Builder) {
    builder.TinyInteger("status", false, true).
        Default(migrate.StringExpr("1")).
        Comment("Status: 0-disabled 1-enabled")
    builder.SoftDeletes("deleted_at", 0)
})

// Run migrations
if err := m.Run(context.TODO()); err != nil {
    panic(err)
}
```

> Each `Schema` version string is a unique identifier. Already-executed migrations will not be re-run. Eloquent automatically creates a `migrations` table to track migration history.

### Available Column Type Methods

| Method | MySQL Type | PostgreSQL Type |
|--------|-----------|-----------------|
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
| `SoftDeletes(column, precision)` | Nullable TIMESTAMP | Nullable TIMESTAMP |
| `RememberToken()` | VARCHAR(100) | VARCHAR(100) |

### Column Modifiers

```go
builder.String("name", 100).
    Nullable(true).          // Allow NULL
    Default(migrate.StringExpr("default_value")). // Default value (string)
    Default(migrate.RawExpr("CURRENT_TIMESTAMP")). // Default value (raw expression)
    Comment("Username").     // Column comment (MySQL only)
    Unique().                // Unique index
    After("id").             // Place after a specific column (MySQL only)
    First().                 // Place as the first column (MySQL only)
    Charset("utf8mb4").      // Character set (MySQL only)
    Collation("utf8mb4_unicode_ci"). // Collation (MySQL only)
    Unsigned().              // Unsigned (MySQL only, ignored on PostgreSQL)
    Change()                 // Modify an existing column (instead of adding)
```

> **PostgreSQL note:** Modifiers marked "MySQL only" are silently ignored when using the PostgreSQL dialect.

### Index Operations

```go
builder.Index("", "name", "email")         // Create index (name auto-generated)
builder.Index("idx_name", "name")          // Create named index
builder.Unique("", "email")                // Create unique index
builder.Primary("", "id")                  // Create primary key

builder.DropIndex("idx_name")              // Drop index
builder.DropUnique("unique_name")          // Drop unique index
builder.DropColumn("column1", "column2")   // Drop columns
```

### Foreign Keys

```go
builder.Foreign("", "user_id").
    References("id").
    On("users").
    OnDelete("CASCADE").
    OnUpdate("CASCADE")

builder.DropForeign("foreign_key_name")    // Drop foreign key
```

### Table Operations

```go
// Drop table
m.Schema("20240102_001").Drop("temp_table")
m.Schema("20240102_002").DropIfExists("temp_table")

// Raw SQL migration
m.Schema("20240102_003").Raw("custom_table", func() []string {
    return []string{
        `CREATE TABLE custom_table (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255))`,
    }
})
```

## Transactions

```go
import "github.com/mylxsw/eloquent/query"

err := eloquent.Transaction(db, func(tx query.Database) error {
    userModel := models.NewUserModel(tx)  // Use tx instead of db

    id, err := userModel.Save(ctx, models.UserN{
        Name: null.StringFrom("Tom"),
    })
    if err != nil {
        return err  // Returning an error triggers automatic rollback
    }

    // ... more operations

    return nil  // Returning nil triggers automatic commit
})
```

> Panics within a transaction also trigger automatic rollback.

## Event System

Eloquent includes a built-in event system for monitoring SQL execution, transactions, and migrations.

### Setting Up Event Listeners

```go
import "github.com/mylxsw/eloquent/event"

em := event.NewEventManager(event.NewMemoryEventStore())
event.SetDispatcher(em)

// Listen for SQL execution events (useful for logging, performance monitoring)
em.Listen(func(evt event.QueryExecutedEvent) {
    fmt.Printf("[SQL] %s | args: %v | time: %s\n", evt.SQL, evt.Bindings, evt.Time)
})

// Listen for transaction events
em.Listen(func(evt event.TransactionBeginningEvent) {
    fmt.Println("Transaction started")
})
em.Listen(func(evt event.TransactionCommittedEvent) {
    fmt.Println("Transaction committed")
})
em.Listen(func(evt event.TransactionRolledBackEvent) {
    fmt.Println("Transaction rolled back")
})

// Listen for migration events
em.Listen(func(evt event.MigrationStartedEvent) {
    fmt.Printf("Running migration: %s\n", evt.SQL)
})
```

### Available Events

| Event | Description |
|-------|-------------|
| `QueryExecutedEvent` | Fired after each SQL query, includes SQL, bindings, and duration |
| `TransactionBeginningEvent` | Fired when a transaction starts |
| `TransactionCommittedEvent` | Fired when a transaction is committed |
| `TransactionRolledBackEvent` | Fired when a transaction is rolled back |
| `MigrationsStartedEvent` | Fired before a batch of migrations starts |
| `MigrationsEndedEvent` | Fired after a batch of migrations completes |
| `MigrationStartedEvent` | Fired before a single migration SQL is executed |
| `MigrationEndedEvent` | Fired after a single migration SQL is executed |

## Direct Database Operations

In addition to model-based operations, Eloquent provides low-level database operation APIs, suitable for complex queries or scenarios where you don't want to define a model.

```go
import "github.com/mylxsw/eloquent"

dbOp := eloquent.DB(db)
ctx := context.TODO()

// Insert
id, err := dbOp.Insert(ctx, "users", query.KV{
    "name":  "Tom",
    "email": "tom@example.com",
})

// Query
results, err := dbOp.Query(ctx,
    eloquent.Build("users").Select("id", "name").Where("age", ">", 18),
    func(row eloquent.Scanner) (any, error) {
        var id int64
        var name string
        err := row.Scan(&id, &name)
        return map[string]any{"id": id, "name": name}, err
    },
)

// Generic query (Go 1.18+)
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

// Raw SQL query
results, err := dbOp.Query(ctx,
    eloquent.Raw("SELECT * FROM users WHERE age > ?", 18),
    func(row eloquent.Scanner) (any, error) {
        // ...
    },
)

// Update
affected, err := dbOp.Update(ctx,
    eloquent.Build("users").Where("id", 1),
    query.KV{"name": "Jerry"},
)

// Delete
affected, err := dbOp.Delete(ctx, eloquent.Build("users").Where("id", 1))

// Execute a statement with no return value
err := dbOp.Statement(ctx, "TRUNCATE TABLE users")
```

## Example Projects

- [_examples](./_examples) - Built-in complete example code
- [tech-share](https://github.com/mylxsw/tech-share) - A simple web project demonstrating Eloquent in a real application

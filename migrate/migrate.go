package migrate

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/mylxsw/eloquent/event"
	"github.com/mylxsw/eloquent/query"
)

type ExprType int

type Expr struct {
	Type  ExprType
	Value string
}

const (
	ExprTypeString ExprType = iota
	ExprTypeRaw
)

type Manager struct {
	db          *sql.DB
	migrateRepo *MigrationsModel

	Engine              string
	Charset             string
	Collation           string
	Prefix              string
	DefaultStringLength int
	MigrationTable      string
	migrateFuncs        []func(ctx context.Context) error

	batch int64
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db:          db,
		migrateRepo: NewMigrationsModel(db),

		Engine:              "InnoDB",
		Charset:             "utf8mb4",
		Collation:           "utf8mb4_unicode_ci",
		MigrationTable:      "migrations",
		DefaultStringLength: 255,
		Prefix:              "",
		migrateFuncs:        make([]func(ctx context.Context) error, 0),
	}
}

func (m *Manager) Init(ctx context.Context) *Manager {
	builder := NewBuilder(m.MigrationTable, "").DefaultStringLength(m.DefaultStringLength)
	builder.Engine(m.Engine)
	builder.Charset(m.Charset)
	builder.Collation(m.Collation)

	builder.Increments("id")
	builder.String("table", 0)
	builder.String("version", 20)
	builder.Text("migration")
	builder.Integer("batch", false, false)

	builder.CreateIfNotExists()

	for _, s := range builder.Build() {
		if _, err := m.db.ExecContext(ctx, s); err != nil {
			panic(err)
		}
	}

	return m
}

func (m *Manager) Schema(version string) *Schema {
	return NewSchema(m, version)
}

func (m *Manager) Execute(builder *Builder, version string) {
	m.migrateFuncs = append(m.migrateFuncs, func(ctx context.Context) error {
		tableName := builder.GetTableName()
		if m.HasVersion(ctx, version, tableName) {
			// fmt.Printf("ignore-version(%s)\n", version)
			return nil
		}

		sqls := builder.Build()

		if err := m.execute(ctx, sqls); err != nil {
			return err
		}

		if err := m.AddVersion(ctx, version, tableName, strings.Join(sqls, ";\n")+";"); err != nil {
			return err
		}

		return nil
	})
}

func (m *Manager) execute(ctx context.Context, sqls []string) error {
	event.Dispatch(event.MigrationsStartedEvent{})

	for _, s := range sqls {
		// fmt.Printf("execute -> %s\n", s)
		event.Dispatch(event.MigrationStartedEvent{SQL: s})
		if _, err := m.db.ExecContext(ctx, s); err != nil {
			return err
		}

		event.Dispatch(event.MigrationEndedEvent{SQL: s})
	}

	event.Dispatch(event.MigrationsEndedEvent{})

	return nil
}

func (m *Manager) HasVersion(ctx context.Context, version string, tableName string) bool {
	existed, err := m.migrateRepo.
		Condition(query.Builder().Where("version", version).Where("table", tableName)).
		Exists(ctx)
	if err != nil {
		panic(err)
	}

	return existed
}

func (m *Manager) AddVersion(ctx context.Context, version string, tableName string, sqlStr string) error {
	// fmt.Printf("add-version(%s)\n", version)
	_, err := m.migrateRepo.Save(
		ctx,
		Migrations{
			Version:   version,
			Table:     tableName,
			Migration: sqlStr,
			Batch:     m.batch,
		}.ToMigrationsN(),
	)

	return err
}

func (m *Manager) Run(ctx context.Context) error {
	m.batch = time.Now().Unix()

	for _, f := range m.migrateFuncs {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

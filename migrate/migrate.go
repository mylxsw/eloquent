package migrate

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
	migrateFuncs        []func() error

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
		migrateFuncs:        make([]func() error, 0),
	}
}

func (m *Manager) Init() *Manager {
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

	if err := m.execute(builder.Build()); err != nil {
		panic(err)
	}

	return m
}

func (m *Manager) Schema(version string) *Schema {
	return NewSchema(m, version)
}

func (m *Manager) Execute(builder *Builder, version string) {
	m.migrateFuncs = append(m.migrateFuncs, func() error {
		tableName := builder.GetTableName()
		if m.HasVersion(version, tableName) {
			fmt.Printf("ignore-version(%s)\n", version)
			return nil
		}

		sqls := builder.Build()

		if err := m.execute(sqls); err != nil {
			return err
		}

		if err := m.AddVersion(version, tableName, strings.Join(sqls, ";\n")+";"); err != nil {
			return err
		}

		return nil
	})
}

func (m *Manager) execute(sqls []string) error {
	for _, s := range sqls {
		fmt.Printf("execute -> %s\n", s)

		if _, err := m.db.Exec(s); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) HasVersion(version string, tableName string) bool {
	existed, err := m.migrateRepo.Query(query.Builder().Where("version", version).Where("table", tableName)).Exists()
	if err != nil {
		panic(err)
	}

	return existed
}

func (m *Manager) AddVersion(version string, tableName string, sqlStr string) error {
	fmt.Printf("add-version(%s)\n", version)

	_, err := m.migrateRepo.Save(Migrations{
		Version:   version,
		Table:     tableName,
		Migration: sqlStr,
		Batch:     m.batch,
	})

	return err
}

func (m *Manager) Run() error {
	m.batch = time.Now().Unix()

	for _, f := range m.migrateFuncs {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

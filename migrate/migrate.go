package migrate

import (
	"fmt"
	"strings"
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
	Engine    string
	Charset   string
	Collation string
	Prefix    string
}

func NewManager() *Manager {
	return &Manager{
		Engine:    "InnoDB",
		Charset:   "utf8mb4",
		Collation: "utf8mb4_unicode_ci",
		Prefix:    "",
	}
}

func (m *Manager) Schema(version string) *Schema {
	return NewSchema(m, version)
}

func (m *Manager) Execute(builder *Builder, version string) {
	if m.HasVersion(version) {
		fmt.Printf("ignore-version(%s)\n", version)
		return
	}

	sqls := builder.Build()

	if err := m.execute(sqls); err != nil {
		panic(err)
	}

	m.AddVersion(version, strings.Join(sqls, ";\n")+";")
}

func (m *Manager) execute(sqls []string) error {
	for _, s := range sqls {
		fmt.Printf("execute -> %s\n", s)
	}

	return nil
}

func (m *Manager) HasVersion(version string) bool {
	return version == "201907150945"
}

func (m *Manager) AddVersion(version string, sqlStr string) {
	fmt.Printf("add-version(%s)\n", version)
}

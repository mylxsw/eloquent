package migrate

import (
	"github.com/mylxsw/eloquent/migrate/schema"
)

type Command struct {
	t *tableBuilder

	CommandName       string
	CommandIndex      string
	CommandParameters []string
	CommandAlgorithm  string
}

func NewCommand(t *tableBuilder) schema.Command {
	return &Command{
		t:                 t,
		CommandParameters: make([]string, 0),
	}
}

func (c *Command) Equal(name string) bool {
	return c.CommandName == name
}

func (c *Command) Name(name string) schema.Command {
	c.CommandName = name
	return c
}

func (c *Command) Index(name string) schema.Command {
	c.CommandIndex = name
	return c
}

func (c *Command) Columns(columns ...string) schema.Command {
	c.CommandParameters = columns
	return c
}

func (c *Command) Algorithm(algorithm string) schema.Command {
	c.CommandAlgorithm = algorithm
	return c
}

func (c *Command) Build() string {
	switch c.CommandName {
	case "index":
		return c.t.compileKey(c, "index")
	case "unique":
		return c.t.compileKey(c, "unique")
	case "primary":
		c.CommandIndex = ""
		return c.t.compileKey(c, "primary key")
	case "spatialIndex":
		return c.t.compileKey(c, "spatial index")
	case "dropTable":
		return c.t.compileDrop(c)
	case "dropColumn":
		return c.t.compileDropColumn(c)
	case "dropIndex", "dropUnique", "dropSpatialIndex":
		return c.t.compileDropIndex(c)
	case "dropPrimary":
		return c.t.compileDropPrimary(c)
	case "drop":
		return c.t.compileDrop(c)
	case "rename":
		return c.t.compileRename(c)
	case "create":
		return c.t.compileCreateCommand()
	case "add":
		return c.t.compileAdd()
	case "change":
		return c.t.compileChange()
	case "renameColumn":
		return c.t.compileRenameColumn()
	}

	return ""
}

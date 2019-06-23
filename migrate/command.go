package migrate

type Command struct {
	t *Builder

	CommandName       string
	CommandIndex      string
	CommandParameters []string
	CommandAlgorithm  string

	CommandReferences            []string
	CommandOnTable               string
	CommandOnDelete              string
	CommandOnUpdate              string
	CommandNotInitiallyImmediate bool
}

func NewCommand(t *Builder) *Command {
	return &Command{
		t:                 t,
		CommandParameters: make([]string, 0),
	}
}

func (c *Command) References(columns ...string) *Command {
	c.CommandReferences = columns
	return c
}

func (c *Command) On(table string) *Command {
	c.CommandOnTable = table
	return c
}

func (c *Command) OnDelete(action string) *Command {
	c.CommandOnDelete = action
	return c
}

func (c *Command) OnUpdate(action string) *Command {
	c.CommandOnUpdate = action
	return c
}

func (c *Command) NotInitiallyImmediate(value bool) *Command {
	c.CommandNotInitiallyImmediate = value
	return c
}

func (c *Command) Equal(name string) bool {
	return c.CommandName == name
}

func (c *Command) Name(name string) *Command {
	c.CommandName = name
	return c
}

func (c *Command) Index(name string) *Command {
	c.CommandIndex = name
	return c
}

func (c *Command) Columns(columns ...string) *Command {
	c.CommandParameters = columns
	return c
}

func (c *Command) Algorithm(algorithm string) *Command {
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
		return c.t.compileCreateCommand(false)
	case "createIfNotExists":
		return c.t.compileCreateCommand(true)
	case "add":
		return c.t.compileAdd()
	case "change":
		return c.t.compileChange()
	case "renameColumn":
		return c.t.compileRenameColumn(c.CommandParameters[0], c.CommandParameters[1])
	case "foreign":
		return c.t.compileForeign(c)
	case "dropForeign":
		return c.t.compileDropForeign(c.CommandParameters[0])
	}

	return ""
}

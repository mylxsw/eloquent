package generator

import (
	"github.com/iancoleman/strcase"
)

type Domain struct {
	Imports     []string `yaml:"imports"`
	PackageName string   `yaml:"package"`
	Models      []Model  `yaml:"models"`
	Meta        Meta     `yaml:"meta"`
}

func (dom Domain) Init() Domain {
	for i, m := range dom.Models {
		for j, f := range m.Definition.Fields {
			dom.Models[i].Definition.Fields[j].Name = strcase.ToCamel(f.Name)
			if f.Type == "" {
				dom.Models[i].Definition.Fields[j].Type = "string"
			}
		}
	}

	return dom
}

type Meta struct {
	TablePrefix string `yaml:"table_prefix"`
}

type Model struct {
	Name       string     `yaml:"name"`
	Relations  []Relation `yaml:"relations"`
	Definition Definition `yaml:"definition"`
}

func (rel Relation) ImportPackages() []string {
	internalPackages := make([]string, 0)

	if rel.Package != "" {
		internalPackages = append(internalPackages, rel.Package)
	}

	if relationRel(rel) == "belongsToMany" {
		internalPackages = append(internalPackages, "github.com/mylxsw/eloquent")
	}

	return unique(internalPackages)
}

type Relation struct {
	Model string `yaml:"model"`
	Rel   string `yaml:"rel"`

	ForeignKey string `yaml:"foreign_key"`
	OwnerKey   string `yaml:"owner_key"`
	LocalKey   string `yaml:"local_key"`

	PivotTable string `yaml:"table"`

	Package string `yaml:"package"`
	Method  string `yaml:"method"`
}

type Definition struct {
	TableName         string            `yaml:"table_name"`
	WithoutCreateTime bool              `yaml:"without_create_time"`
	WithoutUpdateTime bool              `yaml:"without_update_time"`
	SoftDelete        bool              `yaml:"soft_delete"`
	Fields            []DefinitionField `yaml:"fields"`
}

type DefinitionField struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Tag  string `yaml:"tag"`
}


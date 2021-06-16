package generator

import (
	"github.com/iancoleman/strcase"
)

type Domain struct {
	Imports     []string `yaml:"imports,omitempty"`
	PackageName string   `yaml:"package,omitempty"`
	Models      []Model  `yaml:"models,omitempty"`
	Meta        Meta     `yaml:"meta,omitempty"`
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
	Name       string     `yaml:"name,omitempty"`
	Relations  []Relation `yaml:"relations,omitempty"`
	Definition Definition `yaml:"definition,omitempty"`
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
	Model string `yaml:"model,omitempty"`
	Rel   string `yaml:"rel,omitempty"`

	ForeignKey string `yaml:"foreign_key,omitempty"`
	OwnerKey   string `yaml:"owner_key,omitempty"`
	LocalKey   string `yaml:"local_key,omitempty"`

	PivotTable string `yaml:"table,omitempty"`

	Package string `yaml:"package,omitempty"`
	Method  string `yaml:"method,omitempty"`
}

type Definition struct {
	TableName         string            `yaml:"table_name,omitempty"`
	WithoutCreateTime bool              `yaml:"without_create_time,omitempty"`
	WithoutUpdateTime bool              `yaml:"without_update_time,omitempty"`
	SoftDelete        bool              `yaml:"soft_delete,omitempty"`
	Fields            []DefinitionField `yaml:"fields,omitempty"`
}

type DefinitionField struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
	Tag  string `yaml:"tag,omitempty"`
}

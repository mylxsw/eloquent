package generator

import (
	"strings"
)

type DomainContext struct {
	domain Domain
}

func (d DomainContext) Register() {
	AddFunc("table", d.tableName)
	AddFunc("packages", d.importPackages)
	AddFunc("assignable_fields", d.assignableFields)
}

func (d DomainContext) tableName(i int) string {
	m := d.domain.Models[i]
	if m.Definition.TableName != "" {
		return d.domain.Meta.TablePrefix + m.Definition.TableName
	}

	return d.domain.Meta.TablePrefix + strings.ToLower(m.Name)
}

func (d DomainContext) assignableFields(def Definition) []DefinitionField {
	fields := make([]DefinitionField, 0)
	for _, f := range entityFields(def) {
		if f.Name == "Id" {
			continue
		}

		if !def.WithoutCreateTime && f.Name == "CreatedAt" {
			continue
		}

		if !def.WithoutUpdateTime && f.Name == "UpdatedAt" {
			continue
		}

		if def.SoftDelete && f.Name == "DeletedAt" {
			continue
		}

		fields = append(fields, f)
	}

	return fields
}

func (d DomainContext) importPackages() []string {
	var internalPackages = []string{
		"context",
		"gopkg.in/guregu/null.v3",
		"github.com/mylxsw/eloquent/query",
		"github.com/iancoleman/strcase",
	}

	for _, m := range d.domain.Models {
		if m.Definition.SoftDelete || !m.Definition.WithoutCreateTime || !m.Definition.WithoutUpdateTime {
			internalPackages = append(internalPackages, "time")
		}

		for _, rel := range m.Relations {
			internalPackages = append(internalPackages, rel.ImportPackages()...)
		}
	}

	for _, imp := range d.domain.Imports {
		internalPackages = append(internalPackages, imp)
	}

	return unique(internalPackages)
}

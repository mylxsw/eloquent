package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

type Domain struct {
	Imports     []string `yaml:"imports"`
	PackageName string   `yaml:"package"`
	Models      []Model  `yaml:"models"`
	Meta        Meta     `yaml:"meta"`
}

type Meta struct {
	TablePrefix string `yaml:"table_prefix"`
}

type Model struct {
	Name       string     `yaml:"name"`
	Definition Definition `yaml:"definition"`
}

type Definition struct {
	TableName         string            `yaml:"table_name"`
	WithoutCreateTime bool              `yaml:"without_create_time"`
	WithoutUpdateTime bool              `yaml:"without_update_time"`
	SoftDelete        bool              `yaml:"soft_delete"`
	Fields            []DefinitionField `yaml:"fields"`
}

type DefinitionField struct {
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
	RelationField bool   `yaml:"relation_field"`
}

// ParseTemplate 模板解析
func ParseTemplate(templateContent string, data Domain) (string, error) {
	ctx := DomainContext{domain: data}
	funcMap := template.FuncMap{
		"implode":           strings.Join,
		"trim":              strings.Trim,
		"trim_right":        strings.TrimRight,
		"trim_left":         strings.TrimLeft,
		"trim_space":        strings.TrimSpace,
		"lowercase":         strings.ToLower,
		"format":            fmt.Sprintf,
		"assignable_fields": ctx.assignableFields,
		"snake":             strcase.ToSnake,
		"camel":             strcase.ToCamel,
		"lower_camel":       strcase.ToLowerCamel,
		"table":             ctx.tableName,
		"wrap_type":         wrapType,
		"unwrap_type":       unWrapType,
		"unique":            unique,
		"packages":          ctx.importPackages,
	}
	var buffer bytes.Buffer
	if err := template.Must(template.New("").Funcs(funcMap).Parse(templateContent)).Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func wrapType(t string) string {
	switch t {
	case "int64", "int", "int8", "int32":
		return "null.Int"
	case "string":
		return "null.String"
	case "time.Time":
		return "null.Time"
	case "float32", "float64":
		return "null.Float"
	case "bool":
		return "null.Bool"
	}

	return t
}

func unWrapType(name string, t string) string {
	base := "w." + strcase.ToCamel(name)
	// w.{{ camel $f.ColumnName }}{{ wrap_type $f.Type | unwrap_type }}
	if strings.HasPrefix(t, "null.") {
		return base
	}

	switch wrapType(t) {
	case "null.Int":
		if t == "int64" {
			return base + ".Int64"
		}
		return t + "(" + base + ".Int64)"
	case "null.String":
		return base + ".String"
	case "null.Time":
		return base + ".Time"
	case "null.Float":
		if t == "float64" {
			return base + ".Float64"
		}

		return t + "(" + base + ".Float64)"
	case "null.Bool":
		return base + ".Bool"
	}

	return base
}

type DomainContext struct {
	domain Domain
}

func (d DomainContext) tableName(i int) string {
	m := d.domain.Models[i]
	if m.Definition.TableName != "" {
		return d.domain.Meta.TablePrefix + m.Definition.TableName
	}

	return d.domain.Meta.TablePrefix + strings.ToLower(m.Name)
}

func (d DomainContext) assignableFields(fields []DefinitionField) []DefinitionField {
	var res = make([]DefinitionField, 0)
	for _, f := range fields {
		if f.RelationField {
			continue
		}

		res = append(res, f)
	}

	return res
}

func (d DomainContext) importPackages() []string {
	var internalPackages = []string{
		"database/sql",
		"gopkg.in/guregu/null.v3",
		"github.com/mylxsw/eloquent/query",
	}

	for _, m := range d.domain.Models {
		if m.Definition.SoftDelete || !m.Definition.WithoutCreateTime || !m.Definition.WithoutUpdateTime {
			internalPackages = append(internalPackages, "time")
		}
	}

	for _, imp := range d.domain.Imports {
		internalPackages = append(internalPackages, imp)
	}

	return unique(internalPackages)
}

func unique(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}

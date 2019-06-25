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

type Relation struct {
	Model      string `yaml:"model"`
	Rel        string `yaml:"rel"`
	ForeignKey string `yaml:"foreign_key"`
	OwnerKey   string `yaml:"owner_key"`
	LocalKey   string `yaml:"local_key"`

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

// ParseTemplate 模板解析
func ParseTemplate(templateContent string, data Domain) (string, error) {
	ctx := DomainContext{domain: data}
	funcMap := template.FuncMap{
		"implode":            strings.Join,
		"trim":               strings.Trim,
		"trim_right":         strings.TrimRight,
		"trim_left":          strings.TrimLeft,
		"trim_space":         strings.TrimSpace,
		"lowercase":          strings.ToLower,
		"format":             fmt.Sprintf,
		"assignable_fields":  ctx.assignableFields,
		"snake":              strcase.ToSnake,
		"camel":              strcase.ToCamel,
		"lower_camel":        strcase.ToLowerCamel,
		"table":              ctx.tableName,
		"wrap_type":          wrapType,
		"unwrap_type":        unWrapType,
		"unique":             unique,
		"packages":           ctx.importPackages,
		"fields":             entityFields,
		"tag":                entityTags,
		"rel_owner_key":      relationOwnerKey,
		"rel_foreign_key":    relationForeignKey,
		"rel_local_key":      relationLocalKey,
		"rel_package_prefix": relationPackagePrefix,
		"rel_method":         relationMethod,
		"rel":                relationRel,
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
	}

	for _, m := range d.domain.Models {
		if m.Definition.SoftDelete || !m.Definition.WithoutCreateTime || !m.Definition.WithoutUpdateTime {
			internalPackages = append(internalPackages, "time")
		}

		for _, rel := range m.Relations {
			if rel.Package != "" {
				internalPackages = append(internalPackages, rel.Package)
			}
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
		if encountered[elements[v]] {
			continue
		}

		encountered[elements[v]] = true
		result = append(result, elements[v])
	}

	return result
}

func entityFields(def Definition) []DefinitionField {
	fields := make([]DefinitionField, 0)

	for _, f := range def.Fields {
		f.Name = strcase.ToCamel(f.Name)
		fields = append(fields, f)
	}

	fields = append(fields, DefinitionField{Name: "Id", Type: "int64"})

	if !def.WithoutCreateTime {
		fields = append(fields, DefinitionField{Name: "CreatedAt", Type: "time.Time"})
	}

	if !def.WithoutUpdateTime {
		fields = append(fields, DefinitionField{Name: "UpdatedAt", Type: "time.Time"})
	}

	if def.SoftDelete {
		fields = append(fields, DefinitionField{Name: "DeletedAt", Type: "null.Time"})
	}

	return uniqueFields(fields)
}

func uniqueFields(elements []DefinitionField) []DefinitionField {
	encountered := map[string]bool{}
	var result []DefinitionField

	for v := range elements {
		if encountered[elements[v].Name] {
			continue
		}

		encountered[elements[v].Name] = true
		result = append(result, elements[v])
	}

	return result
}

func entityTags(field DefinitionField) string {
	if field.Tag == "" {
		return ""
	}

	return "`" + field.Tag + "`"
}

func relationForeignKey(rel Relation) string {
	if rel.ForeignKey == "" {
		return strings.ToLower(rel.Model) + "_id"
	}

	return rel.ForeignKey
}

func relationOwnerKey(rel Relation) string {
	if rel.OwnerKey == "" {
		return "id"
	}

	return rel.OwnerKey
}

func relationLocalKey(rel Relation) string {
	if rel.LocalKey == "" {
		return "id"
	}

	return rel.LocalKey
}

func relationPackagePrefix(rel Relation) string {
	if rel.Package == "" {
		return ""
	}

	segs := strings.Split(rel.Package, "/")
	lastSeg := segs[len(segs)-1]

	return lastSeg + "."
}

func relationMethod(rel Relation) string {
	if rel.Method == "" {
		switch rel.Rel {
		case "belongsTo":
			return strcase.ToCamel(rel.Model)
		case "hasMany":
			return strcase.ToCamel(rel.Model) + "s"
		}
	}

	return strcase.ToCamel(rel.Method)
}

func relationRel(rel Relation) string {
	switch strings.ToLower(rel.Rel) {
	case "belongsto", "belongs_to", "n-1", "n:1", "*:1", "*-1", "-1":
		return "belongsTo"
	case "hasmany", "has_many", "1-n", "1:n", "1:*", "1-*", "1-":
		return "hasMany"
	}

	panic(fmt.Sprintf("not support: %s", rel.Rel))
}

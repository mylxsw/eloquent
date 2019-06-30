package generator

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

func init() {
	AddFunc("wrap_type", wrapType)
	AddFunc("unwrap_type", unWrapType)
	AddFunc("unique", unique)
	AddFunc("fields", entityFields)
	AddFunc("tag", entityTags)
	AddFunc("rel_owner_key", relationOwnerKey)
	AddFunc("rel_foreign_key", relationForeignKey)
	AddFunc("rel_foreign_key_rev", relationForeignKeyRev)
	AddFunc("rel_local_key", relationLocalKey)
	AddFunc("rel_package_prefix", relationPackagePrefix)
	AddFunc("rel_method", relationMethod)
	AddFunc("rel", relationRel)
	AddFunc("rel_belongs_to_name", relationBelongsToName)
	AddFunc("rel_has_many_name", relationHasManyName)
	AddFunc("rel_has_one_name", relationHasOneName)
	AddFunc("rel_belongs_to_many_name", relationBelongsToManyName)
	AddFunc("rel_pivot_table_name", relationPivotTable)
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

func relationForeignKeyRev(rel Relation, m Model) string {
	if rel.ForeignKey == "" {
		return strcase.ToSnake(m.Name) + "_id"
	}

	return rel.ForeignKey
}

func relationForeignKey(rel Relation) string {
	if rel.ForeignKey == "" {
		return strcase.ToSnake(rel.Model) + "_id"
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
		switch relationRel(rel) {
		case "belongsTo", "hasOne":
			return strcase.ToCamel(rel.Model)
		case "hasMany", "belongsToMany":
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
	case "hasone", "has_one", "1-1", "1:1":
		return "hasOne"
	case "belongs_to_many", "belongstomany", "n:n", "n-n", "*:*", "*-*":
		return "belongsToMany"
	}

	panic(fmt.Sprintf("not support: %s", rel.Rel))
}

func relationBelongsToName(rel Relation, m Model) string {
	return fmt.Sprintf("%sBelongsTo%sRel", strcase.ToCamel(m.Name), strcase.ToCamel(rel.Model))
}

func relationHasManyName(rel Relation, m Model) string {
	return fmt.Sprintf("%sHasMany%sRel", strcase.ToCamel(m.Name), strcase.ToCamel(rel.Model))
}

func relationHasOneName(rel Relation, m Model) string {
	return fmt.Sprintf("%sHasOne%sRel", strcase.ToCamel(m.Name), strcase.ToCamel(rel.Model))
}

func relationBelongsToManyName(rel Relation, m Model) string {
	return fmt.Sprintf("%sBelongsToMany%sRel", strcase.ToCamel(m.Name), strcase.ToCamel(rel.Model))
}

func relationPivotTable(rel Relation, m Model) string {
	if rel.PivotTable != "" {
		return rel.PivotTable
	}

	t1 := strcase.ToSnake(rel.Model)
	t2 := strcase.ToSnake(m.Name)

	if strings.Compare(t1, t2) > 0 {
		return fmt.Sprintf("%s_%s_ref", t1, t2)
	}

	return fmt.Sprintf("%s_%s_ref", t2, t1)
}

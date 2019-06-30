package template

import "fmt"

func GetRelationTemplate() string {
	temp := `{{ range $j, $rel := $m.Relations }}{{ if rel $rel | eq "belongsTo" }}
%s
{{ end }}{{ if rel $rel | eq "hasMany" }}
%s
{{ end }}{{ if rel $rel | eq "hasOne" }}
%s
{{ end }}{{ if rel $rel | eq "belongsToMany" }}
%s
{{ end }}{{ end }}
`
	return fmt.Sprintf(
		temp,
		getRelationBelongsToTemplate(),
		getRelationHasManyTemplate(),
		getRelationhasOneTemplate(),
		getRelationBelongsToManyTemplate(),
	)
}

func getRelationBelongsToTemplate() string {
	return `{{ $relName := rel_belongs_to_name $rel $m }}
func (inst *{{ camel $m.Name }}) {{ rel_method $rel }}() *{{ $relName }} {
	return &{{ $relName }} {
		source: inst,
		relModel: {{ rel_package_prefix $rel }}New{{ camel $rel.Model }}Model(inst.{{ lower_camel $m.Name }}Model.GetDB()),
	}
}

type {{ $relName }} struct {
	source *{{ camel $m.Name }}
	relModel *{{ rel_package_prefix $rel }}{{ camel $rel.Model }}Model
}

func (rel *{{ $relName }}) Create(target {{ camel $rel.Model }}) (int64, error) {
	targetId, err := rel.relModel.Save(target)
	if err != nil {
		return 0, err
	}

	target.Id = targetId

	rel.source.{{ rel_foreign_key $rel | camel }} = target.{{ rel_owner_key $rel | camel }}
	if err := rel.source.Save(); err != nil {
		return targetId, err
	}

	return targetId, nil
}

func (rel *{{ $relName }}) Exists(builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("{{ rel_owner_key $rel | snake }}", rel.source.{{ rel_foreign_key $rel | camel }}).Merge(builders...)
	
	return rel.relModel.Exists(builder)
}

func (rel *{{ $relName }}) First(builders ...query.SQLBuilder) ({{ camel $rel.Model }}, error) {
	builder := query.Builder().Where("{{ rel_owner_key $rel | snake }}", rel.source.{{ rel_foreign_key $rel | camel }}).Limit(1).Merge(builders...)

	return rel.relModel.First(builder)
}

func (rel *{{ $relName }}) Associate(target {{ camel $rel.Model }}) error {
	rel.source.{{ rel_foreign_key $rel | camel }} = target.{{ rel_owner_key $rel | camel }}
	return rel.source.Save()
}

func (rel *{{ $relName }}) Dissociate() error {
	rel.source.{{ rel_foreign_key $rel | camel }} = 0
	return rel.source.Save()
}
`
}

func getRelationHasManyTemplate() string {
	return `{{ $relName := rel_has_many_name $rel $m }}
func (inst *{{ camel $m.Name }}) {{ rel_method $rel }}() *{{ $relName }} {
	return &{{ $relName }} {
		source: inst,
		relModel: {{ rel_package_prefix $rel }}New{{ camel $rel.Model }}Model(inst.{{ lower_camel $m.Name }}Model.GetDB()),
	}
}

type {{ $relName }} struct {
	source *{{ camel $m.Name }}
	relModel *{{ rel_package_prefix $rel }}{{ camel $rel.Model }}Model
}

func (rel *{{ $relName }}) Get(builders ...query.SQLBuilder) ([]{{ camel $rel.Model }}, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Merge(builders...)

	return rel.relModel.Get(builder)
}

func (rel *{{ $relName }}) Count(builders ...query.SQLBuilder) (int64, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Merge(builders...)
	
	return rel.relModel.Count(builder)
}

func (rel *{{ $relName }}) Exists(builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Merge(builders...)
	
	return rel.relModel.Exists(builder)
}

func (rel *{{ $relName }}) First(builders ...query.SQLBuilder) ({{ camel $rel.Model }}, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Limit(1).Merge(builders...)
	return rel.relModel.First(builder)
}

func (rel *{{ $relName }}) Create(target {{ camel $rel.Model }}) (int64, error) {
	target.{{ rel_foreign_key_rev $rel $m | camel }} = rel.source.Id
	return rel.relModel.Save(target)
}
`
}

func getRelationhasOneTemplate() string {
	return `{{ $relName := rel_has_one_name $rel $m }}
func (inst *{{ camel $m.Name }}) {{ rel_method $rel }}() *{{ $relName }} {
	return &{{ $relName }} {
		source: inst,
		relModel: {{ rel_package_prefix $rel }}New{{ camel $rel.Model }}Model(inst.{{ lower_camel $m.Name }}Model.GetDB()),
	}
}

type {{ $relName }} struct {
	source *{{ camel $m.Name }}
	relModel *{{ rel_package_prefix $rel }}{{ camel $rel.Model }}Model
}

func (rel *{{ $relName }}) Exists(builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Merge(builders...)
	
	return rel.relModel.Exists(builder)
}

func (rel *{{ $relName }}) First(builders ...query.SQLBuilder) ({{ camel $rel.Model }}, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Limit(1).Merge(builders...)
	return rel.relModel.First(builder)
}

func (rel *{{ $relName }}) Create(target {{ camel $rel.Model }}) (int64, error) {
	target.{{ rel_foreign_key_rev $rel $m | camel }} = rel.source.{{ rel_local_key $rel | camel }}
	return rel.relModel.Save(target)
}

func (rel *{{ $relName }}) Associate(target {{ camel $rel.Model }}) error {
	_, err := rel.relModel.UpdateFields(
		query.KV {"{{ rel_foreign_key_rev $rel $m | snake }}": rel.source.{{ rel_local_key $rel | camel }}, },
		query.Builder().Where("id", target.Id), 
	)
	return err
}

func (rel *{{ $relName }}) Dissociate() error {
	_, err := rel.relModel.UpdateFields(
		query.KV {"{{ rel_foreign_key_rev $rel $m | snake }}": nil,},
		query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}),
	)

	return err
}
`
}

func getRelationBelongsToManyTemplate() string {
	return `{{ $relName := rel_belongs_to_many_name $rel $m }}
func (inst *{{ camel $m.Name }}) {{ rel_method $rel }}() *{{ $relName }} {
	return &{{ $relName }} {
		source: inst,
		pivotTable: "{{ rel_pivot_table_name $rel $m | snake }}",
		relModel: {{ rel_package_prefix $rel }}New{{ camel $rel.Model }}Model(inst.{{ lower_camel $m.Name }}Model.GetDB()),
	}
}

type {{ $relName }} struct {
	source *{{ camel $m.Name }}
	pivotTable string
	relModel *{{ rel_package_prefix $rel }}{{ camel $rel.Model }}Model
}

func (rel *{{ $relName }}) Get(builders ...query.SQLBuilder) ([]{{ camel $rel.Model }}, error) {
	res, err := eloquent.DB(rel.relModel.GetDB()).Query(
		query.Builder().Table(rel.pivotTable).Select("{{ rel_foreign_key $rel | snake }}").Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_owner_key $rel | camel }}),
		func(row *sql.Rows) (interface{}, error) {
			var k interface{}
			if err := row.Scan(&k); err != nil {
				return nil, err
			}

			return k, nil
		},
	)

	if err != nil {
		return nil, err
	}

	resArr, _ := res.ToArray()
	return rel.relModel.Get(query.Builder().Merge(builders...).WhereIn("{{ rel_owner_key $rel | snake }}", resArr...))
}

func (rel *{{ $relName }}) Count(builders ...query.SQLBuilder) (int64, error) {
	res, err := eloquent.DB(rel.relModel.GetDB()).Query(
		query.Builder().Table(rel.pivotTable).Select(query.Raw("COUNT(1) as c")).Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_owner_key $rel | camel }}),
		func(row *sql.Rows) (interface{}, error) {
			var k int64
			if err := row.Scan(&k); err != nil {
				return nil, err
			}

			return k, nil
		},
	)

	if err != nil {
		return 0, err
	}

	return res.Index(0).(int64), nil
}

func (rel *{{ $relName }}) Exists(builders ...query.SQLBuilder) (bool, error) {
	c, err := rel.Count(builders...)
	if err != nil {
		return false, err
	}
	
	return c > 0, nil
}

func (rel *{{ $relName }}) Attach(target {{ camel $rel.Model }}) error {
	_, err := eloquent.DB(rel.relModel.GetDB()).Insert(rel.pivotTable, query.KV {
		"{{ rel_foreign_key $rel | snake }}": target.{{ rel_owner_key $rel | camel }},
		"{{ rel_foreign_key_rev $rel $m | snake }}": rel.source.{{ rel_owner_key $rel | camel }},
	})

	return err
}

func (rel *{{ $relName }}) Detach(target {{ camel $rel.Model }}) error {
	_, err := eloquent.DB(rel.relModel.GetDB()).
		Delete(eloquent.Build(rel.pivotTable).
			Where("{{ rel_foreign_key $rel | snake }}", target.{{ rel_owner_key $rel | camel }}).
			Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_owner_key $rel | camel }}),)
	
	return err
}

func (rel *{{ $relName }}) DetachAll() error {
	_, err := eloquent.DB(rel.relModel.GetDB()).
		Delete(eloquent.Build(rel.pivotTable).
			Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_owner_key $rel | camel }}),)
	return err
}

func (rel *{{ $relName }}) Create(target {{ camel $rel.Model }}, builders ...query.SQLBuilder) (int64, error) {
	targetId, err := rel.relModel.Save(target)
	if err != nil {
		return 0, err
	}

	target.Id = targetId

	err = rel.Attach(target)

	return targetId, err
}
`
}

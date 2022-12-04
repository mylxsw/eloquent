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
		getRelationHasOneTemplate(),
		getRelationBelongsToManyTemplate(),
	)
}

func getRelationBelongsToTemplate() string {
	return `{{ $relName := rel_belongs_to_name $rel $m }}
func (inst *{{ camel $m.Name }}N) {{ rel_method $rel }}() *{{ $relName }} {
	return &{{ $relName }} {
		source: inst,
		relModel: {{ rel_package_prefix $rel }}New{{ camel $rel.Model }}Model(inst.{{ lower_camel $m.Name }}Model.GetDB()),
	}
}

type {{ $relName }} struct {
	source *{{ camel $m.Name }}N
	relModel *{{ rel_package_prefix $rel }}{{ camel $rel.Model }}Model
}

func (rel *{{ $relName }}) Create(ctx context.Context, target {{ camel $rel.Model }}N) (int64, error) {
	targetId, err := rel.relModel.Save(ctx, target)
	if err != nil {
		return 0, err
	}

	target.Id = null.IntFrom(targetId)

	rel.source.{{ rel_foreign_key $rel | camel }} = target.{{ rel_owner_key $rel | camel }}
	if err := rel.source.Save(ctx); err != nil {
		return targetId, err
	}

	return targetId, nil
}

func (rel *{{ $relName }}) Exists(ctx context.Context, builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("{{ rel_owner_key $rel | snake }}", rel.source.{{ rel_foreign_key $rel | camel }}).Merge(builders...)
	
	return rel.relModel.Exists(ctx, builder)
}

func (rel *{{ $relName }}) First(ctx context.Context, builders ...query.SQLBuilder) (*{{ camel $rel.Model }}N, error) {
	builder := query.Builder().Where("{{ rel_owner_key $rel | snake }}", rel.source.{{ rel_foreign_key $rel | camel }}).Limit(1).Merge(builders...)

	return rel.relModel.First(ctx, builder)
}

func (rel *{{ $relName }}) Associate(ctx context.Context, target {{ camel $rel.Model }}N) error {
	rel.source.{{ rel_foreign_key $rel | camel }} = target.{{ rel_owner_key $rel | camel }}
	return rel.source.Save(ctx)
}

func (rel *{{ $relName }}) Dissociate(ctx context.Context) error {
	rel.source.{{ rel_foreign_key $rel | camel }} = null.IntFrom(0)
	return rel.source.Save(ctx)
}
`
}

func getRelationHasManyTemplate() string {
	return `{{ $relName := rel_has_many_name $rel $m }}
func (inst *{{ camel $m.Name }}N) {{ rel_method $rel }}() *{{ $relName }} {
	return &{{ $relName }} {
		source: inst,
		relModel: {{ rel_package_prefix $rel }}New{{ camel $rel.Model }}Model(inst.{{ lower_camel $m.Name }}Model.GetDB()),
	}
}

type {{ $relName }} struct {
	source *{{ camel $m.Name }}N
	relModel *{{ rel_package_prefix $rel }}{{ camel $rel.Model }}Model
}

func (rel *{{ $relName }}) Get(ctx context.Context, builders ...query.SQLBuilder) ([]{{ camel $rel.Model }}N, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Merge(builders...)

	return rel.relModel.Get(ctx, builder)
}

func (rel *{{ $relName }}) Count(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Merge(builders...)
	
	return rel.relModel.Count(ctx, builder)
}

func (rel *{{ $relName }}) Exists(ctx context.Context, builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Merge(builders...)
	
	return rel.relModel.Exists(ctx, builder)
}

func (rel *{{ $relName }}) First(ctx context.Context, builders ...query.SQLBuilder) (*{{ camel $rel.Model }}N, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Limit(1).Merge(builders...)
	return rel.relModel.First(ctx, builder)
}

func (rel *{{ $relName }}) Create(ctx context.Context, target {{ camel $rel.Model }}N) (int64, error) {
	target.{{ rel_foreign_key_rev $rel $m | camel }} = rel.source.Id
	return rel.relModel.Save(ctx, target)
}
`
}

func getRelationHasOneTemplate() string {
	return `{{ $relName := rel_has_one_name $rel $m }}
func (inst *{{ camel $m.Name }}N) {{ rel_method $rel }}() *{{ $relName }} {
	return &{{ $relName }} {
		source: inst,
		relModel: {{ rel_package_prefix $rel }}New{{ camel $rel.Model }}Model(inst.{{ lower_camel $m.Name }}Model.GetDB()),
	}
}

type {{ $relName }} struct {
	source *{{ camel $m.Name }}N
	relModel *{{ rel_package_prefix $rel }}{{ camel $rel.Model }}Model
}

func (rel *{{ $relName }}) Exists(ctx context.Context, builders ...query.SQLBuilder) (bool, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Merge(builders...)
	
	return rel.relModel.Exists(ctx, builder)
}

func (rel *{{ $relName }}) First(ctx context.Context, builders ...query.SQLBuilder) (*{{ camel $rel.Model }}N, error) {
	builder := query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}).Limit(1).Merge(builders...)
	return rel.relModel.First(ctx, builder)
}

func (rel *{{ $relName }}) Create(ctx context.Context, target {{ camel $rel.Model }}N) (int64, error) {
	target.{{ rel_foreign_key_rev $rel $m | camel }} = rel.source.{{ rel_local_key $rel | camel }}
	return rel.relModel.Save(ctx, target)
}

func (rel *{{ $relName }}) Associate(ctx context.Context, target {{ camel $rel.Model }}N) error {
	_, err := rel.relModel.UpdateFields(
		ctx,
		query.KV {"{{ rel_foreign_key_rev $rel $m | snake }}": rel.source.{{ rel_local_key $rel | camel }}, },
		query.Builder().Where("id", target.Id), 
	)
	return err
}

func (rel *{{ $relName }}) Dissociate(ctx context.Context) error {
	_, err := rel.relModel.UpdateFields(
		ctx,
		query.KV {"{{ rel_foreign_key_rev $rel $m | snake }}": nil,},
		query.Builder().Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_local_key $rel | camel }}),
	)

	return err
}
`
}

func getRelationBelongsToManyTemplate() string {
	return `{{ $relName := rel_belongs_to_many_name $rel $m }}
func (inst *{{ camel $m.Name }}N) {{ rel_method $rel }}() *{{ $relName }} {
	return &{{ $relName }} {
		source: inst,
		pivotTable: "{{ rel_pivot_table_name $rel $m | snake }}",
		relModel: {{ rel_package_prefix $rel }}New{{ camel $rel.Model }}Model(inst.{{ lower_camel $m.Name }}Model.GetDB()),
	}
}

type {{ $relName }} struct {
	source *{{ camel $m.Name }}N
	pivotTable string
	relModel *{{ rel_package_prefix $rel }}{{ camel $rel.Model }}Model
}

func (rel *{{ $relName }}) Get(ctx context.Context, builders ...query.SQLBuilder) ([]{{ camel $rel.Model }}N, error) {
	res, err := eloquent.DB(rel.relModel.GetDB()).Query(
		ctx,
		query.Builder().Table(rel.pivotTable).Select("{{ rel_foreign_key $rel | snake }}").Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_owner_key $rel | camel }}),
		func(row eloquent.Scanner) (interface{}, error) {
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
	return rel.relModel.Get(ctx, query.Builder().Merge(builders...).WhereIn("{{ rel_owner_key $rel | snake }}", resArr...))
}

func (rel *{{ $relName }}) Count(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {
	res, err := eloquent.DB(rel.relModel.GetDB()).Query(
		ctx,
		query.Builder().Table(rel.pivotTable).Select(query.Raw("COUNT(1) as c")).Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_owner_key $rel | camel }}),
		func(row eloquent.Scanner) (interface{}, error) {
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

func (rel *{{ $relName }}) Exists(ctx context.Context, builders ...query.SQLBuilder) (bool, error) {
	c, err := rel.Count(ctx, builders...)
	if err != nil {
		return false, err
	}
	
	return c > 0, nil
}

func (rel *{{ $relName }}) Attach(ctx context.Context, target {{ camel $rel.Model }}N) error {
	_, err := eloquent.DB(rel.relModel.GetDB()).Insert(ctx, rel.pivotTable, query.KV {
		"{{ rel_foreign_key $rel | snake }}": target.{{ rel_owner_key $rel | camel }},
		"{{ rel_foreign_key_rev $rel $m | snake }}": rel.source.{{ rel_owner_key $rel | camel }},
	})

	return err
}

func (rel *{{ $relName }}) Detach(ctx context.Context, target {{ camel $rel.Model }}N) error {
	_, err := eloquent.DB(rel.relModel.GetDB()).
		Delete(ctx, eloquent.Build(rel.pivotTable).
			Where("{{ rel_foreign_key $rel | snake }}", target.{{ rel_owner_key $rel | camel }}).
			Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_owner_key $rel | camel }}),)
	
	return err
}

func (rel *{{ $relName }}) DetachAll(ctx context.Context) error {
	_, err := eloquent.DB(rel.relModel.GetDB()).
		Delete(ctx, eloquent.Build(rel.pivotTable).
			Where("{{ rel_foreign_key_rev $rel $m | snake }}", rel.source.{{ rel_owner_key $rel | camel }}),)
	return err
}

func (rel *{{ $relName }}) Create(ctx context.Context, target {{ camel $rel.Model }}N, builders ...query.SQLBuilder) (int64, error) {
	targetId, err := rel.relModel.Save(ctx, target)
	if err != nil {
		return 0, err
	}

	target.Id = null.IntFrom(targetId)

	err = rel.Attach(ctx, target)

	return targetId, err
}
`
}

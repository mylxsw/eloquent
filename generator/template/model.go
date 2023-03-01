package template

func GetModelTemplate() string {
	return `
// {{ camel $m.Name }}Model is a model which encapsulates the operations of the object
type {{ camel $m.Name }}Model struct {
	db *query.DatabaseWrap
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes []string
	
	query query.SQLBuilder
}

var {{ lower_camel $m.Name }}TableName = "{{ table $i }}"

// {{ camel $m.Name}}Table return table name for {{ camel $m.Name }}
func {{ camel $m.Name}}Table() string {
	return {{ lower_camel $m.Name }}TableName
}

const (
{{ range $j, $f := fields $m.Definition }}
	Field{{ camel $m.Name }}{{ camel $f.Name }} = "{{ snake $f.Name }}"{{ end }}
)


// {{ camel $m.Name }}Fields return all fields in {{ camel $m.Name }} model
func {{ camel $m.Name }}Fields() []string {
	return []string{ {{ range $j, $f := fields $m.Definition }}
	"{{ snake $f.Name }}",{{ end }} 
	}
}

func Set{{ camel $m.Name }}Table (tableName string) {
	{{ lower_camel $m.Name }}TableName = tableName
}

// New{{ camel $m.Name }}Model create a {{ camel $m.Name }}Model
func New{{ camel $m.Name }}Model (db query.Database) *{{ camel $m.Name }}Model {
	return &{{ camel $m.Name }}Model {
		db: query.NewDatabaseWrap(db), 
		tableName: {{ lower_camel $m.Name }}TableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes: make([]string, 0),
		query: query.Builder(),
	}
}

// GetDB return database instance
func (m *{{ camel $m.Name }}Model) GetDB() query.Database {
	return m.db.GetDB()
}

{{ if $m.Definition.SoftDelete }}
// WithTrashed force soft deleted models to appear in a result set
func (m *{{ camel $m.Name }}Model) WithTrashed() *{{ camel $m.Name }}Model {
	return m.WithoutGlobalScopes("soft_delete")
}
{{ end }}

func (m *{{ camel $m.Name }}Model) clone() *{{ camel $m.Name }}Model {
	return &{{ camel $m.Name }}Model{
		db: m.db, 
		tableName: m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes: append([]string{}, m.includeLocalScopes...),
		query: m.query,
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *{{ camel $m.Name }}Model) WithoutGlobalScopes(names ...string) *{{ camel $m.Name }}Model {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *{{ camel $m.Name }}Model) WithLocalScopes(names ...string) *{{ camel $m.Name }}Model {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Condition add query builder to model
func (m *{{ camel $m.Name }}Model) Condition(builder query.SQLBuilder) *{{ camel $m.Name }}Model {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *{{ camel $m.Name }}Model) Find(ctx context.Context, id int64) (*{{ camel $m.Name }}N, error) {
	return m.First(ctx, m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *{{ camel $m.Name }}Model) Exists(ctx context.Context, builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(ctx, builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *{{ camel $m.Name }}Model) Count(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {
	sqlStr, params := m.query.
		Merge(builders...).
		Table(m.tableName).
		AppendCondition(m.applyScope()).
		ResolveCount()
	
	rows, err := m.db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	rows.Next()
	var res int64
	if err := rows.Scan(&res); err != nil {
		return 0, err
	}

	return res, nil
}

func (m *{{ camel $m.Name }}Model) Paginate(ctx context.Context, page int64, perPage int64, builders ...query.SQLBuilder) ([]{{ camel $m.Name }}N, query.PaginateMeta, error) {
	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = 15
	}

	meta := query.PaginateMeta {
		PerPage: perPage,
		Page: page,
	}

	count, err := m.Count(ctx, builders...)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = count
	meta.LastPage = count / perPage
	if count % perPage != 0 {
		meta.LastPage += 1
	}


	res, err := m.Get(ctx, append([]query.SQLBuilder{query.Builder().Limit(perPage).Offset((page - 1) * perPage)}, builders...)...)
	if err != nil {
		return res, meta, err
	}

	return res, meta, nil
}

// Get retrieve all results for given query
func (m *{{ camel $m.Name }}Model) Get(ctx context.Context, builders ...query.SQLBuilder) ([]{{ camel $m.Name }}N, error) {
	b := m.query.Merge(builders...).Table(m.tableName).AppendCondition(m.applyScope())
	if len(b.GetFields()) == 0 {
		b = b.Select({{ range $j, $f := fields $m.Definition }}
			"{{ snake $f.Name }}",{{ end }}
		)
	}

	fields := b.GetFields()
	selectFields := make([]query.Expr, 0)

	for _, f := range fields {
		switch strcase.ToSnake(f.Value) {
		{{ range $j, $f := fields $m.Definition }} 
		case "{{ snake $f.Name }}":
			selectFields = append(selectFields, f){{ end }}
		}
	}

	var createScanVar = func(fields []query.Expr) (*{{ camel $m.Name }}N, []interface{}) {
		var {{ lower_camel $m.Name }}Var {{ camel $m.Name }}N
		scanFields := make([]interface{}, 0)

		for _, f := range fields {
			switch strcase.ToSnake(f.Value) {
			{{ range $j, $f := fields $m.Definition }} 
			case "{{ snake $f.Name }}":
				scanFields = append(scanFields, &{{ lower_camel $m.Name }}Var.{{ camel $f.Name }}){{ end }}
			}
		}

		return &{{ lower_camel $m.Name }}Var, scanFields
	}
	
	sqlStr, params := b.Fields(selectFields...).ResolveQuery()
	
	rows, err := m.db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	{{ lower_camel $m.Name }}s := make([]{{ camel $m.Name }}N, 0)
	for rows.Next() {
		{{ lower_camel $m.Name }}Real, scanFields := createScanVar(fields)
		if err := rows.Scan(scanFields...); err != nil {
			return nil, err
		}

		{{ lower_camel $m.Name }}Real.original = &{{ lower_camel $m.Name }}Original{}
		_ = query.Copy({{ lower_camel $m.Name }}Real, {{ lower_camel $m.Name }}Real.original)

		{{ lower_camel $m.Name }}Real.SetModel(m)
		{{ lower_camel $m.Name }}s = append({{ lower_camel $m.Name }}s, *{{ lower_camel $m.Name }}Real)
	}

	return {{ lower_camel $m.Name }}s, nil
}

// First return first result for given query
func (m *{{ camel $m.Name }}Model) First(ctx context.Context, builders ...query.SQLBuilder) (*{{ camel $m.Name }}N, error) {
	res, err := m.Get(ctx, append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return nil, err 
	}

	if len(res) == 0 {
		return nil, query.ErrNoResult
	}

	return &res[0], nil
}

// Create save a new {{ $m.Name }} to database
func (m *{{ camel $m.Name }}Model) Create(ctx context.Context, kv query.KV) (int64, error) {
	{{ if not $m.Definition.WithoutCreateTime }}
	if _, ok := kv["created_at"]; !ok {
		kv["created_at"] = time.Now()
	}
	{{ end }}
	{{ if not $m.Definition.WithoutUpdateTime }}
	if _, ok := kv["updated_at"]; !ok {
		kv["updated_at"] = time.Now()
	}
	{{ end }}

	sqlStr, params := m.query.Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// SaveAll save all {{ $m.Name }}s to database
func (m *{{ camel $m.Name }}Model) SaveAll(ctx context.Context, {{ lower_camel $m.Name }}s []{{ camel $m.Name }}N) ([]int64, error) {
	ids := make([]int64, 0)
	for _, {{ lower_camel $m.Name }} := range {{ lower_camel $m.Name }}s {
		id, err := m.Save(ctx, {{ lower_camel $m.Name }})
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a {{ $m.Name }} to database
func (m *{{ camel $m.Name }}Model) Save(ctx context.Context, {{ lower_camel $m.Name }} {{ camel $m.Name }}N, onlyFields ...string) (int64, error) {
	return m.Create(ctx, {{ lower_camel $m.Name }}.StaledKV(onlyFields...))
}

// SaveOrUpdate save a new {{ $m.Name }} or update it when it has a id > 0
func (m *{{ camel $m.Name }}Model) SaveOrUpdate(ctx context.Context, {{ lower_camel $m.Name }} {{ camel $m.Name }}N, onlyFields ...string) (id int64, updated bool, err error) {
	if {{ lower_camel $m.Name }}.Id.Int64 > 0 {
		_, _err := m.UpdateById(ctx, {{ lower_camel $m.Name }}.Id.Int64, {{ lower_camel $m.Name }}, onlyFields...)
		return {{ lower_camel $m.Name }}.Id.Int64, true, _err
	}

	_id, _err := m.Save(ctx, {{ lower_camel $m.Name }}, onlyFields...)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *{{ camel $m.Name }}Model) UpdateFields(ctx context.Context, kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	{{ if not $m.Definition.WithoutUpdateTime }}
	kv["updated_at"] = time.Now()
	{{ end }}

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).
		Table(m.tableName).
		ResolveUpdate(kv)

	res, err := m.db.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Update update a model for given query
func (m *{{ camel $m.Name }}Model) Update(ctx context.Context, builder query.SQLBuilder, {{ lower_camel $m.Name }} {{ camel $m.Name }}N, onlyFields ...string) (int64, error) {
	return m.UpdateFields(ctx, {{ lower_camel $m.Name }}.StaledKV(onlyFields...), builder)
}

// UpdateById update a model by id
func (m *{{ camel $m.Name }}Model) UpdateById(ctx context.Context, id int64, {{ lower_camel $m.Name }} {{ camel $m.Name }}N, onlyFields ...string) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).UpdateFields(ctx, {{ lower_camel $m.Name }}.StaledKV(onlyFields...))
}

{{ if $m.Definition.SoftDelete }}
// ForceDelete permanently remove a soft deleted model from the database
func (m *{{ camel $m.Name }}Model) ForceDelete(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()

	sqlStr, params := m2.query.Merge(builders...).AppendCondition(m2.applyScope()).Table(m2.tableName).ResolveDelete()

	res, err := m2.db.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// ForceDeleteById permanently remove a soft deleted model from the database by id
func (m *{{ camel $m.Name }}Model) ForceDeleteById(ctx context.Context, id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).ForceDelete(ctx)
}

// Restore restore a soft deleted model into an active state
func (m *{{ camel $m.Name }}Model) Restore(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()
	return m2.UpdateFields(ctx, query.KV {
		"deleted_at": nil,
	}, builders...)
}

// RestoreById restore a soft deleted model into an active state by id
func (m *{{ camel $m.Name }}Model) RestoreById(ctx context.Context, id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Restore(ctx)
}
{{ end }}

// Delete remove a model
func (m *{{ camel $m.Name }}Model) Delete(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {
	{{ if $m.Definition.SoftDelete }}
	return m.UpdateFields(ctx, query.KV {
		"deleted_at": time.Now(),
	}, builders...)
	{{ else }}
	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
	{{ end }}
}

// DeleteById remove a model by id
func (m *{{ camel $m.Name }}Model) DeleteById(ctx context.Context, id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Delete(ctx)
}
`
}

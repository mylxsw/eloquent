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

// Query add query builder to model
func (m *{{ camel $m.Name }}Model) Query(builder query.SQLBuilder) *{{ camel $m.Name }}Model {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *{{ camel $m.Name }}Model) Find(id int64) ({{ camel $m.Name }}, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *{{ camel $m.Name }}Model) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *{{ camel $m.Name }}Model) Count(builders ...query.SQLBuilder) (int64, error) {
	sqlStr, params := m.query.
		Merge(builders...).
		Table(m.tableName).
		AppendCondition(m.applyScope()).
		ResolveCount()
	
	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
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

func (m *{{ camel $m.Name }}Model) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]{{ camel $m.Name }}, query.PaginateMeta, error) {
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

	count, err := m.Count(builders...)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = count
	meta.LastPage = count / perPage
	if count % perPage != 0 {
		meta.LastPage += 1
	}


	res, err := m.Get(append([]query.SQLBuilder{query.Builder().Limit(perPage).Offset((page - 1) * perPage)}, builders...)...)
	if err != nil {
		return res, meta, err
	}

	return res, meta, nil
}

// Get retrieve all results for given query
func (m *{{ camel $m.Name }}Model) Get(builders ...query.SQLBuilder) ([]{{ camel $m.Name }}, error) {
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

	var createScanVar = func(fields []query.Expr) (*{{ lower_camel $m.Name }}Wrap, []interface{}) {
		var {{ lower_camel $m.Name }}Var {{ lower_camel $m.Name }}Wrap
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
	
	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	{{ lower_camel $m.Name }}s := make([]{{ camel $m.Name }}, 0)
	for rows.Next() {
		{{ lower_camel $m.Name }}Var, scanFields := createScanVar(fields)
		if err := rows.Scan(scanFields...); err != nil {
			return nil, err
		}

		{{ lower_camel $m.Name }}Real := {{ lower_camel $m.Name }}Var.To{{ camel $m.Name }}()
		{{ lower_camel $m.Name }}Real.SetModel(m)
		{{ lower_camel $m.Name }}s = append({{ lower_camel $m.Name }}s, {{ lower_camel $m.Name }}Real)
	}

	return {{ lower_camel $m.Name }}s, nil
}

// First return first result for given query
func (m *{{ camel $m.Name }}Model) First(builders ...query.SQLBuilder) ({{ camel $m.Name }}, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return {{ camel $m.Name }}{}, err 
	}

	if len(res) == 0 {
		return {{ camel $m.Name }}{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new {{ $m.Name }} to database
func (m *{{ camel $m.Name }}Model) Create(kv query.KV) (int64, error) {
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

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// SaveAll save all {{ $m.Name }}s to database
func (m *{{ camel $m.Name }}Model) SaveAll({{ lower_camel $m.Name }}s []{{ camel $m.Name }}) ([]int64, error) {
	ids := make([]int64, 0)
	for _, {{ lower_camel $m.Name }} := range {{ lower_camel $m.Name }}s {
		id, err := m.Save({{ lower_camel $m.Name }})
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a {{ $m.Name }} to database
func (m *{{ camel $m.Name }}Model) Save({{ lower_camel $m.Name }} {{ camel $m.Name }}) (int64, error) {
	return m.Create({{ lower_camel $m.Name }}.StaledKV())
}

// SaveOrUpdate save a new {{ $m.Name }} or update it when it has a id > 0
func (m *{{ camel $m.Name }}Model) SaveOrUpdate({{ lower_camel $m.Name }} {{ camel $m.Name }}) (id int64, updated bool, err error) {
	if {{ lower_camel $m.Name }}.Id > 0 {
		_, _err := m.UpdateById({{ lower_camel $m.Name }}.Id, {{ lower_camel $m.Name }})
		return {{ lower_camel $m.Name }}.Id, true, _err
	}

	_id, _err := m.Save({{ lower_camel $m.Name }})
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *{{ camel $m.Name }}Model) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	{{ if not $m.Definition.WithoutUpdateTime }}
	kv["updated_at"] = time.Now()
	{{ end }}

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).
		Table(m.tableName).
		ResolveUpdate(kv)

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Update update a model for given query
func (m *{{ camel $m.Name }}Model) Update({{ lower_camel $m.Name }} {{ camel $m.Name }}) (int64, error) {
	return m.UpdateFields({{ lower_camel $m.Name }}.StaledKV())
}

// UpdateById update a model by id
func (m *{{ camel $m.Name }}Model) UpdateById(id int64, {{ lower_camel $m.Name }} {{ camel $m.Name }}) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Update({{ lower_camel $m.Name }})
}

{{ if $m.Definition.SoftDelete }}
// ForceDelete permanently remove a soft deleted model from the database
func (m *{{ camel $m.Name }}Model) ForceDelete(builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()

	sqlStr, params := m2.query.Merge(builders...).AppendCondition(m2.applyScope()).Table(m2.tableName).ResolveDelete()

	res, err := m2.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// ForceDeleteById permanently remove a soft deleted model from the database by id
func (m *{{ camel $m.Name }}Model) ForceDeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).ForceDelete()
}

// Restore restore a soft deleted model into an active state
func (m *{{ camel $m.Name }}Model) Restore(builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()
	return m2.UpdateFields(query.KV {
		"deleted_at": nil,
	}, builders...)
}

// RestoreById restore a soft deleted model into an active state by id
func (m *{{ camel $m.Name }}Model) RestoreById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Restore()
}
{{ end }}

// Delete remove a model
func (m *{{ camel $m.Name }}Model) Delete(builders ...query.SQLBuilder) (int64, error) {
	{{ if $m.Definition.SoftDelete }}
	return m.UpdateFields(query.KV {
		"deleted_at": time.Now(),
	}, builders...)
	{{ else }}
	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
	{{ end }}
}

// DeleteById remove a model by id
func (m *{{ camel $m.Name }}Model) DeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Delete()
}
`
}

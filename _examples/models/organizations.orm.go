package models

// !!! DO NOT EDIT THIS FILE

import (
	"context"
	"encoding/json"
	"github.com/iancoleman/strcase"
	"github.com/mylxsw/eloquent"
	"github.com/mylxsw/eloquent/query"
	"gopkg.in/guregu/null.v3"
	"time"
)

func init() {

}

// OrganizationN is a Organization object, all fields are nullable
type OrganizationN struct {
	original          *organizationOriginal
	organizationModel *OrganizationModel

	Id        null.Int
	Name      null.String
	CreatedAt null.Time
	UpdatedAt null.Time
}

// As convert object to other type
// dst must be a pointer to struct
func (inst *OrganizationN) As(dst interface{}) error {
	return query.Copy(inst, dst)
}

// SetModel set model for Organization
func (inst *OrganizationN) SetModel(organizationModel *OrganizationModel) {
	inst.organizationModel = organizationModel
}

// organizationOriginal is an object which stores original Organization from database
type organizationOriginal struct {
	Id        null.Int
	Name      null.String
	CreatedAt null.Time
	UpdatedAt null.Time
}

// Staled identify whether the object has been modified
func (inst *OrganizationN) Staled(onlyFields ...string) bool {
	if inst.original == nil {
		inst.original = &organizationOriginal{}
	}

	if len(onlyFields) == 0 {

		if inst.Id != inst.original.Id {
			return true
		}
		if inst.Name != inst.original.Name {
			return true
		}
		if inst.CreatedAt != inst.original.CreatedAt {
			return true
		}
		if inst.UpdatedAt != inst.original.UpdatedAt {
			return true
		}
	} else {
		for _, f := range onlyFields {
			switch strcase.ToSnake(f) {

			case "id":
				if inst.Id != inst.original.Id {
					return true
				}
			case "name":
				if inst.Name != inst.original.Name {
					return true
				}
			case "created_at":
				if inst.CreatedAt != inst.original.CreatedAt {
					return true
				}
			case "updated_at":
				if inst.UpdatedAt != inst.original.UpdatedAt {
					return true
				}
			default:
			}
		}
	}

	return false
}

// StaledKV return all fields has been modified
func (inst *OrganizationN) StaledKV(onlyFields ...string) query.KV {
	kv := make(query.KV, 0)

	if inst.original == nil {
		inst.original = &organizationOriginal{}
	}

	if len(onlyFields) == 0 {

		if inst.Id != inst.original.Id {
			kv["id"] = inst.Id
		}
		if inst.Name != inst.original.Name {
			kv["name"] = inst.Name
		}
		if inst.CreatedAt != inst.original.CreatedAt {
			kv["created_at"] = inst.CreatedAt
		}
		if inst.UpdatedAt != inst.original.UpdatedAt {
			kv["updated_at"] = inst.UpdatedAt
		}
	} else {
		for _, f := range onlyFields {
			switch strcase.ToSnake(f) {

			case "id":
				if inst.Id != inst.original.Id {
					kv["id"] = inst.Id
				}
			case "name":
				if inst.Name != inst.original.Name {
					kv["name"] = inst.Name
				}
			case "created_at":
				if inst.CreatedAt != inst.original.CreatedAt {
					kv["created_at"] = inst.CreatedAt
				}
			case "updated_at":
				if inst.UpdatedAt != inst.original.UpdatedAt {
					kv["updated_at"] = inst.UpdatedAt
				}
			default:
			}
		}
	}

	return kv
}

// Save create a new model or update it
func (inst *OrganizationN) Save(ctx context.Context, onlyFields ...string) error {
	if inst.organizationModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := inst.organizationModel.SaveOrUpdate(ctx, *inst, onlyFields...)
	if err != nil {
		return err
	}

	inst.Id = null.IntFrom(id)
	return nil
}

// Delete remove a organization
func (inst *OrganizationN) Delete(ctx context.Context) error {
	if inst.organizationModel == nil {
		return query.ErrModelNotSet
	}

	_, err := inst.organizationModel.DeleteById(ctx, inst.Id.Int64)
	if err != nil {
		return err
	}

	return nil
}

// String convert instance to json string
func (inst *OrganizationN) String() string {
	rs, _ := json.Marshal(inst)
	return string(rs)
}

func (inst *OrganizationN) Users() *OrganizationBelongsToManyUserRel {
	return &OrganizationBelongsToManyUserRel{
		source:     inst,
		pivotTable: "user_organization_ref",
		relModel:   NewUserModel(inst.organizationModel.GetDB()),
	}
}

type OrganizationBelongsToManyUserRel struct {
	source     *OrganizationN
	pivotTable string
	relModel   *UserModel
}

func (rel *OrganizationBelongsToManyUserRel) Get(ctx context.Context, builders ...query.SQLBuilder) ([]UserN, error) {
	res, err := eloquent.DB(rel.relModel.GetDB()).Query(
		ctx,
		query.Builder().Table(rel.pivotTable).Select("user_id").Where("organization_id", rel.source.Id),
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

	return rel.relModel.Get(ctx, query.Builder().Merge(builders...).WhereIn("id", res...))
}

func (rel *OrganizationBelongsToManyUserRel) Count(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {
	res, err := eloquent.DB(rel.relModel.GetDB()).Query(
		ctx,
		query.Builder().Table(rel.pivotTable).Select(query.Raw("COUNT(1) as c")).Where("organization_id", rel.source.Id),
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

	return res[0].(int64), nil
}

func (rel *OrganizationBelongsToManyUserRel) Exists(ctx context.Context, builders ...query.SQLBuilder) (bool, error) {
	c, err := rel.Count(ctx, builders...)
	if err != nil {
		return false, err
	}

	return c > 0, nil
}

func (rel *OrganizationBelongsToManyUserRel) Attach(ctx context.Context, target UserN) error {
	_, err := eloquent.DB(rel.relModel.GetDB()).Insert(ctx, rel.pivotTable, query.KV{
		"user_id":         target.Id,
		"organization_id": rel.source.Id,
	})

	return err
}

func (rel *OrganizationBelongsToManyUserRel) Detach(ctx context.Context, target UserN) error {
	_, err := eloquent.DB(rel.relModel.GetDB()).
		Delete(ctx, eloquent.Build(rel.pivotTable).
			Where("user_id", target.Id).
			Where("organization_id", rel.source.Id))

	return err
}

func (rel *OrganizationBelongsToManyUserRel) DetachAll(ctx context.Context) error {
	_, err := eloquent.DB(rel.relModel.GetDB()).
		Delete(ctx, eloquent.Build(rel.pivotTable).
			Where("organization_id", rel.source.Id))
	return err
}

func (rel *OrganizationBelongsToManyUserRel) Create(ctx context.Context, target UserN, builders ...query.SQLBuilder) (int64, error) {
	targetId, err := rel.relModel.Save(ctx, target)
	if err != nil {
		return 0, err
	}

	target.Id = null.IntFrom(targetId)

	err = rel.Attach(ctx, target)

	return targetId, err
}

type organizationScope struct {
	name  string
	apply func(builder query.Condition)
}

var organizationGlobalScopes = make([]organizationScope, 0)
var organizationLocalScopes = make([]organizationScope, 0)

// AddGlobalScopeForOrganization assign a global scope to a model
func AddGlobalScopeForOrganization(name string, apply func(builder query.Condition)) {
	organizationGlobalScopes = append(organizationGlobalScopes, organizationScope{name: name, apply: apply})
}

// AddLocalScopeForOrganization assign a local scope to a model
func AddLocalScopeForOrganization(name string, apply func(builder query.Condition)) {
	organizationLocalScopes = append(organizationLocalScopes, organizationScope{name: name, apply: apply})
}

func (m *OrganizationModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range organizationGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range organizationLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *OrganizationModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *OrganizationModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type Organization struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w Organization) ToOrganizationN(allows ...string) OrganizationN {
	if len(allows) == 0 {
		return OrganizationN{

			Id:        null.IntFrom(int64(w.Id)),
			Name:      null.StringFrom(w.Name),
			CreatedAt: null.TimeFrom(w.CreatedAt),
			UpdatedAt: null.TimeFrom(w.UpdatedAt),
		}
	}

	res := OrganizationN{}
	for _, al := range allows {
		switch strcase.ToSnake(al) {

		case "id":
			res.Id = null.IntFrom(int64(w.Id))
		case "name":
			res.Name = null.StringFrom(w.Name)
		case "created_at":
			res.CreatedAt = null.TimeFrom(w.CreatedAt)
		case "updated_at":
			res.UpdatedAt = null.TimeFrom(w.UpdatedAt)
		default:
		}
	}

	return res
}

// As convert object to other type
// dst must be a pointer to struct
func (w Organization) As(dst interface{}) error {
	return query.Copy(w, dst)
}

func (w *OrganizationN) ToOrganization() Organization {
	return Organization{

		Id:        w.Id.Int64,
		Name:      w.Name.String,
		CreatedAt: w.CreatedAt.Time,
		UpdatedAt: w.UpdatedAt.Time,
	}
}

// OrganizationModel is a model which encapsulates the operations of the object
type OrganizationModel struct {
	db        *query.DatabaseWrap
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string

	query query.SQLBuilder
}

var organizationTableName = "wz_organization"

// OrganizationTable return table name for Organization
func OrganizationTable() string {
	return organizationTableName
}

const (
	FieldOrganizationId        = "id"
	FieldOrganizationName      = "name"
	FieldOrganizationCreatedAt = "created_at"
	FieldOrganizationUpdatedAt = "updated_at"
)

// OrganizationFields return all fields in Organization model
func OrganizationFields() []string {
	return []string{
		"id",
		"name",
		"created_at",
		"updated_at",
	}
}

func SetOrganizationTable(tableName string) {
	organizationTableName = tableName
}

// NewOrganizationModel create a OrganizationModel
func NewOrganizationModel(db query.Database) *OrganizationModel {
	return &OrganizationModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           organizationTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *OrganizationModel) GetDB() query.Database {
	return m.db.GetDB()
}

func (m *OrganizationModel) clone() *OrganizationModel {
	return &OrganizationModel{
		db:                  m.db,
		tableName:           m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes:  append([]string{}, m.includeLocalScopes...),
		query:               m.query,
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *OrganizationModel) WithoutGlobalScopes(names ...string) *OrganizationModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *OrganizationModel) WithLocalScopes(names ...string) *OrganizationModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Condition add query builder to model
func (m *OrganizationModel) Condition(builder query.SQLBuilder) *OrganizationModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *OrganizationModel) Find(ctx context.Context, id int64) (*OrganizationN, error) {
	return m.First(ctx, m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *OrganizationModel) Exists(ctx context.Context, builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(ctx, builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *OrganizationModel) Count(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {
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

func (m *OrganizationModel) Paginate(ctx context.Context, page int64, perPage int64, builders ...query.SQLBuilder) ([]OrganizationN, query.PaginateMeta, error) {
	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = 15
	}

	meta := query.PaginateMeta{
		PerPage: perPage,
		Page:    page,
	}

	count, err := m.Count(ctx, builders...)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = count
	meta.LastPage = count / perPage
	if count%perPage != 0 {
		meta.LastPage += 1
	}

	res, err := m.Get(ctx, append([]query.SQLBuilder{query.Builder().Limit(perPage).Offset((page - 1) * perPage)}, builders...)...)
	if err != nil {
		return res, meta, err
	}

	return res, meta, nil
}

// Get retrieve all results for given query
func (m *OrganizationModel) Get(ctx context.Context, builders ...query.SQLBuilder) ([]OrganizationN, error) {
	b := m.query.Merge(builders...).Table(m.tableName).AppendCondition(m.applyScope())
	if len(b.GetFields()) == 0 {
		b = b.Select(
			"id",
			"name",
			"created_at",
			"updated_at",
		)
	}

	fields := b.GetFields()
	selectFields := make([]query.Expr, 0)

	for _, f := range fields {
		switch strcase.ToSnake(f.Value) {

		case "id":
			selectFields = append(selectFields, f)
		case "name":
			selectFields = append(selectFields, f)
		case "created_at":
			selectFields = append(selectFields, f)
		case "updated_at":
			selectFields = append(selectFields, f)
		}
	}

	var createScanVar = func(fields []query.Expr) (*OrganizationN, []interface{}) {
		var organizationVar OrganizationN
		scanFields := make([]interface{}, 0)

		for _, f := range fields {
			switch strcase.ToSnake(f.Value) {

			case "id":
				scanFields = append(scanFields, &organizationVar.Id)
			case "name":
				scanFields = append(scanFields, &organizationVar.Name)
			case "created_at":
				scanFields = append(scanFields, &organizationVar.CreatedAt)
			case "updated_at":
				scanFields = append(scanFields, &organizationVar.UpdatedAt)
			}
		}

		return &organizationVar, scanFields
	}

	sqlStr, params := b.Fields(selectFields...).ResolveQuery()

	rows, err := m.db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	organizations := make([]OrganizationN, 0)
	for rows.Next() {
		organizationReal, scanFields := createScanVar(fields)
		if err := rows.Scan(scanFields...); err != nil {
			return nil, err
		}

		organizationReal.original = &organizationOriginal{}
		_ = query.Copy(organizationReal, organizationReal.original)

		organizationReal.SetModel(m)
		organizations = append(organizations, *organizationReal)
	}

	return organizations, nil
}

// First return first result for given query
func (m *OrganizationModel) First(ctx context.Context, builders ...query.SQLBuilder) (*OrganizationN, error) {
	res, err := m.Get(ctx, append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, query.ErrNoResult
	}

	return &res[0], nil
}

// Create save a new organization to database
func (m *OrganizationModel) Create(ctx context.Context, kv query.KV) (int64, error) {

	if _, ok := kv["created_at"]; !ok {
		kv["created_at"] = time.Now()
	}

	if _, ok := kv["updated_at"]; !ok {
		kv["updated_at"] = time.Now()
	}

	sqlStr, params := m.query.Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// SaveAll save all organizations to database
func (m *OrganizationModel) SaveAll(ctx context.Context, organizations []OrganizationN) ([]int64, error) {
	ids := make([]int64, 0)
	for _, organization := range organizations {
		id, err := m.Save(ctx, organization)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a organization to database
func (m *OrganizationModel) Save(ctx context.Context, organization OrganizationN, onlyFields ...string) (int64, error) {
	return m.Create(ctx, organization.StaledKV(onlyFields...))
}

// SaveOrUpdate save a new organization or update it when it has a id > 0
func (m *OrganizationModel) SaveOrUpdate(ctx context.Context, organization OrganizationN, onlyFields ...string) (id int64, updated bool, err error) {
	if organization.Id.Int64 > 0 {
		_, _err := m.UpdateById(ctx, organization.Id.Int64, organization, onlyFields...)
		return organization.Id.Int64, true, _err
	}

	_id, _err := m.Save(ctx, organization, onlyFields...)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *OrganizationModel) UpdateFields(ctx context.Context, kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	kv["updated_at"] = time.Now()

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
func (m *OrganizationModel) Update(ctx context.Context, builder query.SQLBuilder, organization OrganizationN, onlyFields ...string) (int64, error) {
	return m.UpdateFields(ctx, organization.StaledKV(onlyFields...), builder)
}

// UpdateById update a model by id
func (m *OrganizationModel) UpdateById(ctx context.Context, id int64, organization OrganizationN, onlyFields ...string) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).UpdateFields(ctx, organization.StaledKV(onlyFields...))
}

// Delete remove a model
func (m *OrganizationModel) Delete(ctx context.Context, builders ...query.SQLBuilder) (int64, error) {

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

// DeleteById remove a model by id
func (m *OrganizationModel) DeleteById(ctx context.Context, id int64) (int64, error) {
	return m.Condition(query.Builder().Where("id", "=", id)).Delete(ctx)
}

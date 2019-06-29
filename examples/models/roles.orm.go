package models

// !!! DO NOT EDIT THIS FILE

import (
	"context"
	"github.com/mylxsw/eloquent/query"
	"gopkg.in/guregu/null.v3"
	"time"
)

func init() {

}

// Role is a Role object
type Role struct {
	original  *roleOriginal
	roleModel *RoleModel

	Name        string
	Description string
	Id          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SetModel set model for Role
func (roleSelf *Role) SetModel(roleModel *RoleModel) {
	roleSelf.roleModel = roleModel
}

// roleOriginal is an object which stores original Role from database
type roleOriginal struct {
	Name        string
	Description string
	Id          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Staled identify whether the object has been modified
func (roleSelf *Role) Staled() bool {
	if roleSelf.original == nil {
		roleSelf.original = &roleOriginal{}
	}

	if roleSelf.Name != roleSelf.original.Name {
		return true
	}
	if roleSelf.Description != roleSelf.original.Description {
		return true
	}
	if roleSelf.Id != roleSelf.original.Id {
		return true
	}
	if roleSelf.CreatedAt != roleSelf.original.CreatedAt {
		return true
	}
	if roleSelf.UpdatedAt != roleSelf.original.UpdatedAt {
		return true
	}

	return false
}

// StaledKV return all fields has been modified
func (roleSelf *Role) StaledKV() query.KV {
	kv := make(query.KV, 0)

	if roleSelf.original == nil {
		roleSelf.original = &roleOriginal{}
	}

	if roleSelf.Name != roleSelf.original.Name {
		kv["name"] = roleSelf.Name
	}
	if roleSelf.Description != roleSelf.original.Description {
		kv["description"] = roleSelf.Description
	}
	if roleSelf.Id != roleSelf.original.Id {
		kv["id"] = roleSelf.Id
	}
	if roleSelf.CreatedAt != roleSelf.original.CreatedAt {
		kv["created_at"] = roleSelf.CreatedAt
	}
	if roleSelf.UpdatedAt != roleSelf.original.UpdatedAt {
		kv["updated_at"] = roleSelf.UpdatedAt
	}

	return kv
}

func (roleSelf *Role) Users() *RoleHasManyUserRel {
	return &RoleHasManyUserRel{
		source:   roleSelf,
		relModel: NewUserModel(roleSelf.roleModel.GetDB()),
	}
}

type RoleHasManyUserRel struct {
	source   *Role
	relModel *UserModel
}

func (rel *RoleHasManyUserRel) Get(builders ...query.SQLBuilder) ([]User, error) {
	builder := query.Builder().Where("role_id", rel.source.Id).Merge(builders...)

	return rel.relModel.Get(builder)
}

func (rel *RoleHasManyUserRel) First(builders ...query.SQLBuilder) (User, error) {
	builder := query.Builder().Where("role_id", rel.source.Id).Limit(1).Merge(builders...)
	return rel.relModel.First(builder)
}

func (rel *RoleHasManyUserRel) Create(target User) (int64, error) {
	target.RoleId = rel.source.Id
	return rel.relModel.Save(target)
}

// Save create a new model or update it
func (roleSelf *Role) Save() error {
	if roleSelf.roleModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := roleSelf.roleModel.SaveOrUpdate(*roleSelf)
	if err != nil {
		return err
	}

	roleSelf.Id = id
	return nil
}

// Delete remove a Role
func (roleSelf *Role) Delete() error {
	if roleSelf.roleModel == nil {
		return query.ErrModelNotSet
	}

	_, err := roleSelf.roleModel.DeleteById(roleSelf.Id)
	if err != nil {
		return err
	}

	return nil
}

type roleScope struct {
	name  string
	apply func(builder query.Condition)
}

var roleGlobalScopes = make([]roleScope, 0)
var roleLocalScopes = make([]roleScope, 0)

// AddRoleGlobalScope assign a global scope to a model
func AddRoleGlobalScope(name string, apply func(builder query.Condition)) {
	roleGlobalScopes = append(roleGlobalScopes, roleScope{name: name, apply: apply})
}

// AddRoleLocalScope assign a local scope to a model
func AddRoleLocalScope(name string, apply func(builder query.Condition)) {
	roleLocalScopes = append(roleLocalScopes, roleScope{name: name, apply: apply})
}

func (m *RoleModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range roleGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range roleLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *RoleModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *RoleModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type roleWrap struct {
	Name        null.String
	Description null.String
	Id          null.Int
	CreatedAt   null.Time
	UpdatedAt   null.Time
}

func (w roleWrap) ToRole() Role {
	return Role{
		original: &roleOriginal{
			Name:        w.Name.String,
			Description: w.Description.String,
			Id:          w.Id.Int64,
			CreatedAt:   w.CreatedAt.Time,
			UpdatedAt:   w.UpdatedAt.Time,
		},

		Name:        w.Name.String,
		Description: w.Description.String,
		Id:          w.Id.Int64,
		CreatedAt:   w.CreatedAt.Time,
		UpdatedAt:   w.UpdatedAt.Time,
	}
}

// RoleModel is a model which encapsulates the operations of the object
type RoleModel struct {
	db        *query.DatabaseWrap
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string

	query query.SQLBuilder

	beforeCreate func(kv query.KV) error
	afterCreate  func(id int64) error
	beforeUpdate func(kv query.KV) error
	beforeDelete func() error
	afterDelete  func() error
}

func (m *RoleModel) BeforeCreate(f func(kv query.KV) error) {
	m.beforeCreate = f
}

func (m *RoleModel) AfterCreate(f func(id int64) error) {
	m.afterCreate = f
}

func (m *RoleModel) BeforeUpdate(f func(kv query.KV) error) {
	m.beforeUpdate = f
}

func (m *RoleModel) BeforeDelete(f func() error) {
	m.beforeDelete = f
}

func (m *RoleModel) AfterDelete(f func() error) {
	m.afterDelete = f
}

var roleTableName = "wz_role"

func SetRoleTable(tableName string) {
	roleTableName = tableName
}

// NewRoleModel create a RoleModel
func NewRoleModel(db query.Database) *RoleModel {
	return &RoleModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           roleTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *RoleModel) GetDB() query.Database {
	return m.db.GetDB()
}

func (m *RoleModel) clone() *RoleModel {
	return &RoleModel{
		db:                  m.db,
		tableName:           m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes:  append([]string{}, m.includeLocalScopes...),
		query:               m.query,
		beforeCreate:        m.beforeCreate,
		afterCreate:         m.afterCreate,
		beforeUpdate:        m.beforeUpdate,
		beforeDelete:        m.beforeDelete,
		afterDelete:         m.afterDelete,
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *RoleModel) WithoutGlobalScopes(names ...string) *RoleModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *RoleModel) WithLocalScopes(names ...string) *RoleModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Query add query builder to model
func (m *RoleModel) Query(builder query.SQLBuilder) *RoleModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *RoleModel) Find(id int64) (Role, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *RoleModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *RoleModel) Count(builders ...query.SQLBuilder) (int64, error) {
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

func (m *RoleModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]Role, query.PaginateMeta, error) {
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

	count, err := m.Count(builders...)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = count
	meta.LastPage = count / perPage
	if count%perPage != 0 {
		meta.LastPage += 1
	}

	res, err := m.Get(append([]query.SQLBuilder{query.Builder().Limit(perPage).Offset((page - 1) * perPage)}, builders...)...)
	if err != nil {
		return res, meta, err
	}

	return res, meta, nil
}

// Get retrieve all results for given query
func (m *RoleModel) Get(builders ...query.SQLBuilder) ([]Role, error) {
	sqlStr, params := m.query.Merge(builders...).
		Table(m.tableName).
		Select(
			"name",
			"description",
			"id",
			"created_at",
			"updated_at",
		).AppendCondition(m.applyScope()).
		ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := make([]Role, 0)
	for rows.Next() {
		var roleVar roleWrap
		if err := rows.Scan(
			&roleVar.Name,
			&roleVar.Description,
			&roleVar.Id,
			&roleVar.CreatedAt,
			&roleVar.UpdatedAt); err != nil {
			return nil, err
		}

		roleReal := roleVar.ToRole()
		roleReal.SetModel(m)
		roles = append(roles, roleReal)
	}

	return roles, nil
}

// First return first result for given query
func (m *RoleModel) First(builders ...query.SQLBuilder) (Role, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return Role{}, err
	}

	if len(res) == 0 {
		return Role{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new Role to database
func (m *RoleModel) Create(kv query.KV) (int64, error) {
	kv["created_at"] = time.Now()
	kv["updated_at"] = time.Now()

	if m.beforeCreate != nil {
		if err := m.beforeCreate(kv); err != nil {
			return 0, err
		}
	}

	sqlStr, params := m.query.Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if m.afterCreate != nil {
		if err := m.afterCreate(lastInsertId); err != nil {
			return lastInsertId, err
		}
	}

	return lastInsertId, nil
}

// SaveAll save all Roles to database
func (m *RoleModel) SaveAll(roles []Role) ([]int64, error) {
	ids := make([]int64, 0)
	for _, role := range roles {
		id, err := m.Save(role)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a Role to database
func (m *RoleModel) Save(role Role) (int64, error) {
	return m.Create(query.KV{
		"name":        role.Name,
		"description": role.Description,
	})
}

// SaveOrUpdate save a new Role or update it when it has a id > 0
func (m *RoleModel) SaveOrUpdate(role Role) (id int64, updated bool, err error) {
	if role.Id > 0 {
		_, _err := m.UpdateById(role.Id, role)
		return role.Id, true, _err
	}

	_id, _err := m.Save(role)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *RoleModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	kv["updated_at"] = time.Now()

	if m.beforeUpdate != nil {
		if err := m.beforeUpdate(kv); err != nil {
			return 0, err
		}
	}

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
func (m *RoleModel) Update(role Role) (int64, error) {
	return m.UpdateFields(role.StaledKV())
}

// UpdateById update a model by id
func (m *RoleModel) UpdateById(id int64, role Role) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Update(role)
}

// Delete remove a model
func (m *RoleModel) Delete(builders ...query.SQLBuilder) (int64, error) {
	if m.beforeDelete != nil {
		if err := m.beforeDelete(); err != nil {
			return 0, err
		}
	}

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return affectedRows, err
	}

	if m.afterDelete != nil {
		if err := m.afterDelete(); err != nil {
			return affectedRows, err
		}
	}

	return affectedRows, nil

}

// DeleteById remove a model by id
func (m *RoleModel) DeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Delete()
}
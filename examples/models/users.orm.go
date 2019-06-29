package models

// !!! DO NOT EDIT THIS FILE

import (
	"context"
	"github.com/mylxsw/eloquent/query"
	"gopkg.in/guregu/null.v3"
	"time"
)

func init() {

	// AddUserGlobalScope assign a global scope to a model for soft delete
	AddUserGlobalScope("soft_delete", func(builder query.Condition) {
		builder.WhereNull("deleted_at")
	})

}

// User is a User object
type User struct {
	original  *userOriginal
	userModel *UserModel

	Id            int64 `json:"id"`
	Name          string
	Email         string `json:"email"`
	Password      string `json:"password" yaml:"password"`
	RoleId        int64
	RememberToken string `json:"remember_token" yaml:"remember_token"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     null.Time
}

// SetModel set model for User
func (userSelf *User) SetModel(userModel *UserModel) {
	userSelf.userModel = userModel
}

// userOriginal is an object which stores original User from database
type userOriginal struct {
	Id            int64
	Name          string
	Email         string
	Password      string
	RoleId        int64
	RememberToken string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     null.Time
}

// Staled identify whether the object has been modified
func (userSelf *User) Staled() bool {
	if userSelf.original == nil {
		userSelf.original = &userOriginal{}
	}

	if userSelf.Id != userSelf.original.Id {
		return true
	}
	if userSelf.Name != userSelf.original.Name {
		return true
	}
	if userSelf.Email != userSelf.original.Email {
		return true
	}
	if userSelf.Password != userSelf.original.Password {
		return true
	}
	if userSelf.RoleId != userSelf.original.RoleId {
		return true
	}
	if userSelf.RememberToken != userSelf.original.RememberToken {
		return true
	}
	if userSelf.CreatedAt != userSelf.original.CreatedAt {
		return true
	}
	if userSelf.UpdatedAt != userSelf.original.UpdatedAt {
		return true
	}
	if userSelf.DeletedAt != userSelf.original.DeletedAt {
		return true
	}

	return false
}

// StaledKV return all fields has been modified
func (userSelf *User) StaledKV() query.KV {
	kv := make(query.KV, 0)

	if userSelf.original == nil {
		userSelf.original = &userOriginal{}
	}

	if userSelf.Id != userSelf.original.Id {
		kv["id"] = userSelf.Id
	}
	if userSelf.Name != userSelf.original.Name {
		kv["name"] = userSelf.Name
	}
	if userSelf.Email != userSelf.original.Email {
		kv["email"] = userSelf.Email
	}
	if userSelf.Password != userSelf.original.Password {
		kv["password"] = userSelf.Password
	}
	if userSelf.RoleId != userSelf.original.RoleId {
		kv["role_id"] = userSelf.RoleId
	}
	if userSelf.RememberToken != userSelf.original.RememberToken {
		kv["remember_token"] = userSelf.RememberToken
	}
	if userSelf.CreatedAt != userSelf.original.CreatedAt {
		kv["created_at"] = userSelf.CreatedAt
	}
	if userSelf.UpdatedAt != userSelf.original.UpdatedAt {
		kv["updated_at"] = userSelf.UpdatedAt
	}
	if userSelf.DeletedAt != userSelf.original.DeletedAt {
		kv["deleted_at"] = userSelf.DeletedAt
	}

	return kv
}

func (userSelf *User) Role() *UserBelongsToRoleRel {
	return &UserBelongsToRoleRel{
		source:   userSelf,
		relModel: NewRoleModel(userSelf.userModel.GetDB()),
	}
}

type UserBelongsToRoleRel struct {
	source   *User
	relModel *RoleModel
}

func (rel *UserBelongsToRoleRel) Create(target Role) (int64, error) {
	targetId, err := rel.relModel.Save(target)
	if err != nil {
		return 0, err
	}

	target.Id = targetId

	rel.source.RoleId = target.Id
	if err := rel.source.Save(); err != nil {
		return targetId, err
	}

	return targetId, nil
}

func (rel *UserBelongsToRoleRel) Get(builders ...query.SQLBuilder) ([]Role, error) {
	builder := query.Builder().Where("id", rel.source.RoleId).Merge(builders...)

	return rel.relModel.Get(builder)
}

func (rel *UserBelongsToRoleRel) First(builders ...query.SQLBuilder) (Role, error) {
	builder := query.Builder().Where("id", rel.source.RoleId).Limit(1).Merge(builders...)

	return rel.relModel.First(builder)
}

func (rel *UserBelongsToRoleRel) Associate(target Role) error {
	rel.source.RoleId = target.Id
	return rel.source.Save()
}

func (rel *UserBelongsToRoleRel) Dissociate() error {
	rel.source.RoleId = 0
	return rel.source.Save()
}

// Save create a new model or update it
func (userSelf *User) Save() error {
	if userSelf.userModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := userSelf.userModel.SaveOrUpdate(*userSelf)
	if err != nil {
		return err
	}

	userSelf.Id = id
	return nil
}

// Delete remove a User
func (userSelf *User) Delete() error {
	if userSelf.userModel == nil {
		return query.ErrModelNotSet
	}

	_, err := userSelf.userModel.DeleteById(userSelf.Id)
	if err != nil {
		return err
	}

	return nil
}

type userScope struct {
	name  string
	apply func(builder query.Condition)
}

var userGlobalScopes = make([]userScope, 0)
var userLocalScopes = make([]userScope, 0)

// AddUserGlobalScope assign a global scope to a model
func AddUserGlobalScope(name string, apply func(builder query.Condition)) {
	userGlobalScopes = append(userGlobalScopes, userScope{name: name, apply: apply})
}

// AddUserLocalScope assign a local scope to a model
func AddUserLocalScope(name string, apply func(builder query.Condition)) {
	userLocalScopes = append(userLocalScopes, userScope{name: name, apply: apply})
}

func (m *UserModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range userGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range userLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *UserModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *UserModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type userWrap struct {
	Id            null.Int
	Name          null.String
	Email         null.String
	Password      null.String
	RoleId        null.Int
	RememberToken null.String
	CreatedAt     null.Time
	UpdatedAt     null.Time
	DeletedAt     null.Time
}

func (w userWrap) ToUser() User {
	return User{
		original: &userOriginal{
			Id:            w.Id.Int64,
			Name:          w.Name.String,
			Email:         w.Email.String,
			Password:      w.Password.String,
			RoleId:        w.RoleId.Int64,
			RememberToken: w.RememberToken.String,
			CreatedAt:     w.CreatedAt.Time,
			UpdatedAt:     w.UpdatedAt.Time,
			DeletedAt:     w.DeletedAt,
		},

		Id:            w.Id.Int64,
		Name:          w.Name.String,
		Email:         w.Email.String,
		Password:      w.Password.String,
		RoleId:        w.RoleId.Int64,
		RememberToken: w.RememberToken.String,
		CreatedAt:     w.CreatedAt.Time,
		UpdatedAt:     w.UpdatedAt.Time,
		DeletedAt:     w.DeletedAt,
	}
}

// UserModel is a model which encapsulates the operations of the object
type UserModel struct {
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

func (m *UserModel) BeforeCreate(f func(kv query.KV) error) {
	m.beforeCreate = f
}

func (m *UserModel) AfterCreate(f func(id int64) error) {
	m.afterCreate = f
}

func (m *UserModel) BeforeUpdate(f func(kv query.KV) error) {
	m.beforeUpdate = f
}

func (m *UserModel) BeforeDelete(f func() error) {
	m.beforeDelete = f
}

func (m *UserModel) AfterDelete(f func() error) {
	m.afterDelete = f
}

var userTableName = "wz_user"

func SetUserTable(tableName string) {
	userTableName = tableName
}

// NewUserModel create a UserModel
func NewUserModel(db query.Database) *UserModel {
	return &UserModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           userTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *UserModel) GetDB() query.Database {
	return m.db.GetDB()
}

// WithTrashed force soft deleted models to appear in a result set
func (m *UserModel) WithTrashed() *UserModel {
	return m.WithoutGlobalScopes("soft_delete")
}

func (m *UserModel) clone() *UserModel {
	return &UserModel{
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
func (m *UserModel) WithoutGlobalScopes(names ...string) *UserModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *UserModel) WithLocalScopes(names ...string) *UserModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Query add query builder to model
func (m *UserModel) Query(builder query.SQLBuilder) *UserModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *UserModel) Find(id int64) (User, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *UserModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *UserModel) Count(builders ...query.SQLBuilder) (int64, error) {
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

func (m *UserModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]User, query.PaginateMeta, error) {
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
func (m *UserModel) Get(builders ...query.SQLBuilder) ([]User, error) {
	sqlStr, params := m.query.Merge(builders...).
		Table(m.tableName).
		Select(
			"id",
			"name",
			"email",
			"password",
			"role_id",
			"remember_token",
			"created_at",
			"updated_at",
			"deleted_at",
		).AppendCondition(m.applyScope()).
		ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var userVar userWrap
		if err := rows.Scan(
			&userVar.Id,
			&userVar.Name,
			&userVar.Email,
			&userVar.Password,
			&userVar.RoleId,
			&userVar.RememberToken,
			&userVar.CreatedAt,
			&userVar.UpdatedAt,
			&userVar.DeletedAt); err != nil {
			return nil, err
		}

		userReal := userVar.ToUser()
		userReal.SetModel(m)
		users = append(users, userReal)
	}

	return users, nil
}

// First return first result for given query
func (m *UserModel) First(builders ...query.SQLBuilder) (User, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return User{}, err
	}

	if len(res) == 0 {
		return User{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new User to database
func (m *UserModel) Create(kv query.KV) (int64, error) {
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

// SaveAll save all Users to database
func (m *UserModel) SaveAll(users []User) ([]int64, error) {
	ids := make([]int64, 0)
	for _, user := range users {
		id, err := m.Save(user)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a User to database
func (m *UserModel) Save(user User) (int64, error) {
	return m.Create(query.KV{
		"name":           user.Name,
		"email":          user.Email,
		"password":       user.Password,
		"role_id":        user.RoleId,
		"remember_token": user.RememberToken,
	})
}

// SaveOrUpdate save a new User or update it when it has a id > 0
func (m *UserModel) SaveOrUpdate(user User) (id int64, updated bool, err error) {
	if user.Id > 0 {
		_, _err := m.UpdateById(user.Id, user)
		return user.Id, true, _err
	}

	_id, _err := m.Save(user)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *UserModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
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
func (m *UserModel) Update(user User) (int64, error) {
	return m.UpdateFields(user.StaledKV())
}

// UpdateById update a model by id
func (m *UserModel) UpdateById(id int64, user User) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Update(user)
}

// ForceDelete permanently remove a soft deleted model from the database
func (m *UserModel) ForceDelete(builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()

	sqlStr, params := m2.query.Merge(builders...).AppendCondition(m2.applyScope()).Table(m2.tableName).ResolveDelete()

	res, err := m2.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// ForceDeleteById permanently remove a soft deleted model from the database by id
func (m *UserModel) ForceDeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).ForceDelete()
}

// Restore restore a soft deleted model into an active state
func (m *UserModel) Restore(builders ...query.SQLBuilder) (int64, error) {
	m2 := m.WithTrashed()
	return m2.UpdateFields(query.KV{
		"deleted_at": nil,
	}, builders...)
}

// RestoreById restore a soft deleted model into an active state by id
func (m *UserModel) RestoreById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Restore()
}

// Delete remove a model
func (m *UserModel) Delete(builders ...query.SQLBuilder) (int64, error) {
	if m.beforeDelete != nil {
		if err := m.beforeDelete(); err != nil {
			return 0, err
		}
	}

	affectedRows, err := m.UpdateFields(query.KV{
		"deleted_at": time.Now(),
	}, builders...)

	if err == nil && m.afterDelete != nil {
		if err2 := m.afterDelete(); err2 != nil {
			return 0, err2
		}
	}

	return affectedRows, err

}

// DeleteById remove a model by id
func (m *UserModel) DeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Delete()
}

// PasswordReset is a PasswordReset object
type PasswordReset struct {
	original           *passwordresetOriginal
	passwordresetModel *PasswordResetModel

	Email     string
	Token     string
	Id        int64
	CreatedAt time.Time
}

// SetModel set model for PasswordReset
func (passwordresetSelf *PasswordReset) SetModel(passwordresetModel *PasswordResetModel) {
	passwordresetSelf.passwordresetModel = passwordresetModel
}

// passwordresetOriginal is an object which stores original PasswordReset from database
type passwordresetOriginal struct {
	Email     string
	Token     string
	Id        int64
	CreatedAt time.Time
}

// Staled identify whether the object has been modified
func (passwordresetSelf *PasswordReset) Staled() bool {
	if passwordresetSelf.original == nil {
		passwordresetSelf.original = &passwordresetOriginal{}
	}

	if passwordresetSelf.Email != passwordresetSelf.original.Email {
		return true
	}
	if passwordresetSelf.Token != passwordresetSelf.original.Token {
		return true
	}
	if passwordresetSelf.Id != passwordresetSelf.original.Id {
		return true
	}
	if passwordresetSelf.CreatedAt != passwordresetSelf.original.CreatedAt {
		return true
	}

	return false
}

// StaledKV return all fields has been modified
func (passwordresetSelf *PasswordReset) StaledKV() query.KV {
	kv := make(query.KV, 0)

	if passwordresetSelf.original == nil {
		passwordresetSelf.original = &passwordresetOriginal{}
	}

	if passwordresetSelf.Email != passwordresetSelf.original.Email {
		kv["email"] = passwordresetSelf.Email
	}
	if passwordresetSelf.Token != passwordresetSelf.original.Token {
		kv["token"] = passwordresetSelf.Token
	}
	if passwordresetSelf.Id != passwordresetSelf.original.Id {
		kv["id"] = passwordresetSelf.Id
	}
	if passwordresetSelf.CreatedAt != passwordresetSelf.original.CreatedAt {
		kv["created_at"] = passwordresetSelf.CreatedAt
	}

	return kv
}

// Save create a new model or update it
func (passwordresetSelf *PasswordReset) Save() error {
	if passwordresetSelf.passwordresetModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := passwordresetSelf.passwordresetModel.SaveOrUpdate(*passwordresetSelf)
	if err != nil {
		return err
	}

	passwordresetSelf.Id = id
	return nil
}

// Delete remove a PasswordReset
func (passwordresetSelf *PasswordReset) Delete() error {
	if passwordresetSelf.passwordresetModel == nil {
		return query.ErrModelNotSet
	}

	_, err := passwordresetSelf.passwordresetModel.DeleteById(passwordresetSelf.Id)
	if err != nil {
		return err
	}

	return nil
}

type passwordresetScope struct {
	name  string
	apply func(builder query.Condition)
}

var passwordresetGlobalScopes = make([]passwordresetScope, 0)
var passwordresetLocalScopes = make([]passwordresetScope, 0)

// AddPasswordResetGlobalScope assign a global scope to a model
func AddPasswordResetGlobalScope(name string, apply func(builder query.Condition)) {
	passwordresetGlobalScopes = append(passwordresetGlobalScopes, passwordresetScope{name: name, apply: apply})
}

// AddPasswordResetLocalScope assign a local scope to a model
func AddPasswordResetLocalScope(name string, apply func(builder query.Condition)) {
	passwordresetLocalScopes = append(passwordresetLocalScopes, passwordresetScope{name: name, apply: apply})
}

func (m *PasswordResetModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range passwordresetGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range passwordresetLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *PasswordResetModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *PasswordResetModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type passwordResetWrap struct {
	Email     null.String
	Token     null.String
	Id        null.Int
	CreatedAt null.Time
}

func (w passwordResetWrap) ToPasswordReset() PasswordReset {
	return PasswordReset{
		original: &passwordresetOriginal{
			Email:     w.Email.String,
			Token:     w.Token.String,
			Id:        w.Id.Int64,
			CreatedAt: w.CreatedAt.Time,
		},

		Email:     w.Email.String,
		Token:     w.Token.String,
		Id:        w.Id.Int64,
		CreatedAt: w.CreatedAt.Time,
	}
}

// PasswordResetModel is a model which encapsulates the operations of the object
type PasswordResetModel struct {
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

func (m *PasswordResetModel) BeforeCreate(f func(kv query.KV) error) {
	m.beforeCreate = f
}

func (m *PasswordResetModel) AfterCreate(f func(id int64) error) {
	m.afterCreate = f
}

func (m *PasswordResetModel) BeforeUpdate(f func(kv query.KV) error) {
	m.beforeUpdate = f
}

func (m *PasswordResetModel) BeforeDelete(f func() error) {
	m.beforeDelete = f
}

func (m *PasswordResetModel) AfterDelete(f func() error) {
	m.afterDelete = f
}

var passwordresetTableName = "wz_passwordreset"

func SetPasswordResetTable(tableName string) {
	passwordresetTableName = tableName
}

// NewPasswordResetModel create a PasswordResetModel
func NewPasswordResetModel(db query.Database) *PasswordResetModel {
	return &PasswordResetModel{
		db:                  query.NewDatabaseWrap(db),
		tableName:           passwordresetTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *PasswordResetModel) GetDB() query.Database {
	return m.db.GetDB()
}

func (m *PasswordResetModel) clone() *PasswordResetModel {
	return &PasswordResetModel{
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
func (m *PasswordResetModel) WithoutGlobalScopes(names ...string) *PasswordResetModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *PasswordResetModel) WithLocalScopes(names ...string) *PasswordResetModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Query add query builder to model
func (m *PasswordResetModel) Query(builder query.SQLBuilder) *PasswordResetModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *PasswordResetModel) Find(id int64) (PasswordReset, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *PasswordResetModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *PasswordResetModel) Count(builders ...query.SQLBuilder) (int64, error) {
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

func (m *PasswordResetModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]PasswordReset, query.PaginateMeta, error) {
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
func (m *PasswordResetModel) Get(builders ...query.SQLBuilder) ([]PasswordReset, error) {
	sqlStr, params := m.query.Merge(builders...).
		Table(m.tableName).
		Select(
			"email",
			"token",
			"id",
			"created_at",
		).AppendCondition(m.applyScope()).
		ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	passwordresets := make([]PasswordReset, 0)
	for rows.Next() {
		var passwordresetVar passwordResetWrap
		if err := rows.Scan(
			&passwordresetVar.Email,
			&passwordresetVar.Token,
			&passwordresetVar.Id,
			&passwordresetVar.CreatedAt); err != nil {
			return nil, err
		}

		passwordresetReal := passwordresetVar.ToPasswordReset()
		passwordresetReal.SetModel(m)
		passwordresets = append(passwordresets, passwordresetReal)
	}

	return passwordresets, nil
}

// First return first result for given query
func (m *PasswordResetModel) First(builders ...query.SQLBuilder) (PasswordReset, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return PasswordReset{}, err
	}

	if len(res) == 0 {
		return PasswordReset{}, query.ErrNoResult
	}

	return res[0], nil
}

// Create save a new PasswordReset to database
func (m *PasswordResetModel) Create(kv query.KV) (int64, error) {
	kv["created_at"] = time.Now()

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

// SaveAll save all PasswordResets to database
func (m *PasswordResetModel) SaveAll(passwordresets []PasswordReset) ([]int64, error) {
	ids := make([]int64, 0)
	for _, passwordreset := range passwordresets {
		id, err := m.Save(passwordreset)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a PasswordReset to database
func (m *PasswordResetModel) Save(passwordreset PasswordReset) (int64, error) {
	return m.Create(query.KV{
		"email": passwordreset.Email,
		"token": passwordreset.Token,
	})
}

// SaveOrUpdate save a new PasswordReset or update it when it has a id > 0
func (m *PasswordResetModel) SaveOrUpdate(passwordreset PasswordReset) (id int64, updated bool, err error) {
	if passwordreset.Id > 0 {
		_, _err := m.UpdateById(passwordreset.Id, passwordreset)
		return passwordreset.Id, true, _err
	}

	_id, _err := m.Save(passwordreset)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *PasswordResetModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

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
func (m *PasswordResetModel) Update(passwordreset PasswordReset) (int64, error) {
	return m.UpdateFields(passwordreset.StaledKV())
}

// UpdateById update a model by id
func (m *PasswordResetModel) UpdateById(id int64, passwordreset PasswordReset) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Update(passwordreset)
}

// Delete remove a model
func (m *PasswordResetModel) Delete(builders ...query.SQLBuilder) (int64, error) {
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
func (m *PasswordResetModel) DeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Delete()
}

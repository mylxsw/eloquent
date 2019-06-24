package models

// !!! DO NOT EDIT THIS FILE

import (
	"context"
	"database/sql"
	"github.com/mylxsw/eloquent/query"
	"gopkg.in/guregu/null.v3"
	"time"
)

func init() {

}

// Page is a Page object
type Page struct {
	original  *pageOriginal
	pageModel *PageModel

	Id              int64
	Pid             int64
	Title           string
	Description     string
	Content         string
	ProjectId       int64
	UserId          int64
	Type            int
	Status          int
	LastModifiedUid int64
	HistoryId       int64
	SortLevel       int
	SyncUrl         string
	LastSyncAt      time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// SetModel set model for Page
func (pageSelf *Page) SetModel(pageModel *PageModel) {
	pageSelf.pageModel = pageModel
}

// pageOriginal is an object which stores original Page from database
type pageOriginal struct {
	Id              int64
	Pid             int64
	Title           string
	Description     string
	Content         string
	ProjectId       int64
	UserId          int64
	Type            int
	Status          int
	LastModifiedUid int64
	HistoryId       int64
	SortLevel       int
	SyncUrl         string
	LastSyncAt      time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Staled identify whether the object has been modified
func (pageSelf *Page) Staled() bool {
	if pageSelf.original == nil {
		pageSelf.original = &pageOriginal{}
	}

	if pageSelf.Id != pageSelf.original.Id {
		return true
	}

	if pageSelf.Pid != pageSelf.original.Pid {
		return true
	}
	if pageSelf.Title != pageSelf.original.Title {
		return true
	}
	if pageSelf.Description != pageSelf.original.Description {
		return true
	}
	if pageSelf.Content != pageSelf.original.Content {
		return true
	}
	if pageSelf.ProjectId != pageSelf.original.ProjectId {
		return true
	}
	if pageSelf.UserId != pageSelf.original.UserId {
		return true
	}
	if pageSelf.Type != pageSelf.original.Type {
		return true
	}
	if pageSelf.Status != pageSelf.original.Status {
		return true
	}
	if pageSelf.LastModifiedUid != pageSelf.original.LastModifiedUid {
		return true
	}
	if pageSelf.HistoryId != pageSelf.original.HistoryId {
		return true
	}
	if pageSelf.SortLevel != pageSelf.original.SortLevel {
		return true
	}
	if pageSelf.SyncUrl != pageSelf.original.SyncUrl {
		return true
	}
	if pageSelf.LastSyncAt != pageSelf.original.LastSyncAt {
		return true
	}

	if pageSelf.CreatedAt != pageSelf.original.CreatedAt {
		return true
	}
	if pageSelf.UpdatedAt != pageSelf.original.UpdatedAt {
		return true
	}

	return false
}

// StaledKV return all fields has been modified
func (pageSelf *Page) StaledKV() query.KV {
	kv := make(query.KV, 0)

	if pageSelf.original == nil {
		pageSelf.original = &pageOriginal{}
	}

	if pageSelf.Id != pageSelf.original.Id {
		kv["id"] = pageSelf.Id
	}

	if pageSelf.Pid != pageSelf.original.Pid {
		kv["pid"] = pageSelf.Pid
	}
	if pageSelf.Title != pageSelf.original.Title {
		kv["title"] = pageSelf.Title
	}
	if pageSelf.Description != pageSelf.original.Description {
		kv["description"] = pageSelf.Description
	}
	if pageSelf.Content != pageSelf.original.Content {
		kv["content"] = pageSelf.Content
	}
	if pageSelf.ProjectId != pageSelf.original.ProjectId {
		kv["project_id"] = pageSelf.ProjectId
	}
	if pageSelf.UserId != pageSelf.original.UserId {
		kv["user_id"] = pageSelf.UserId
	}
	if pageSelf.Type != pageSelf.original.Type {
		kv["type"] = pageSelf.Type
	}
	if pageSelf.Status != pageSelf.original.Status {
		kv["status"] = pageSelf.Status
	}
	if pageSelf.LastModifiedUid != pageSelf.original.LastModifiedUid {
		kv["last_modified_uid"] = pageSelf.LastModifiedUid
	}
	if pageSelf.HistoryId != pageSelf.original.HistoryId {
		kv["history_id"] = pageSelf.HistoryId
	}
	if pageSelf.SortLevel != pageSelf.original.SortLevel {
		kv["sort_level"] = pageSelf.SortLevel
	}
	if pageSelf.SyncUrl != pageSelf.original.SyncUrl {
		kv["sync_url"] = pageSelf.SyncUrl
	}
	if pageSelf.LastSyncAt != pageSelf.original.LastSyncAt {
		kv["last_sync_at"] = pageSelf.LastSyncAt
	}

	if pageSelf.CreatedAt != pageSelf.original.CreatedAt {
		kv["created_at"] = pageSelf.CreatedAt
	}
	if pageSelf.UpdatedAt != pageSelf.original.UpdatedAt {
		kv["updated_at"] = pageSelf.UpdatedAt
	}

	return kv
}

func (pageSelf *Page) Project() *ProjectModel {

	q := query.Builder().Where("id", pageSelf.ProjectId)

	return NewProjectModel(pageSelf.pageModel.GetDB()).Query(q)
}

// Save create a new model or update it
func (pageSelf *Page) Save() error {
	if pageSelf.pageModel == nil {
		return query.ErrModelNotSet
	}

	id, _, err := pageSelf.pageModel.SaveOrUpdate(*pageSelf)
	if err != nil {
		return err
	}

	pageSelf.Id = id
	return nil
}

// Delete remove a page
func (pageSelf *Page) Delete() error {
	if pageSelf.pageModel == nil {
		return query.ErrModelNotSet
	}

	_, err := pageSelf.pageModel.DeleteById(pageSelf.Id)
	if err != nil {
		return err
	}

	return nil
}

type pageScope struct {
	name  string
	apply func(builder query.Condition)
}

var pageGlobalScopes = make([]pageScope, 0)
var pageLocalScopes = make([]pageScope, 0)

// AddPageGlobalScope assign a global scope to a model
func AddPageGlobalScope(name string, apply func(builder query.Condition)) {
	pageGlobalScopes = append(pageGlobalScopes, pageScope{name: name, apply: apply})
}

// AddPageLocalScope assign a local scope to a model
func AddPageLocalScope(name string, apply func(builder query.Condition)) {
	pageLocalScopes = append(pageLocalScopes, pageScope{name: name, apply: apply})
}

func (m *PageModel) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range pageGlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range pageLocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *PageModel) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *PageModel) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}

	return true
}

type pageWrap struct {
	Id              null.Int
	Pid             null.Int
	Title           null.String
	Description     null.String
	Content         null.String
	ProjectId       null.Int
	UserId          null.Int
	Type            null.Int
	Status          null.Int
	LastModifiedUid null.Int
	HistoryId       null.Int
	SortLevel       null.Int
	SyncUrl         null.String
	LastSyncAt      null.Time

	CreatedAt null.Time
	UpdatedAt null.Time
}

func (w pageWrap) ToPage() Page {
	return Page{
		original: &pageOriginal{
			Id:              w.Id.Int64,
			Pid:             w.Pid.Int64,
			Title:           w.Title.String,
			Description:     w.Description.String,
			Content:         w.Content.String,
			ProjectId:       w.ProjectId.Int64,
			UserId:          w.UserId.Int64,
			Type:            int(w.Type.Int64),
			Status:          int(w.Status.Int64),
			LastModifiedUid: w.LastModifiedUid.Int64,
			HistoryId:       w.HistoryId.Int64,
			SortLevel:       int(w.SortLevel.Int64),
			SyncUrl:         w.SyncUrl.String,
			LastSyncAt:      w.LastSyncAt.Time,

			CreatedAt: w.CreatedAt.Time,
			UpdatedAt: w.UpdatedAt.Time,
		},
		Id:              w.Id.Int64,
		Pid:             w.Pid.Int64,
		Title:           w.Title.String,
		Description:     w.Description.String,
		Content:         w.Content.String,
		ProjectId:       w.ProjectId.Int64,
		UserId:          w.UserId.Int64,
		Type:            int(w.Type.Int64),
		Status:          int(w.Status.Int64),
		LastModifiedUid: w.LastModifiedUid.Int64,
		HistoryId:       w.HistoryId.Int64,
		SortLevel:       int(w.SortLevel.Int64),
		SyncUrl:         w.SyncUrl.String,
		LastSyncAt:      w.LastSyncAt.Time,

		CreatedAt: w.CreatedAt.Time,
		UpdatedAt: w.UpdatedAt.Time,
	}
}

// PageModel is a model which encapsulates the operations of the object
type PageModel struct {
	db        query.Database
	tableName string

	excludeGlobalScopes []string
	includeLocalScopes  []string

	query query.SQLBuilder
}

var pageTableName = "wz_pages"

func SetPageTable(tableName string) {
	pageTableName = tableName
}

// NewPageModel create a PageModel
func NewPageModel(db query.Database) *PageModel {
	return &PageModel{
		db:                  db,
		tableName:           pageTableName,
		excludeGlobalScopes: make([]string, 0),
		includeLocalScopes:  make([]string, 0),
		query:               query.Builder(),
	}
}

// GetDB return database instance
func (m *PageModel) GetDB() query.Database {
	return m.db
}

func (m *PageModel) clone() *PageModel {
	return &PageModel{
		db:                  m.db,
		tableName:           m.tableName,
		excludeGlobalScopes: append([]string{}, m.excludeGlobalScopes...),
		includeLocalScopes:  append([]string{}, m.includeLocalScopes...),
		query:               m.query,
	}
}

// WithoutGlobalScopes remove a global scope for given query
func (m *PageModel) WithoutGlobalScopes(names ...string) *PageModel {
	mc := m.clone()
	mc.excludeGlobalScopes = append(mc.excludeGlobalScopes, names...)

	return mc
}

// WithLocalScopes add a local scope for given query
func (m *PageModel) WithLocalScopes(names ...string) *PageModel {
	mc := m.clone()
	mc.includeLocalScopes = append(mc.includeLocalScopes, names...)

	return mc
}

// Query add query builder to model
func (m *PageModel) Query(builder query.SQLBuilder) *PageModel {
	mm := m.clone()
	mm.query = mm.query.Merge(builder)

	return mm
}

// Find retrieve a model by its primary key
func (m *PageModel) Find(id int64) (Page, error) {
	return m.First(m.query.Where("id", "=", id))
}

// Exists return whether the records exists for a given query
func (m *PageModel) Exists(builders ...query.SQLBuilder) (bool, error) {
	count, err := m.Count(builders...)
	return count > 0, err
}

// Count return model count for a given query
func (m *PageModel) Count(builders ...query.SQLBuilder) (int64, error) {
	sqlStr, params := m.query.Merge(builders...).Table(m.tableName).ResolveCount()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	rows.Next()
	var res int64
	if err := rows.Scan(&res); err != nil {
		return 0, err
	}

	return res, nil
}

func (m *PageModel) Paginate(page int64, perPage int64, builders ...query.SQLBuilder) ([]Page, query.PaginateMeta, error) {
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
func (m *PageModel) Get(builders ...query.SQLBuilder) ([]Page, error) {
	sqlStr, params := m.query.Merge(builders...).
		Table(m.tableName).
		Select("id", "created_at", "updated_at", "pid", "title", "description", "content", "project_id", "user_id", "type", "status", "last_modified_uid", "history_id", "sort_level", "sync_url", "last_sync_at").
		AppendCondition(m.applyScope()).
		ResolveQuery()

	rows, err := m.db.QueryContext(context.Background(), sqlStr, params...)
	if err != nil {
		return nil, err
	}

	pages := make([]Page, 0)
	for rows.Next() {
		var pageVar pageWrap
		if err := rows.Scan(
			&pageVar.Id,
			&pageVar.CreatedAt,
			&pageVar.UpdatedAt,
			&pageVar.Pid,
			&pageVar.Title,
			&pageVar.Description,
			&pageVar.Content,
			&pageVar.ProjectId,
			&pageVar.UserId,
			&pageVar.Type,
			&pageVar.Status,
			&pageVar.LastModifiedUid,
			&pageVar.HistoryId,
			&pageVar.SortLevel,
			&pageVar.SyncUrl,
			&pageVar.LastSyncAt); err != nil {
			return nil, err
		}

		pageReal := pageVar.ToPage()
		pageReal.SetModel(m)
		pages = append(pages, pageReal)
	}

	return pages, nil
}

// First return first result for given query
func (m *PageModel) First(builders ...query.SQLBuilder) (Page, error) {
	res, err := m.Get(append(builders, query.Builder().Limit(1))...)
	if err != nil {
		return Page{}, err
	}

	if len(res) == 0 {
		return Page{}, sql.ErrNoRows
	}

	return res[0], nil
}

// Create save a new page to database
func (m *PageModel) Create(kv query.KV) (int64, error) {
	kv["created_at"] = time.Now()
	kv["updated_at"] = time.Now()

	sqlStr, params := m.query.Table(m.tableName).ResolveInsert(kv)

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// SaveAll save all pages to database
func (m *PageModel) SaveAll(pages []Page) ([]int64, error) {
	ids := make([]int64, 0)
	for _, page := range pages {
		id, err := m.Save(page)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// Save save a page to database
func (m *PageModel) Save(page Page) (int64, error) {
	return m.Create(query.KV{
		"pid":               page.Pid,
		"title":             page.Title,
		"description":       page.Description,
		"content":           page.Content,
		"project_id":        page.ProjectId,
		"user_id":           page.UserId,
		"type":              page.Type,
		"status":            page.Status,
		"last_modified_uid": page.LastModifiedUid,
		"history_id":        page.HistoryId,
		"sort_level":        page.SortLevel,
		"sync_url":          page.SyncUrl,
		"last_sync_at":      page.LastSyncAt,
	})
}

// SaveOrUpdate save a new page or update it when it has a id > 0
func (m *PageModel) SaveOrUpdate(page Page) (id int64, updated bool, err error) {
	if page.Id > 0 {
		_, _err := m.UpdateById(page.Id, page)
		return page.Id, true, _err
	}

	_id, _err := m.Save(page)
	return _id, false, _err
}

// UpdateFields update kv for a given query
func (m *PageModel) UpdateFields(kv query.KV, builders ...query.SQLBuilder) (int64, error) {
	if len(kv) == 0 {
		return 0, nil
	}

	kv["updated_at"] = time.Now()

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
func (m *PageModel) Update(page Page) (int64, error) {
	return m.UpdateFields(page.StaledKV())
}

// UpdateById update a model by id
func (m *PageModel) UpdateById(id int64, page Page) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Update(page)
}

// Delete remove a model
func (m *PageModel) Delete(builders ...query.SQLBuilder) (int64, error) {

	sqlStr, params := m.query.Merge(builders...).AppendCondition(m.applyScope()).Table(m.tableName).ResolveDelete()

	res, err := m.db.ExecContext(context.Background(), sqlStr, params...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

// DeleteById remove a model by id
func (m *PageModel) DeleteById(id int64) (int64, error) {
	return m.Query(query.Builder().Where("id", "=", id)).Delete()
}

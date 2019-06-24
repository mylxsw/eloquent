package query

import (
	"fmt"
	"strings"
)

type SQLBuilder struct {
	tableName  string
	conditions Condition
	fields     []expr
	limit      int64
	offset     int64
	orders     orderBys
	groups     []string
	having     ConditionGroup
	joins      []sqlJoin
	unions     []sqlUnion
}

type sqlUnion struct {
	Type     string
	SubQuery SubQuery
}

type expr struct {
	Value    string
	Type     exprType
	Bindings []interface{}
}

type exprType int

const (
	exprTypeString exprType = iota
	exprTypeRaw
)

func Raw(rawStr string, bindings ...interface{}) expr {
	return expr{
		Type:     exprTypeRaw,
		Value:    rawStr,
		Bindings: bindings,
	}
}

func Builder() SQLBuilder {
	return SQLBuilder{
		conditions: ConditionBuilder(),
		fields:     make([]expr, 0),
		orders:     make([]sqlOrderBy, 0),
		groups:     make([]string, 0),
		joins:      make([]sqlJoin, 0),
		unions:     make([]sqlUnion, 0),
		limit:      -1,
		offset:     -1,
	}
}


func (builder SQLBuilder) Merge(b2s ...SQLBuilder) SQLBuilder {
	b := builder.Clone()
	for _, b2 := range b2s {
		b.conditions.Append(b2.conditions)
		b.fields = append(b.fields, b2.fields...)
		b.orders = append(b.orders, b2.orders...)
		b.groups = append(b.groups, b2.groups...)
		b.joins = append(b.joins, b2.joins...)
		b.unions = append(b.unions, b2.unions...)

		if b2.limit != -1 {
			b.limit = b2.limit
		}

		if b2.offset != -1 {
			b.offset = b2.offset
		}
	}

	return b
}

func (builder SQLBuilder) Clone() SQLBuilder {
	b := SQLBuilder{
		conditions: builder.conditions.Clone(),
		tableName:  builder.tableName,
		limit:      builder.limit,
		offset:     builder.offset,
		having:     builder.having,
		fields:     append([]expr{}, builder.fields...),
		orders:     append([]sqlOrderBy{}, builder.orders...),
		groups:     append([]string{}, builder.groups...),
		joins:      append([]sqlJoin{}, builder.joins...),
		unions:     append([]sqlUnion{}, builder.unions...),
	}

	return b
}

type KV map[string]interface{}

func (builder SQLBuilder) ResolveDelete() (string, []interface{}) {
	sqlStr := fmt.Sprintf("DELETE FROM %s", builder.tableName)
	values := make([]interface{}, 0)

	tableAlias := resolveTableAlias(builder.tableName)
	if !builder.conditions.Empty() {
		conditions, p := builder.conditions.Resolve(tableAlias)
		sqlStr += fmt.Sprintf(" WHERE %s", strings.TrimLeft(conditions, " "))
		values = append(values, p...)
	}

	if len(builder.orders) > 0 {
		sqlStr += fmt.Sprintf(" ORDER BY %s", builder.orders.String(tableAlias))
	}

	if builder.limit >= 0 {
		sqlStr += fmt.Sprintf(" LIMIT %d", builder.limit)
	}

	return sqlStr, values
}

func (builder SQLBuilder) ResolveUpdate(kvPairs KV) (string, []interface{}) {

	sqlStr := fmt.Sprintf("UPDATE %s", builder.tableName)

	fields, values := builder.resolveKvPairsForUpdate(kvPairs)
	sqlStr += fmt.Sprintf(" SET %s", strings.Join(fields, ", "))

	if !builder.conditions.Empty() {
		conditions, p := builder.conditions.Resolve(resolveTableAlias(builder.tableName))
		sqlStr += fmt.Sprintf(" WHERE %s", strings.TrimLeft(conditions, " "))
		values = append(values, p...)
	}

	return sqlStr, values
}

func (builder SQLBuilder) resolveKvPairsForUpdate(kvPairs KV) ([]string, []interface{}) {
	fields := make([]string, len(kvPairs))
	values := make([]interface{}, 0)
	var i = 0
	for k, v := range kvPairs {
		vv, ok := v.(expr)
		if ok {
			switch vv.Type {
			case exprTypeString:
				fields[i] = "`" + k + "` = ?"
				values = append(values, vv.Value)
			case exprTypeRaw:
				fields[i] = "`" + k + "` = " + vv.Value
				values = append(values, vv.Bindings...)
			}
		} else {
			fields[i] = "`" + k + "` = ?"
			values = append(values, v)
		}

		i++
	}
	return fields, values
}

func (builder SQLBuilder) resolveKvPairsForInsert(kvPairs KV) ([]string, []interface{}) {
	fields := make([]string, len(kvPairs))
	values := make([]interface{}, 0)
	var i = 0
	for k, v := range kvPairs {
		fields[i] = "`" + k + "`"

		vv, ok := v.(expr)
		if ok {
			values = append(values, vv.Value)
			if vv.Bindings != nil && len(vv.Bindings) > 0 {
				values = append(values, vv.Bindings...)
			}
		} else {
			values = append(values, v)
		}

		i++
	}
	return fields, values
}

func (builder SQLBuilder) ResolveInsert(kvPairs KV) (string, []interface{}) {
	sqlStr := fmt.Sprintf("INSERT INTO %s", builder.tableName)

	fields, values := builder.resolveKvPairsForInsert(kvPairs)
	sqlStr += fmt.Sprintf(" (%s) VALUES (%s)", strings.Join(fields, ", "), strings.Trim(strings.Repeat("?,", len(values)), ","))

	return sqlStr, values
}

func (builder SQLBuilder) ResolveCount() (string, []interface{}) {
	b := builder.Clone()
	b.fields = make([]expr, 1)
	b.fields[0] = Raw("COUNT(1) as count")

	return b.ResolveQuery()
}

func (builder SQLBuilder) ResolveMax(field string) (string, []interface{}) {
	b := builder.Clone()
	b.fields = make([]expr, 1)
	b.fields[0] = Raw(fmt.Sprintf("MAX(%s) as max", field))

	return b.ResolveQuery()
}

func (builder SQLBuilder) ResolveMin(field string) (string, []interface{}) {
	b := builder.Clone()
	b.fields = make([]expr, 1)
	b.fields[0] = Raw(fmt.Sprintf("MIN(%s) as min", field))

	return b.ResolveQuery()
}

func (builder SQLBuilder) ResolveAvg(field string) (string, []interface{}) {
	b := builder.Clone()
	b.fields = make([]expr, 1)
	b.fields[0] = Raw(fmt.Sprintf("AVG(%s) as avg", field))

	return b.ResolveQuery()
}

func (builder SQLBuilder) ResolveSum(field string) (string, []interface{}) {
	b := builder.Clone()
	b.fields = make([]expr, 1)
	b.fields[0] = Raw(fmt.Sprintf("SUM(%s) as sum", field))

	return b.ResolveQuery()
}

func (builder SQLBuilder) ResolveQuery() (string, []interface{}) {
	tableAlias := resolveTableAlias(builder.tableName)

	params := make([]interface{}, 0)

	fields, p := builder.getFields(tableAlias)
	params = append(params, p...)
	sqlStr := fmt.Sprintf("SELECT %s FROM %s", fields, builder.tableName)

	if len(builder.joins) > 0 {
		for _, j := range builder.joins {
			s, p := j.String(tableAlias)
			sqlStr += " " + s
			params = append(params, p...)
		}
	}

	if !builder.conditions.Empty() {
		conditions, p := builder.conditions.Resolve(tableAlias)
		sqlStr += fmt.Sprintf(" WHERE %s", strings.TrimLeft(conditions, " "))
		params = append(params, p...)
	}

	if len(builder.groups) > 0 {
		var groupBys = make([]string, len(builder.groups))
		for i, g := range builder.groups {
			groupBys[i] = replaceTableField(tableAlias, g)
		}
		sqlStr += fmt.Sprintf(" GROUP BY %s", strings.Join(groupBys, ", "))
	}

	if builder.having != nil {
		newBuilder := ConditionBuilder()
		builder.having(newBuilder)

		havingCond, havingParams := newBuilder.Resolve(tableAlias)
		if havingCond != "" {
			sqlStr += fmt.Sprintf(" HAVING %s", havingCond)
			params = append(params, havingParams...)
		}
	}

	if len(builder.orders) > 0 {
		sqlStr += fmt.Sprintf(" ORDER BY %s", builder.orders.String(tableAlias))
	}

	if builder.limit >= 0 {
		sqlStr += fmt.Sprintf(" LIMIT %d", builder.limit)
	}
	if builder.offset >= 0 {
		sqlStr += fmt.Sprintf(" OFFSET %d", builder.offset)
	}

	if len(builder.unions) > 0 {
		sqlStr = "(" + sqlStr + ")"
		for _, u := range builder.unions {
			s, p := u.SubQuery.ResolveQuery()
			sqlStr += fmt.Sprintf(" UNION %s (%s)", u.Type, s)
			params = append(params, p...)
		}
	}

	return sqlStr, params
}

func (builder SQLBuilder) Condition(where Condition) SQLBuilder {
	b := builder.Clone()
	b.conditions = where

	return b
}

func (builder SQLBuilder) Table(name string) SQLBuilder {
	b := builder.Clone()
	b.tableName = name

	return b
}

func (builder SQLBuilder) LeftJoin(tableName string, on ConditionGroup) SQLBuilder {
	return builder.join(sqlJoin{
		joinType:  "LEFT JOIN",
		tableName: tableName,
		on:        on,
	})
}

func (builder SQLBuilder) RightJoin(tableName string, on ConditionGroup) SQLBuilder {
	return builder.join(sqlJoin{
		joinType:  "RIGHT JOIN",
		tableName: tableName,
		on:        on,
	})
}

func (builder SQLBuilder) InnerJoin(tableName string, on ConditionGroup) SQLBuilder {
	return builder.join(sqlJoin{
		joinType:  "INNER JOIN",
		tableName: tableName,
		on:        on,
	})
}

func (builder SQLBuilder) CrossJoin(tableName string, on ConditionGroup) SQLBuilder {
	return builder.join(sqlJoin{
		joinType:  "CROSS JOIN",
		tableName: tableName,
		on:        on,
	})
}

func (builder SQLBuilder) join(join sqlJoin) SQLBuilder {
	b := builder.Clone()
	b.joins = append(builder.joins, join)
	return b
}

func (builder SQLBuilder) Union(b2 SubQuery, distinct bool) SQLBuilder {
	union := sqlUnion{
		SubQuery: b2,
	}

	if distinct {
		union.Type = "DISTINCT"
	} else {
		union.Type = "ALL"
	}

	b := builder.Clone()
	b.unions = append(builder.unions, union)

	return b
}

func (builder SQLBuilder) Select(fields ...interface{}) SQLBuilder {
	b := builder.Clone()
	for _, f := range fields {
		f1, ok := f.(expr)
		if ok {
			b.fields = append(b.fields, f1)
		} else {
			b.fields = append(b.fields, expr{
				Type:  exprTypeString,
				Value: f.(string),
			})
		}
	}
	return b
}

func (builder SQLBuilder) Limit(limit int64) SQLBuilder {
	b := builder.Clone()
	b.limit = limit
	return b
}

func (builder SQLBuilder) Offset(offset int64) SQLBuilder {
	b := builder.Clone()
	b.offset = offset
	return b
}

func (builder SQLBuilder) OrderBy(field string, direction string) SQLBuilder {
	b := builder.Clone()
	b.orders = append(builder.orders, sqlOrderBy{raw: false, Field: field, Direction: direction})
	return b
}

func (builder SQLBuilder) OrderByRaw(raw string) SQLBuilder {
	b := builder.Clone()
	b.orders = append(builder.orders, sqlOrderBy{raw: true, Raw: raw})
	return b
}

func (builder SQLBuilder) GroupBy(fields ...string) SQLBuilder {
	b := builder.Clone()
	b.groups = append(builder.groups, fields...)
	return b
}

func (builder SQLBuilder) Having(closure ConditionGroup) SQLBuilder {
	b := builder.Clone()
	b.having = closure
	return b
}

func (builder SQLBuilder) getFields(tableAlias string) (string, []interface{}) {
	var params = make([]interface{}, 0)
	if len(builder.fields) == 0 {
		return "*", params
	}

	fields := make([]string, len(builder.fields))
	for i, f := range builder.fields {
		switch f.Type {
		case exprTypeString:
			fields[i] = replaceTableField(tableAlias, f.Value)
		case exprTypeRaw:
			fields[i] = f.Value
			params = append(params, f.Bindings...)
		}
	}

	return strings.Join(fields, ", "), params
}

func (builder SQLBuilder) WhereColumn(field, operator string, value string) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereColumn(field, operator, value)

	return b
}

func (builder SQLBuilder) OrWhereColumn(field, operator string, value string) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereColumn(field, operator, value)

	return b
}

func (builder SQLBuilder) OrWhereNotExist(subQuery SubQuery) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereNotExist(subQuery)

	return b
}

func (builder SQLBuilder) OrWhereExist(subQuery SubQuery) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereExist(subQuery)

	return b
}

func (builder SQLBuilder) WhereNotExist(subQuery SubQuery) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereNotExist(subQuery)

	return b
}

func (builder SQLBuilder) WhereExist(subQuery SubQuery) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereExist(subQuery)

	return b
}

func (builder SQLBuilder) OrWhereNotNull(field string) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereNotNull(field)

	return b
}

func (builder SQLBuilder) OrWhereNull(field string) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereNull(field)

	return b
}

func (builder SQLBuilder) WhereNotNull(field string) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereNotNull(field)

	return b
}

func (builder SQLBuilder) WhereNull(field string) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereNull(field)

	return b
}

func (builder SQLBuilder) OrWhereRaw(raw string, items ...interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereRaw(raw, items...)

	return b
}

func (builder SQLBuilder) WhereRaw(raw string, items ...interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereRaw(raw, items...)

	return b
}

func (builder SQLBuilder) OrWhereNotIn(field string, items ...interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereNotIn(field, items...)

	return b
}

func (builder SQLBuilder) OrWhereIn(field string, items ...interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereIn(field, items...)

	return b
}

func (builder SQLBuilder) WhereNotIn(field string, items ...interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereNotIn(field, items...)

	return b
}

func (builder SQLBuilder) WhereIn(field string, items ...interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereIn(field, items...)

	return b
}

func (builder SQLBuilder) WhereGroup(wc ConditionGroup) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereGroup(wc)

	return b
}

func (builder SQLBuilder) OrWhereGroup(wc ConditionGroup) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereGroup(wc)

	return b
}

func (builder SQLBuilder) Where(field string, value ...interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.Where(field, value...)

	return b
}

func (builder SQLBuilder) OrWhere(field string, value ...interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhere(field, value...)

	return b
}

func (builder SQLBuilder) When(when When, cg ConditionGroup) SQLBuilder {
	b := builder.Clone()
	b.conditions.When(when, cg)

	return b
}

func (builder SQLBuilder) OrWhen(when When, cg ConditionGroup) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhen(when, cg)

	return b
}

func (builder SQLBuilder) WhereBetween(field string, min interface{}, max interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereBetween(field, min, max)

	return b
}

func (builder SQLBuilder) OrWhereBetween(field string, min interface{}, max interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereBetween(field, min, max)

	return b
}

func (builder SQLBuilder) WhereNotBetween(field string, min interface{}, max interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.WhereNotBetween(field, min, max)

	return b
}

func (builder SQLBuilder) OrWhereNotBetween(field string, min interface{}, max interface{}) SQLBuilder {
	b := builder.Clone()
	b.conditions.OrWhereNotBetween(field, min, max)

	return b
}

func (builder SQLBuilder) AppendCondition(cond Condition) SQLBuilder {
	b := builder.Clone()
	b.conditions.Append(cond)

	return b
}

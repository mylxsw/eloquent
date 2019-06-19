package query

import (
	"fmt"
)

type sqlJoin struct {
	joinType  string
	tableName string
	on        ConditionGroup
}

func (t sqlJoin) String(tableAlias string) (string, []interface{}) {
	newBuilder := ConditionBuilder()
	t.on(newBuilder)

	newCond, newParams := newBuilder.Resolve(tableAlias)
	sql := fmt.Sprintf("%s %s ON %s", t.joinType, t.tableName, newCond)
	return sql, newParams
}

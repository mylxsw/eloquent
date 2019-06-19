package query

import (
	"fmt"
	"strings"
)

type orderBys []sqlOrderBy

func (ob orderBys) String(tableAlias string) string {
	var obs = make([]string, len(ob))
	for i, o := range ob {
		obs[i] = fmt.Sprintf("%s %s", replaceTableField(tableAlias, o.Field), o.Direction)
	}

	return strings.Join(obs, ",")
}

type sqlOrderBy struct {
	Field     string
	Direction string
}

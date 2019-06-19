package query

import (
	"strings"
)

func isSubQuery(values []interface{}) bool {
	if len(values) != 1 {
		return false
	}

	if _, ok := values[0].(SubQuery); ok {
		return true
	}

	return false
}

func replaceTableField(tableAlias string, name string) string {
	segs1 := strings.Split(name, " ")
	org := segs1[0]
	segs1len := len(segs1)
	if segs1len == 3 {
		return resolveOrgTableField(tableAlias, org) + " AS " + segs1[2]
	} else if segs1len == 2 {
		return resolveOrgTableField(tableAlias, org) + " AS " + segs1[1]
	}

	// a.b      => a.`b`
	// b        => alias.`b`
	// b as c   => alias.`b` as c
	// a.b as c => a.`b` as c

	return resolveOrgTableField(tableAlias, org)
}

func resolveOrgTableField(tableAlias string, org string) string {
	segs := strings.Split(org, ".")
	if len(segs) > 1 {
		if segs[1] != "*" {
			segs[1] = "`" + segs[1] + "`"
		}
	} else if segs[0] != "*" {
		segs[0] = "`" + segs[0] + "`"
	}

	if tableAlias != "" && len(segs) == 1 {
		return tableAlias + "." + strings.Join(segs, ".")
	}

	return strings.Join(segs, ".")
}

func resolveTableAlias(name string) string {
	segs := strings.Split(name, " ")
	if len(segs) == 3 && strings.ToUpper(segs[1]) == "AS" {
		return segs[2]
	} else if len(segs) == 2 {
		return segs[1]
	}

	return segs[0]
}

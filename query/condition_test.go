package query

import (
	"fmt"
	"testing"
)

func TestConditionBuilder_Clone(t *testing.T) {
	builder := ConditionBuilder()
	builder.Where("field1", "=", 123)
	builder.WhereIn("field2", 134, 521, 341)

	sql, params := builder.Resolve("")
	res1 := fmt.Sprint(sql, params)

	newBuilder := builder.Clone()
	newBuilder.Where("field3", ">", 199)

	sql, params = newBuilder.Resolve("")
	res2 := fmt.Sprint(sql, params)

	sql, params = builder.Resolve("")
	res3 := fmt.Sprint(sql, params)

	if res1 != res3 || res1 == res2 {
		t.Errorf("test failed")
	}
}

package template_test

import (
	"encoding/json"
	"testing"

	"github.com/mylxsw/go-utils/assert"
	"gopkg.in/guregu/null.v3"
)

type testStruct struct {
	S1 null.String `json:"s1"`
	S2 null.Int
}

func TestEntity(t *testing.T) {
	s1 := null.StringFrom("Hello, world")
	s2 := null.StringFrom("Hello, world")

	assert.True(t, s1 == s2)

	s1.Valid = false
	assert.True(t, s1 != s2)

	ts := testStruct{S1: null.StringFrom("Hello, world"), S2: null.IntFrom(1230)}
	{
		data, _ := json.Marshal(ts)
		assert.Equal(t, `{"s1":"Hello, world","S2":1230}`, string(data))
	}
	{
		ts.S2 = null.NewInt(0, false)

		data, _ := json.Marshal(ts)
		assert.Equal(t, `{"s1":"Hello, world","S2":null}`, string(data))
	}
}

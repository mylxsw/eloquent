package generator_test

import (
	"fmt"
	"testing"

	"github.com/mylxsw/eloquent/generator"
	"gopkg.in/yaml.v2"
)

func TestDomain(t *testing.T) {
	domain := generator.Domain{
		Imports: []string{
			"github.com/mylxsw/eloquent",
		},
		PackageName: "models",
		Models: []generator.Model{
			{
				Name: "user",
				Relations: []generator.Relation{
					{
						Model:      "role",
						Rel:        "n-1",
						ForeignKey: "role_id",
						OwnerKey:   "id",
					},
				},
				Definition: generator.Definition{
					TableName:         "user",
					WithoutCreateTime: false,
					WithoutUpdateTime: false,
					SoftDelete:        false,
					Fields: []generator.DefinitionField{
						{
							Name: "id",
							Type: "int64",
							Tag:  `json:"id"`,
						},
						{
							Name: "name",
							Type: "string",
							Tag:  `json:"name"`,
						},
						{
							Name: "age",
							Type: "int64",
							Tag:  `json:"age"`,
						},
					},
				},
			},
		},
		Meta: generator.Meta{
			TablePrefix: "el_",
		},
	}

	marshal, err := yaml.Marshal(domain)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshal))
}

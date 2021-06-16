package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mylxsw/eloquent/generator"
	"github.com/mylxsw/eloquent/generator/template"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func main() {
	app := &cli.App{
		Name: "Eloquent 命令行工具",
		Commands: []cli.Command{
			{
				Name: "gen",
                Usage: "根据模型文件定义生成模型对象",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "source",
						Required: true,
						Usage:    "模型定义所在文件的 Glob 表达式，比如 ./models/*.yml",
					},
				},
				Action: func(c *cli.Context) error {
					source := c.String("source")
					matches, err := filepath.Glob(source)
					if err != nil {
						return err
					}

					for _, m := range matches {
						dest := replaceExt(m, ".orm.go")

						input, err := ioutil.ReadFile(m)
						assertError(err)

						var domain generator.Domain
						assertError(yaml.Unmarshal(input, &domain))

						res, err := generator.ParseTemplate(template.GetTemplate(), domain.Init())
						assertError(err)

						assertError(ioutil.WriteFile(dest, []byte(res), os.ModePerm))

						fmt.Println(dest)
					}

					return nil
				},
			},
			{
				Name:  "create-model",
				Usage: "创建表模型定义文件",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "table", Required: true, Usage: "表名"},
					cli.StringFlag{Name: "package", Required: true, Usage: "包名"},
					cli.StringFlag{Name: "output", Usage: "输出目录"},
					cli.StringFlag{Name: "table-prefix", Usage: "表前缀"},
					cli.BoolFlag{Name: "no-created_at", Usage: "不自动添加 created_at 字段"},
					cli.BoolFlag{Name: "no-updated_at", Usage: "不自动添加 updated_at 字段"},
					cli.BoolFlag{Name: "soft-delete", Usage: "启用软删除支持"},
					cli.StringSliceFlag{Name: "import", Usage: "引入包"},
				},
				Action: func(c *cli.Context) error {
					domain := generator.Domain{
						PackageName: c.String("package"),
						Models: []generator.Model{
							{
								Name: c.String("table"),
								Definition: generator.Definition{
									TableName:         c.String("table"),
									WithoutCreateTime: c.Bool("no-created_at"),
									WithoutUpdateTime: c.Bool("no-updated_at"),
									SoftDelete:        c.BoolT("soft-delete"),
									Fields: []generator.DefinitionField{
										{Name: "id", Type: "int64", Tag: `json:"id"`},
									},
								},
							},
						},
					}

					imports := c.StringSlice("import")
					if len(imports) > 0 {
						domain.Imports = imports
					}

					tablePrefix := c.String("table-prefix")
					if tablePrefix != "" {
						domain.Meta = generator.Meta{TablePrefix: tablePrefix}
					}

					data, err := yaml.Marshal(domain)
					if err != nil {
						return err
					}

					return ioutil.WriteFile(filepath.Join(c.String("output"), c.String("table")+".yml"), data, os.ModePerm)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}

}

// replaceExt replace ext for src
func replaceExt(src string, ext string) string {
	ext1 := path.Ext(src)

	return fmt.Sprintf("%s%s", src[:len(src)-len(ext1)], ext)
}

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

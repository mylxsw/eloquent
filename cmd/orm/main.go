package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mylxsw/eloquent/generator"
	"github.com/mylxsw/eloquent/generator/template"
	"gopkg.in/yaml.v2"
)

func main() {
	source := os.Args[1]
	matches, err := filepath.Glob(source)
	assertError(err)
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

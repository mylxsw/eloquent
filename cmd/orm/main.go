package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mylxsw/eloquent/generator"
	"github.com/mylxsw/go-toolkit/file"
	"gopkg.in/yaml.v2"
)

func main() {
	source := os.Args[1]
	matches, err := filepath.Glob(source)
	assertError(err)
	for _, m := range matches {
		dest := file.ReplaceExt(m, ".orm.go")

		input, err := ioutil.ReadFile(m)
		assertError(err)

		var domain generator.Domain
		assertError(yaml.Unmarshal(input, &domain))

		res, err := generator.ParseTemplate(generator.GetTemplate(), domain.Init())
		assertError(err)

		assertError(ioutil.WriteFile(dest, []byte(res), os.ModePerm))

		fmt.Println(dest)
	}
}

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var funcMap = template.FuncMap{
	"implode":        strings.Join,
	"trim":           strings.Trim,
	"trim_right":     strings.TrimRight,
	"trim_left":      strings.TrimLeft,
	"trim_space":     strings.TrimSpace,
	"lowercase":      strings.ToLower,
	"format":         fmt.Sprintf,
	"snake":          strcase.ToSnake,
	"camel":          strcase.ToCamel,
	"lower_camel":    strcase.ToLowerCamel,
}

func AddFunc(name string, f interface{}) {
	funcMap[name] = f
}

// ParseTemplate 模板解析
func ParseTemplate(templateContent string, data Domain) (string, error) {
	ctx := DomainContext{domain: data}
	ctx.Register()

	var buffer bytes.Buffer
	if err := template.Must(template.New("").Funcs(funcMap).Parse(templateContent)).Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

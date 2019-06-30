package template

func GetEntityWrapTemplate() string {
	return `
type {{ lower_camel $m.Name }}Wrap struct { {{ range $j, $f := fields $m.Definition }}	
	{{ camel $f.Name }} {{ wrap_type $f.Type }}{{ end }}
}

func (w {{ lower_camel $m.Name }}Wrap) To{{ camel $m.Name }} () {{ camel $m.Name }} {
	return {{ camel $m.Name }} {
		original: &{{ lower_camel $m.Name }}Original { {{ range $j, $f := fields $m.Definition }}
			{{ camel $f.Name }}: {{ unwrap_type $f.Name $f.Type }},{{ end }}
		},
		{{ range $j, $f := fields $m.Definition }}
		{{ camel $f.Name }}: {{ unwrap_type $f.Name $f.Type }},{{ end }}
	}
}
`
}

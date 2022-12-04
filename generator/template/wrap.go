package template

func GetEntityPlainTemplate() string {
	return `
type {{ camel $m.Name }} struct { {{ range $j, $f := fields $m.Definition }}	
	{{ camel $f.Name }} {{ $f.Type }}{{ end }}
}

func (w {{ camel $m.Name }}) To{{ camel $m.Name }}N(allows ...string) {{ camel $m.Name }}N {
	if len(allows) == 0 {
		return {{ camel $m.Name }}N {
			{{ range $j, $f := fields $m.Definition }}
			{{ camel $f.Name }}: {{ wrap_type (printf "w.%s" $f.Name) $f.Type }},{{ end }}
		}	
	}

	res := {{ camel $m.Name }}N{}
	for _, al := range allows {
		switch strcase.ToSnake(al) {
		{{ range $j, $f := fields $m.Definition }}
		case "{{ snake $f.Name }}":
			res.{{ camel $f.Name }} = {{ wrap_type (printf "w.%s" $f.Name) $f.Type }}{{ end }}
		default:
		}
	}

	return res
}

// As convert object to other type
// dst must be a pointer to struct
func (w {{ camel $m.Name }}) As(dst interface{}) error {
	return coll.CopyProperties(w, dst)
}


func (w *{{ camel $m.Name }}N) To{{ camel $m.Name }} () {{ camel $m.Name }} {
	return {{ camel $m.Name }} {
		{{ range $j, $f := fields $m.Definition }}
		{{ camel $f.Name }}: {{ unwrap_type $f.Name $f.Type }},{{ end }}
	}
}
`
}

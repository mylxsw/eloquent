package template

func GetEntityPlainTemplate() string {
	return `
type {{ camel $m.Name }}Plain struct { {{ range $j, $f := fields $m.Definition }}	
	{{ camel $f.Name }} {{ $f.Type }}{{ end }}
}

func (w {{ camel $m.Name }}Plain) To{{ camel $m.Name }}() {{ camel $m.Name }} {
	return {{ camel $m.Name }} {
		{{ range $j, $f := fields $m.Definition }}
		{{ camel $f.Name }}: {{ wrap_type (printf "w.%s" $f.Name) $f.Type }},{{ end }}
	}
}

// As convert object to other type
// dst must be a pointer to struct
func (w {{ camel $m.Name }}Plain) As(dst interface{}) error {
	return coll.CopyProperties(w, dst)
}


func (w *{{ camel $m.Name }}) To{{ camel $m.Name }}Plain () {{ camel $m.Name }}Plain {
	return {{ camel $m.Name }}Plain {
		{{ range $j, $f := fields $m.Definition }}
		{{ camel $f.Name }}: {{ unwrap_type $f.Name $f.Type }},{{ end }}
	}
}
`
}

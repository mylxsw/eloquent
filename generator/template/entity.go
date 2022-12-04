package template

func GetEntityTemplate() string {
	return `
// {{ camel $m.Name }}N is a {{ camel $m.Name }} object, all fields are nullable
type {{ camel $m.Name }}N struct {
	original *{{ lower_camel $m.Name }}Original
	{{ lower_camel $m.Name }}Model *{{ camel $m.Name }}Model

	{{ range $j, $f := fields $m.Definition }}
	{{ camel $f.Name }} {{ sql_field_type $f.Type }} {{ tag $f }}{{ end }}
}

// As convert object to other type
// dst must be a pointer to struct
func (inst *{{ camel $m.Name }}N) As(dst interface{}) error {
	return coll.CopyProperties(inst, dst)
}

// SetModel set model for {{ camel $m.Name }}
func (inst *{{ camel $m.Name }}N) SetModel({{ lower_camel $m.Name }}Model *{{ camel $m.Name }}Model) {
	inst.{{ lower_camel $m.Name }}Model = {{ lower_camel $m.Name }}Model
}

// {{ lower_camel $m.Name }}Original is an object which stores original {{ camel $m.Name }} from database
type {{ lower_camel $m.Name }}Original struct {
	{{ range $j, $f := fields $m.Definition }}
	{{ camel $f.Name }} {{ sql_field_type $f.Type }}{{ end }}
}

// Staled identify whether the object has been modified
func (inst *{{ camel $m.Name }}N) Staled(onlyFields ...string) bool {
	if inst.original == nil {
		inst.original = &{{ lower_camel $m.Name }}Original {}
	}

	if len(onlyFields) == 0 {
		{{ range $j, $f := fields $m.Definition }}
		if inst.{{ camel $f.Name }} != inst.original.{{ camel $f.Name }} {
			return true
		}{{ end }}
	} else {
		for _, f := range onlyFields {
			switch strcase.ToSnake(f) {
			{{ range $j, $f := fields $m.Definition }}
			case "{{ snake $f.Name }}":
				if inst.{{ camel $f.Name }} != inst.original.{{ camel $f.Name }} {
					return true
				}{{ end }}
			default:
			}
		}
	}

	return false
}

// StaledKV return all fields has been modified
func (inst *{{ camel $m.Name }}N) StaledKV(onlyFields ...string) query.KV {
	kv := make(query.KV, 0)

	if inst.original == nil {
		inst.original = &{{ lower_camel $m.Name }}Original {}
	}

	if len(onlyFields) == 0 {
		{{ range $j, $f := fields $m.Definition }}
		if inst.{{ camel $f.Name }} != inst.original.{{ camel $f.Name }} {
			kv["{{ snake $f.Name }}"] = inst.{{ camel $f.Name }}
		}{{ end }}
	} else {
		for _, f := range onlyFields {
			switch strcase.ToSnake(f) {
			{{ range $j, $f := fields $m.Definition }}
			case "{{ snake $f.Name }}":
				if inst.{{ camel $f.Name }} != inst.original.{{ camel $f.Name }} {
					kv["{{ snake $f.Name }}"] = inst.{{ camel $f.Name }}
				}{{ end }}
			default:
			}
		}
	}


	return kv
}

// Save create a new model or update it 
func (inst *{{ camel $m.Name }}N) Save(ctx context.Context, onlyFields ...string) error {
	if inst.{{ lower_camel $m.Name }}Model == nil {
		return query.ErrModelNotSet
	}

	id, _, err := inst.{{ lower_camel $m.Name }}Model.SaveOrUpdate(ctx, *inst, onlyFields...)
	if err != nil {
		return err 
	}

	inst.Id = null.IntFrom(id)
	return nil
}

// Delete remove a {{ $m.Name }}
func (inst *{{ camel $m.Name }}N) Delete(ctx context.Context) error {
	if inst.{{ lower_camel $m.Name }}Model == nil {
		return query.ErrModelNotSet
	}

	_, err := inst.{{ lower_camel $m.Name }}Model.DeleteById(ctx, inst.Id.Int64)
	if err != nil {
		return err 
	}

	return nil
}

// String convert instance to json string
func (inst *{{ camel $m.Name }}N) String() string {
	rs, _ := json.Marshal(inst)
	return string(rs)
}
`
}

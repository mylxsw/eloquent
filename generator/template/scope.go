package template


func GetScopeTemplate() string {
	return `
type {{ lower_camel $m.Name }}Scope struct {
	name  string
	apply func(builder query.Condition)
}

var {{ lower_camel $m.Name }}GlobalScopes = make([]{{ lower_camel $m.Name }}Scope, 0)
var {{ lower_camel $m.Name }}LocalScopes = make([]{{ lower_camel $m.Name }}Scope, 0)

// AddGlobalScopeFor{{ camel $m.Name }} assign a global scope to a model
func AddGlobalScopeFor{{ camel $m.Name }}(name string, apply func(builder query.Condition)) {
	{{ lower_camel $m.Name }}GlobalScopes = append({{ lower_camel $m.Name }}GlobalScopes, {{ lower_camel $m.Name }}Scope{name: name, apply: apply})
}

// AddLocalScopeFor{{ camel $m.Name }} assign a local scope to a model
func AddLocalScopeFor{{ camel $m.Name }}(name string, apply func(builder query.Condition)) {
	{{ lower_camel $m.Name }}LocalScopes = append({{ lower_camel $m.Name }}LocalScopes, {{ lower_camel $m.Name }}Scope{name: name, apply: apply})
}

func (m *{{ camel $m.Name }}Model) applyScope() query.Condition {
	scopeCond := query.ConditionBuilder()
	for _, g := range {{ lower_camel $m.Name }}GlobalScopes {
		if m.globalScopeEnabled(g.name) {
			g.apply(scopeCond)
		}
	}

	for _, s := range {{ lower_camel $m.Name }}LocalScopes {
		if m.localScopeEnabled(s.name) {
			s.apply(scopeCond)
		}
	}

	return scopeCond
}

func (m *{{ camel $m.Name }}Model) localScopeEnabled(name string) bool {
	for _, n := range m.includeLocalScopes {
		if name == n {
			return true
		}
	}

	return false
}

func (m *{{ camel $m.Name }}Model) globalScopeEnabled(name string) bool {
	for _, n := range m.excludeGlobalScopes {
		if name == n {
			return false
		}
	}
	
	return true
}
`
}
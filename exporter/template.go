package exporter

import (
	"bytes"
	_ "embed"
	"fmt"
	"reflect"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/marvinpeter95/reexporter/exporter/exports"
)

var (
	//go:embed templates/exported.gotpl
	codeTemplateSource string

	// codeTemplate is the parsed template for generation.
	codeTemplate = sync.OnceValue(func() *template.Template {
		fm := sprig.FuncMap()
		fm["mapProperty"] = templateMapProperty
		fm["parenthesize"] = templateParenthesize

		tpl, err := template.New("").Funcs(fm).Parse(codeTemplateSource)
		if err != nil {
			panic(err)
		}
		return tpl
	})
)

// renderTemplate renders the code template with the given export data.
func renderTemplate(data *exports.Exports) (string, error) {
	var buf bytes.Buffer
	if err := codeTemplate().Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// templateMapProperty extracts the property p from each element in the slice s.
// It supports slices of structs and maps.
func templateMapProperty(p string, s any) []any {
	sl := reflect.ValueOf(s)
	if sl.Kind() != reflect.Slice {
		return nil
	}

	var result []any
	for i := range sl.Len() {
		v := sl.Index(i)
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Map:
			result = append(result, v.MapIndex(reflect.ValueOf(p)).Interface())
		case reflect.Struct:
			if field := v.FieldByName(p); field.IsValid() {
				result = append(result, field.Interface())
			} else if m := v.MethodByName(p); m.IsValid() && m.Type().NumIn() == 0 && m.Type().NumOut() == 1 {
				result = append(result, m.Call(nil)[0].Interface())
			}
		}
	}

	return result
}

// templateParenthesize conditionally wraps the given string in parentheses.
// The last parameter is the string to wrap. Optional parameters before that:
// - A boolean condition (default: true) to determine whether to wrap.
// - A string of two characters representing the parentheses to use (default: "()").
func templateParenthesize(params ...any) string {
	if len(params) == 0 || len(params) > 3 {
		return ""
	}

	var (
		s           string
		condition   = true
		parentheses = "()"
	)

	if len(params) > 1 {
		if v, ok := params[0].(string); ok {
			parentheses = v
		}
	}

	// Determine the condition
	if len(params) > 2 {
		switch v := params[1].(type) {
		case bool:
			condition = v
		case int:
			condition = v != 0
		case string:
			condition = v != ""
		default:
			rv := reflect.ValueOf(v)
			condition = rv.IsValid() && !rv.IsZero()

			switch rv.Kind() {
			case reflect.Slice, reflect.Map, reflect.Array:
				condition = rv.Len() != 0
			}
		}
	}

	// Get the string to wrap
	last := params[len(params)-1]
	if str, ok := last.(string); ok {
		s = str
	} else {
		s = fmt.Sprintf("%v", last)
	}

	// Do not wrap if condition is false
	if !condition {
		return s
	}

	return string(parentheses[0]) + s + string(parentheses[1])
}

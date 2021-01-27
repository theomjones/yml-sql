package sql

import (
	"reflect"
	"text/template"
)

var funcMap = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len() - 1
	},
	"notLast": func(x int, a interface{}) bool {
		return x != reflect.ValueOf(a).Len() - 1
	},
}

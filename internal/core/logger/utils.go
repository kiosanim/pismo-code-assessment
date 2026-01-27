package logger

import (
	"reflect"
	"strings"
)

// ComponentFromStruct returns the type name of a struct without its package path.
// Ex: &AccountService{} → "AccountService"
func ComponentNameFromStruct(i any) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	parts := strings.Split(t.String(), ".")
	return parts[len(parts)-1]
}

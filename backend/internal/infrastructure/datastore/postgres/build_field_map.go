package postgres

import "reflect"

// buildFieldMap builds a map from db tag to field index
func buildFieldMap(t reflect.Type) map[string]int {
	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag != "" {
			fieldMap[tag] = i
		}
	}
	return fieldMap
}

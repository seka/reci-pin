package postgres

import (
	"reflect"

	"github.com/jackc/pgx/v5"
)

// scanAll scans all rows into a slice of T using reflection
func scanAll[T any](rows pgx.Rows) ([]T, error) {
	var results []T

	// Get type information for T
	typeOfT := reflect.TypeOf((*T)(nil)).Elem()

	// Build field map from db tags
	fieldMap := buildFieldMap(typeOfT)

	// Get column descriptions
	fieldDescriptions := rows.FieldDescriptions()

	for rows.Next() {
		// Create a new instance of T
		item := reflect.New(typeOfT).Elem()

		// Build scan destination slice
		scanDests := make([]interface{}, len(fieldDescriptions))
		for i, fd := range fieldDescriptions {
			colName := string(fd.Name)
			if fieldIndex, ok := fieldMap[colName]; ok {
				// Map to struct field
				scanDests[i] = item.Field(fieldIndex).Addr().Interface()
			} else {
				// No matching field, use dummy variable
				var dummy interface{}
				scanDests[i] = &dummy
			}
		}

		if err := rows.Scan(scanDests...); err != nil {
			return nil, err
		}

		results = append(results, item.Interface().(T))
	}

	return results, rows.Err()
}

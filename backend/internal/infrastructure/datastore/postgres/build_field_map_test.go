package postgres

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFieldMap(t *testing.T) {
	tests := []struct {
		name     string
		input    reflect.Type
		expected map[string]int
	}{
		{
			name: "basic struct with db tags",
			input: reflect.TypeOf(struct {
				ID   int64  `db:"id"`
				Name string `db:"name"`
			}{}),
			expected: map[string]int{
				"id":   0,
				"name": 1,
			},
		},
		{
			name: "struct with missing db tag",
			input: reflect.TypeOf(struct {
				ID       int64  `db:"id"`
				Name     string `db:"name"`
				Internal string // no db tag
			}{}),
			expected: map[string]int{
				"id":   0,
				"name": 1,
			},
		},
		{
			name:     "empty struct",
			input:    reflect.TypeOf(struct{}{}),
			expected: map[string]int{},
		},
		{
			name: "struct with only untagged fields",
			input: reflect.TypeOf(struct {
				Internal1 string
				Internal2 int
			}{}),
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFieldMap(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

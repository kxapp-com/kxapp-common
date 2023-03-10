package utilz

import (
	"fmt"
	"testing"
)

func TestGetMapArrayValue(t *testing.T) {
	m := map[string]any{
		"foo": []any{
			map[string]any{"bar": "baz"},
			map[string]any{"bar": "qux"},
		},
	}
	val, err := GetMapArrayValue(m, "foo", "1", "bar")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", val)
	}
}

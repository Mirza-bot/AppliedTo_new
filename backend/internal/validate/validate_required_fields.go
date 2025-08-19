package validate

import (
	"errors"
	"strings"
)

type Field struct {
	Name  string
	Value any
}

func Required(fields ...Field) error {
	for _, f := range fields {
		if isBlank(f.Value) {
			return errors.New("a " + f.Name + " is required")
		}
	}
	return nil
}

func isBlank(v any) bool {
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x) == ""
	case *string:
		return x == nil || strings.TrimSpace(*x) == ""
	case nil:
		return true
	default:
		// Everything else is considered present (numbers, bools, structs, etc.)
		return false
	}
}


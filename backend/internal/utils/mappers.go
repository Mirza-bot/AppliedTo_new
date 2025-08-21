package utils

import (
	"encoding/json"

	"gorm.io/datatypes"
)

// ---- helpers: Tags (JSONB) ----

func ToJSONTags(tags []string) datatypes.JSON {
	if tags == nil {
		return datatypes.JSON([]byte("[]"))
	}
	b, _ := json.Marshal(tags)
	return datatypes.JSON(b)
}

func FromJSONTags(j datatypes.JSON) ([]string, error) {
	if len(j) == 0 || string(j) == "null" {
		return []string{}, nil
	}
	var out []string
	if err := json.Unmarshal([]byte(j), &out); err != nil {
		return nil, err
	}
	return out, nil
}

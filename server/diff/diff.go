package diff

import (
	"encoding/json"
	"reflect"
)

type DiffType string

const (
	DiffAdded        DiffType = "ADDED"
	DiffRemoved      DiffType = "REMOVED"
	DiffTypeMismatch DiffType = "TYPE_MISMATCH"
	DiffValueChange  DiffType = "VALUE_CHANGE"
)

type DiffItem struct {
	Path     string   `json:"path"`
	Type     DiffType `json:"type"`
	Val1     interface{} `json:"val1,omitempty"`
	Val2     interface{} `json:"val2,omitempty"`
}

func CompareJSON(json1, json2 []byte) ([]DiffItem, error) {
	var val1, val2 interface{}

	if err := json.Unmarshal(json1, &val1); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(json2, &val2); err != nil {
		return nil, err
	}

	var diffs []DiffItem
	compareRecursive("", val1, val2, &diffs)
	return diffs, nil
}

func compareRecursive(path string, val1, val2 interface{}, diffs *[]DiffItem) {
	if val1 == nil && val2 == nil {
		return
	}

	if val1 == nil {
		*diffs = append(*diffs, DiffItem{Path: path, Type: DiffAdded, Val2: val2})
		return
	}

	if val2 == nil {
		*diffs = append(*diffs, DiffItem{Path: path, Type: DiffRemoved, Val1: val1})
		return
	}

	t1 := reflect.TypeOf(val1)
	t2 := reflect.TypeOf(val2)

	if t1 != t2 {
		*diffs = append(*diffs, DiffItem{Path: path, Type: DiffTypeMismatch, Val1: val1, Val2: val2})
		return
	}

	switch v1 := val1.(type) {
	case map[string]interface{}:
		v2 := val2.(map[string]interface{})
		// Check for removed or modified keys
		for k, val1Elem := range v1 {
			elemPath := k
			if path != "" {
				elemPath = path + "." + k
			}
			if val2Elem, exists := v2[k]; exists {
				compareRecursive(elemPath, val1Elem, val2Elem, diffs)
			} else {
				*diffs = append(*diffs, DiffItem{Path: elemPath, Type: DiffRemoved, Val1: val1Elem})
			}
		}
		// Check for added keys
		for k, val2Elem := range v2 {
			elemPath := k
			if path != "" {
				elemPath = path + "." + k
			}
			if _, exists := v1[k]; !exists {
				*diffs = append(*diffs, DiffItem{Path: elemPath, Type: DiffAdded, Val2: val2Elem})
			}
		}

	case []interface{}:
		v2 := val2.([]interface{})
		len1 := len(v1)
		len2 := len(v2)
		maxLen := len1
		if len2 > maxLen {
			maxLen = len2
		}

		for i := 0; i < maxLen; i++ {
			elemPath := reflect.ValueOf(i).String()
			if path != "" {
				elemPath = path + "[" + reflect.ValueOf(i).String() + "]"
			}
			if i < len1 && i < len2 {
				compareRecursive(elemPath, v1[i], v2[i], diffs)
			} else if i < len1 {
				*diffs = append(*diffs, DiffItem{Path: elemPath, Type: DiffRemoved, Val1: v1[i]})
			} else {
				*diffs = append(*diffs, DiffItem{Path: elemPath, Type: DiffAdded, Val2: v2[i]})
			}
		}

	default:
		if !reflect.DeepEqual(val1, val2) {
			*diffs = append(*diffs, DiffItem{Path: path, Type: DiffValueChange, Val1: val1, Val2: val2})
		}
	}
}

package main

import (
	"fmt"
	"reflect"
)

// MergeParams merges source map into target map, preserving types (including booleans and nested maps).
func MergeParams(target, source map[string]any) map[string]any {
	if target == nil {
		target = make(map[string]any)
	}
	for k, v := range source {
		if v == nil {
			continue
		}
		tg, targetExists := target[k]
		if !targetExists {
			target[k] = cloneValue(v)
			continue
		}

		// If both are maps, deep merge them
		sourceMap, sourceIsMap := toMap(v)
		targetMap, targetIsMap := toMap(tg)
		if sourceIsMap && targetIsMap {
			target[k] = MergeParams(targetMap, sourceMap)
		} else {
			// Otherwise, target takes precedence, do not overwrite
		}
	}
	return target
}

func toMap(v any) (map[string]any, bool) {
	if m, ok := v.(map[string]any); ok {
		return m, true
	}
	// Handle other map types if necessary
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Map {
		res := make(map[string]any)
		for _, key := range val.MapKeys() {
			res[fmt.Sprint(key.Interface())] = val.MapIndex(key).Interface()
		}
		return res, true
	}
	return nil, false
}

func cloneValue(v any) any {
	if m, ok := toMap(v); ok {
		cloned := make(map[string]any)
		for mk, mv := range m {
			cloned[mk] = cloneValue(mv)
		}
		return cloned
	}
	return v
}

func main() {
	// Simulate cascade inheritance test
	parentCascade := map[string]any{
		"boolean_true":  true,
		"boolean_false": false,
		"nested": map[string]any{
			"bool_val": true,
			"str_val":  "hello",
		},
	}

	childParams := map[string]any{
		"existing_param": "value",
	}

	merged := MergeParams(childParams, parentCascade)

	// Verify types
	bt, ok1 := merged["boolean_true"].(bool)
	bf, ok2 := merged["boolean_false"].(bool)
	
	nestedMap, ok3 := merged["nested"].(map[string]any)
	var nestedBool bool
	var ok4 bool
	if ok3 {
		nestedBool, ok4 = nestedMap["bool_val"].(bool)
	}

	if ok1 && bt && ok2 && !bf && ok3 && ok4 && nestedBool {
		fmt.Println("Success: Boolean values in cascade inheritance preserved correctly!")
	} else {
		fmt.Println("Failure: Boolean values were coerced or lost during cascade inheritance.")
	}
}

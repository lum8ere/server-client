package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type ANY_DATA map[string]any

func AnyDataRef(anyData ANY_DATA) *ANY_DATA {
	return &anyData
}

// DeepCopy creates a full deep copy of ANY_DATA
func (j *ANY_DATA) DeepCopy() (ANY_DATA, error) {
	if j == nil {
		return nil, nil
	}

	// Serialize to JSON and back to ensure full deep copy
	data, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ANY_DATA: %w", err)
	}

	var copy ANY_DATA
	if err := json.Unmarshal(data, &copy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ANY_DATA: %w", err)
	}

	return copy, nil
}

func (a ANY_DATA) ToMap() map[string]any {
	return map[string]any(a)
}

func (a *ANY_DATA) Scan(value interface{}) error {
	if value == nil {
		// Handle NULL values
		*a = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("ANY_DATA.Scan: expected []byte or string, got %T", value)
	}
	// Unmarshal JSON to the map
	return json.Unmarshal(bytes, a)
}

func (a *ANY_DATA) GetIntValue(argName string) (int64, error) {
	if a == nil {
		// Return 0, NULL if the map is nil
		return 0, nil
	}

	val, ok := (*a)[argName]
	if !ok {
		// Если ключ не найден, возвращаем ошибку
		return 0, fmt.Errorf("key '%s' not found", argName)
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if f == float64(int64(f)) {
			return int64(f), nil
		}
	case reflect.String:
		idInt, err := strconv.Atoi(v.String())
		if err == nil {
			return int64(idInt), nil
		}
	case reflect.Map:
		idVal := v.MapIndex(reflect.ValueOf("Id"))
		if idVal.IsValid() {
			id := idVal.Interface()

			v := reflect.ValueOf(id)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return v.Int(), nil
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return int64(v.Uint()), nil
			case reflect.Float32, reflect.Float64:
				f := v.Float()
				if f == float64(int64(f)) {
					return int64(f), nil
				}
			case reflect.String:
				idInt, err := strconv.Atoi(v.String())
				if err == nil {
					return int64(idInt), nil
				}
			}
		}
	}

	return 0, nil
}

func (a *ANY_DATA) GetBoolValue(argName string) (bool, bool) {
	if a == nil {
		return false, false
	}

	val, ok := (*a)[argName]
	if !ok {
		return false, false
	}

	switch v := val.(type) {
	case string:
		if v == "true" {
			return true, true
		}

		if v == "false" {
			return false, true
		}

		return false, false
	case bool:
		return v, true
	default:
		return false, false
	}
}

func (a *ANY_DATA) GetStringValue(argName string) (string, bool) {
	if a == nil {
		return "", false
	}

	val, ok := (*a)[argName]
	if !ok {
		return "", false
	}

	switch v := val.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}
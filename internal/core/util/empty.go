package utils

import "reflect"

// IsEmpty 检查值是否为空（零值）。
// 字符串为空、数值为 0、bool 为 false、slice/map/chan 长度为 0 时返回 true。
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch value := value.(type) {
	case int:
		return value == 0
	case int8:
		return value == 0
	case int16:
		return value == 0
	case int32:
		return value == 0
	case int64:
		return value == 0
	case uint:
		return value == 0
	case uint8:
		return value == 0
	case uint16:
		return value == 0
	case uint32:
		return value == 0
	case uint64:
		return value == 0
	case float32:
		return value == 0
	case float64:
		return value == 0
	case bool:
		return value == false
	case string:
		return value == ""
	case []byte:
		return len(value) == 0
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
			return rv.Len() == 0
		case reflect.Func, reflect.Ptr, reflect.Interface, reflect.UnsafePointer:
			if rv.IsNil() {
				return true
			}
		}
	}
	return false
}

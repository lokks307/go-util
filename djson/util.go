package djson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func indent(v interface{}) string {
	indent, _ := json.MarshalIndent(v, "", "   ")
	return string(indent)
}

func _string(v interface{}) string {
	indent, _ := json.Marshal(v)
	return string(indent)
}

func _stringBase(v interface{}) (string, bool) {

	fmt.Println("_stringBase=", reflect.TypeOf(v).String())

	if v == nil {
		return "", false
	}

	switch t := v.(type) {
	case string:
		return t, true
	case bool:
		return strconv.FormatBool(t), true
	case int:
		return strconv.FormatInt(int64(t), 10), true
	case uint:
		return strconv.FormatUint(uint64(t), 10), true
	case int8:
		return strconv.FormatInt(int64(t), 10), true
	case uint8:
		return strconv.FormatUint(uint64(t), 10), true
	case int16:
		return strconv.FormatInt(int64(t), 10), true
	case uint16:
		return strconv.FormatUint(uint64(t), 10), true
	case int32:
		return strconv.FormatInt(int64(t), 10), true
	case uint32:
		return strconv.FormatUint(uint64(t), 10), true
	case int64:
		return strconv.FormatInt(t, 10), true
	case uint64:
		return strconv.FormatUint(t, 10), true
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32), true
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), true
	}

	return "", false
}

func IsBaseType(v interface{}) bool {
	switch v.(type) {
	case string:
		return true
	case bool:
		return true
	case int:
		return true
	case uint:
		return true
	case int8:
		return true
	case uint8:
		return true
	case int16:
		return true
	case uint16:
		return true
	case int32:
		return true
	case uint32:
		return true
	case int64:
		return true
	case uint64:
		return true
	case float32:
		return true
	case float64:
		return true
	}

	return false
}

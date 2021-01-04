package djson

import (
	"fmt"
	"reflect"
	"strings"

	gov "github.com/asaskevich/govalidator"
)

func ConverMapToObject(dmap map[string]interface{}) *DO {
	nObj := NewObject()
	for k, v := range dmap {
		nObj.Put(k, v)
	}
	return nObj
}

func ConvertSliceToArray(dslice []interface{}) *DA {
	nArr := NewArray()
	nArr.Put(dslice)
	return nArr
}

func ConverObjectToMap(obj *DO) map[string]interface{} {
	wMap := make(map[string]interface{})

	for k, v := range obj.Map {
		switch t := v.(type) {
		case DA:
			wMap[k] = ConvertArrayToSlice(&t)
		case DO:
			wMap[k] = ConverObjectToMap(&t)
		case *DA:
			wMap[k] = ConvertArrayToSlice(t)
		case *DO:
			wMap[k] = ConverObjectToMap(t)
		default:
			wMap[k] = v
		}
	}

	return wMap
}

func ConvertArrayToSlice(arr *DA) []interface{} {

	wArray := make([]interface{}, 0)

	for idx := range arr.Element {
		switch t := arr.Element[idx].(type) {
		case DA:
			wArray = append(wArray, ConvertArrayToSlice(&t))
		case DO:
			wArray = append(wArray, ConverObjectToMap(&t))
		case *DA:
			wArray = append(wArray, ConvertArrayToSlice(t))
		case *DO:
			wArray = append(wArray, ConverObjectToMap(t))
		default:
			wArray = append(wArray, t)
		}
	}

	return wArray
}

func getStringBase(v interface{}) (string, bool) {
	if v == nil {
		return "nil", true
	}

	if IsInTypes(v, "string", "bool", "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64", "float32", "float64") {
		return fmt.Sprintf("%v", v), true
	}

	return "", false
}

func getBoolBase(v interface{}) (bool, bool) {
	if IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		intVal, _ := gov.ToInt(v)
		if intVal == 0 {
			return false, true
		}
	}

	if IsInTypes(v, "string") {
		if strVal, ok := v.(string); ok {
			if strings.EqualFold(strVal, "true") {
				return true, true
			} else if strings.EqualFold(strVal, "false") {
				return false, true
			}
		}
	}

	if IsInTypes(v, "bool") {
		if boolVal, ok := v.(bool); ok {
			return boolVal, true
		}
	}

	return false, false
}

func getFloatBase(v interface{}) (float64, bool) {
	if floatVal, err := gov.ToFloat(v); err != nil {
		return 0, false
	} else {
		return floatVal, true
	}
}

func getIntBase(v interface{}) (int64, bool) {
	if intVal, err := gov.ToInt(v); err != nil {
		return 0, false
	} else {
		return intVal, true
	}
}

func IsBaseType(v interface{}) bool {
	return IsInTypes(v, "string", "bool", "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64", "float32", "float64")
}

func IsInTypes(v interface{}, types ...string) bool {
	var vTypeStr string
	if v == nil {
		vTypeStr = "nil"
	} else {
		vTypeStr = reflect.TypeOf(v).String()
	}

	for idx := range types {
		if vTypeStr == types[idx] {
			return true
		}
	}

	return false
}

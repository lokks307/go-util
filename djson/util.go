package djson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

	if IsInTypes(v, "string", "bool", "float32", "float64") {
		return fmt.Sprintf("%v", v), true
	}

	if IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		return fmt.Sprintf("%d", v), true
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

func IsIntType(v interface{}) bool {
	return IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64")
}

func IsFloatType(v interface{}) bool {
	return IsInTypes(v, "float32", "float64")
}

func IsBoolType(v interface{}) bool {
	return IsInTypes(v, "bool")
}

func IsStringType(v interface{}) bool {
	return IsInTypes(v, "string")
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

func ParseToObject(doc string) (*DO, error) {
	var data map[string]interface{}

	d := json.NewDecoder(strings.NewReader(doc))
	d.UseNumber()

	if err := d.Decode(&data); err != nil {
		return nil, errors.New("not Object")
	}

	return ParseObject(data), nil

}

func ParseToArray(doc string) (*DA, error) {
	var data []interface{}

	d := json.NewDecoder(strings.NewReader(doc))
	d.UseNumber()

	if err := d.Decode(&data); err != nil {
		return nil, errors.New("not Array")
	}

	return ParseArray(data), nil
}

func ParseObject(data map[string]interface{}) *DO {
	obj := NewObject()
	for k, v := range data {
		if IsBaseType(v) {
			obj.Put(k, v)
			continue
		}

		if n, ok := v.(json.Number); ok {
			if i, err := n.Int64(); err == nil {
				obj.Put(k, i)
				continue
			}
			if f, err := n.Float64(); err == nil {
				obj.Put(k, f)
				continue
			}
		}

		switch tValue := v.(type) {
		case []interface{}: // Array
			nArr := ParseArray(tValue)
			obj.Put(k, nArr)
		case map[string]interface{}: // Object
			nObj := ParseObject(tValue)
			obj.Put(k, nObj)
		case nil: // null
			obj.Put(k, nil)
		}
	}

	return obj
}

func ParseArray(data []interface{}) *DA {
	arr := NewArray()

	for idx := range data {
		if IsBaseType(data[idx]) {
			arr.Put(data[idx])
			continue
		}

		if n, ok := data[idx].(json.Number); ok {
			if i, err := n.Int64(); err == nil {
				arr.Put(i)
				continue
			}
			if f, err := n.Float64(); err == nil {
				arr.Put(f)
				continue
			}
		}

		switch tValue := data[idx].(type) {
		case []interface{}: // Array
			nArr := ParseArray(tValue)
			arr.PutAsArray(nArr)
		case map[string]interface{}: // Object
			nObj := ParseObject(tValue)
			arr.Put(nObj)
		case nil: // null
			arr.Put(nil)
		}
	}

	return arr
}

func PathTokenizer(path string) []interface{} {
	rstack := NewRuneStack()
	token := make([]rune, 0)
	inTokens := make([]string, 0)

	prev := rune(0)
	var depthL int

	for _, each := range path {

		peek := rstack.Peek()

		if depthL == 0 {
			if each == '[' && prev != '\\' {
				rstack.Push(each)
				token = make([]rune, 0)
				depthL = 1
			} else {
				token = append(token, each)
			}
		} else if depthL == 1 {
			if peek == '[' && each == ']' && prev != '\\' {
				if len(token) > 0 {
					inTokens = append(inTokens, string(token))
					token = make([]rune, 0)
				}
				rstack.Pop()
				depthL = 0
			} else if (each == '"' || each == '\'') && prev != '\\' {
				rstack.Push(each)
				depthL = 2
			} else {
				token = append(token, each)
			}
		} else if depthL == 2 {

			if (peek == '"' && each == '"' && prev != '\\') || (peek == '\'' && each == '\'' && prev != '\\') {
				if len(token) > 0 {
					inTokens = append(inTokens, string(token))
					token = make([]rune, 0)
				}
				rstack.Pop()
				depthL = 1
			} else {
				token = append(token, each)
			}
		}

		prev = each
	}

	outTokens := make([]interface{}, 0)
	for idx := range inTokens {
		if intVal, err := strconv.Atoi(inTokens[idx]); err == nil {
			outTokens = append(outTokens, intVal)
		} else {
			outTokens = append(outTokens, inTokens[idx])
		}
	}

	return outTokens
}

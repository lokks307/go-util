package djson

import (
	"strconv"
	"strings"

	gov "github.com/asaskevich/govalidator"
)

const (
	JSON_NULL   = 0
	JSON_OBJECT = 1
	JSON_ARRAY  = 2
	JSON_STRING = 3
	JSON_INT    = 4
	JSON_FLOAT  = 5
	JSON_BOOL   = 6
)

type DJSON struct {
	Object   *O
	Array    *A
	String   string
	Int      int64
	Float    float64
	Bool     bool
	JsonType int
}

func NewDJSON() *DJSON {
	return &DJSON{
		JsonType: JSON_NULL,
	}
}

func (m *DJSON) Parse(doc string) *DJSON {

	if m.JsonType != JSON_NULL {
		return m
	}

	tdoc := strings.TrimSpace(doc)
	if tdoc == "" {
		return m
	}

	var err error

	if tdoc[0] == '{' {
		m.Object, err = ParseToObject(tdoc)
		if err == nil {
			m.JsonType = JSON_OBJECT
		}
	} else if tdoc[0] == '[' {
		m.Array, err = ParseToArray(tdoc)
		if err == nil {
			m.JsonType = JSON_ARRAY
		}
	} else {
		if strings.EqualFold(tdoc, "null") {
			m.JsonType = JSON_NULL
		} else if strings.EqualFold(tdoc, "true") || strings.EqualFold(tdoc, "false") {
			m.JsonType = JSON_BOOL
			m.Bool, _ = gov.ToBoolean(tdoc)
		} else {
			if gov.IsNumeric(tdoc) {
				if gov.IsInt(tdoc) {
					m.Int, _ = strconv.ParseInt(tdoc, 10, 64)
					m.JsonType = JSON_INT
				} else {
					m.Float, _ = strconv.ParseFloat(tdoc, 64)
					m.JsonType = JSON_FLOAT
				}
			} else {
				m.String = tdoc
				m.JsonType = JSON_STRING
			}
		}
	}

	return m
}

func (m *DJSON) Put(v interface{}) *DJSON {
	if m.JsonType != JSON_NULL || v == nil {
		return m
	}

	if IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		m.Int, _ = getIntBase(v)
		m.JsonType = JSON_INT
		return m
	}

	if IsInTypes(v, "float32", "float64") {
		m.Float, _ = getFloatBase(v)
		m.JsonType = JSON_FLOAT
		return m
	}

	if IsInTypes(v, "bool") {
		m.Bool, _ = getBoolBase(v)
		m.JsonType = JSON_BOOL
		return m
	}

	switch t := v.(type) {
	case map[string]interface{}:
		m.Object = ConverMapToObject(t)
		m.JsonType = JSON_OBJECT
	case []interface{}:
		m.Array = ConvertSliceToArray(t)
		m.JsonType = JSON_ARRAY
	case _Object:
		m.Object = ConverMapToObject(t)
		m.JsonType = JSON_OBJECT
	case _Array:
		m.Array = ConvertSliceToArray(t)
		m.JsonType = JSON_ARRAY
	}

	return m
}

func (m *DJSON) PutAsArray(value ...interface{}) *DJSON {
	if m.JsonType == JSON_NULL {
		m.Array = NewArray()
		m.JsonType = JSON_ARRAY
	}

	if m.JsonType == JSON_ARRAY {
		m.Array.Put(value)
	}

	return m
}

func (m *DJSON) PutAsObject(key string, value interface{}) *DJSON {
	if m.JsonType == JSON_NULL {
		m.Object = NewObject()
		m.JsonType = JSON_OBJECT
	}

	if m.JsonType == JSON_OBJECT {
		m.Object.Put(key, value)
	}

	return m
}

func (m *DJSON) Get(key ...interface{}) (*DJSON, bool) {
	if len(key) == 0 {
		return m, true
	} else {

		r := NewDJSON()
		var element interface{}
		var retOk bool

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				element, retOk = m.Object.Get(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				element, retOk = m.Array.Get(tkey)
			}
		}

		if !retOk {
			return nil, false
		}

		switch t := element.(type) {
		case nil:
			r.JsonType = JSON_NULL
		case string:
			r.String = t
			r.JsonType = JSON_STRING
		case bool:
			r.Bool = t
			r.JsonType = JSON_BOOL
		case uint8:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case int8:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case uint16:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case int16:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case uint32:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case int32:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case uint64:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case int64:
			r.Int = t
			r.JsonType = JSON_INT
		case uint:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case int:
			r.Int = int64(t)
			r.JsonType = JSON_INT
		case float32:
			r.Float = float64(t)
			r.JsonType = JSON_FLOAT
		case float64:
			r.Float = t
			r.JsonType = JSON_FLOAT
		case A:
			r.Array = &t
			r.JsonType = JSON_ARRAY
		case O:
			r.Object = &t
			r.JsonType = JSON_OBJECT
		case *A:
			r.Array = t
			r.JsonType = JSON_ARRAY
		case *O:
			r.Object = t
			r.JsonType = JSON_OBJECT
		default:
			return nil, false
		}

		return r, true
	}
}

func (m *DJSON) GetAsObject(key ...interface{}) (*DJSON, bool) {

	if m.JsonType == JSON_STRING || m.JsonType == JSON_INT || m.JsonType == JSON_FLOAT {
		return nil, false
	}

	if len(key) == 0 {
		if m.JsonType == JSON_OBJECT {
			return m, true
		}
	} else {

		var ok bool
		var nObj *O

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				nObj, ok = m.Object.GetAsObject(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				nObj, ok = m.Array.GetAsObject(tkey)
			}
		}

		if !ok {
			return nil, false
		}

		if nObj != nil {
			return &DJSON{
				Object:   nObj,
				Array:    nil,
				JsonType: JSON_OBJECT,
			}, true
		}
	}

	return nil, false
}

func (m *DJSON) GetAsArray(key ...interface{}) (*DJSON, bool) {

	if m.JsonType == JSON_STRING || m.JsonType == JSON_INT || m.JsonType == JSON_FLOAT {
		return nil, false
	}

	if len(key) == 0 {
		if m.JsonType == JSON_ARRAY {
			return m, true
		}
	} else {

		var ok bool
		var nArr *A

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				nArr, ok = m.Object.GetAsArray(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				nArr, ok = m.Array.GetAsArray(tkey)
			}
		}

		if !ok {
			return nil, false
		}

		if nArr != nil {
			return &DJSON{
				Object:   nil,
				Array:    nArr,
				JsonType: JSON_ARRAY,
			}, true
		}

	}

	return nil, false
}

func (m *DJSON) GetAsInt(key ...interface{}) int64 {
	if len(key) == 0 {

		switch m.JsonType {
		case JSON_NULL:
			return 0
		case JSON_STRING:
			return 0
		case JSON_INT:
			return m.Int
		case JSON_FLOAT:
			return int64(m.Float)
		}

	} else {
		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				return m.Object.GetAsInt(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				return m.Array.GetAsInt(tkey)
			}
		}
	}

	return 0
}

func (m *DJSON) GetAsBool(key ...interface{}) bool {
	if len(key) == 0 {

		switch m.JsonType {
		case JSON_NULL:
			return false
		case JSON_STRING:
			return false
		case JSON_INT:
			return false
		case JSON_FLOAT:
			return false
		case JSON_BOOL:
			return m.Bool
		}

	} else {
		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				return m.Object.GetAsBool(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				return m.Array.GetAsBool(tkey)
			}
		}
	}

	return false
}

func (m *DJSON) GetAsFloat(key ...interface{}) float64 {
	if len(key) == 0 {

		switch m.JsonType {
		case JSON_NULL:
			return 0
		case JSON_STRING:
			return 0
		case JSON_INT:
			return float64(m.Int)
		case JSON_FLOAT:
			return m.Float
		}

	} else {
		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				return m.Object.GetAsFloat(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				return m.Array.GetAsFloat(tkey)
			}
		}
	}

	return 0
}

func (m *DJSON) GetAsString(key ...interface{}) string {

	if len(key) == 0 {

		switch m.JsonType {
		case JSON_NULL:
			return "null"
		case JSON_STRING:
			return m.String
		case JSON_INT:
			intStr, ok := getStringBase(m.Int)
			if !ok {
				return ""
			}
			return intStr
		case JSON_FLOAT:
			floatStr, ok := getStringBase(m.Float)
			if !ok {
				return ""
			}
			return floatStr
		case JSON_BOOL:
			return gov.ToString(m.Bool)
		case JSON_OBJECT:
			return m.Object.ToString()
		case JSON_ARRAY:
			return m.Array.ToString()
		}

	} else {
		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				return m.Object.GetAsString(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				return m.Array.GetAsString(tkey)
			}
		}
	}

	return ""
}

func (m *DJSON) IsBool() bool {
	return m.JsonType == JSON_BOOL
}

func (m *DJSON) IsInt() bool {
	return m.JsonType == JSON_INT
}

func (m *DJSON) IsNumeric() bool {
	return m.JsonType == JSON_FLOAT || m.JsonType == JSON_INT
}

func (m *DJSON) IsFloat() bool {
	return m.JsonType == JSON_FLOAT
}

func (m *DJSON) IsString() bool {
	return m.JsonType == JSON_STRING
}

func (m *DJSON) IsNull() bool {
	return m.JsonType == JSON_NULL
}

func (m *DJSON) IsObject() bool {
	return m.JsonType == JSON_OBJECT
}

func (m *DJSON) IsArray() bool {
	return m.JsonType == JSON_ARRAY
}

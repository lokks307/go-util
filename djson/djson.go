package djson

import (
	"reflect"
	"regexp"
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
	Object   *DO
	Array    *DA
	String   string
	Int      int64
	Float    float64
	Bool     bool
	JsonType int
}

var XPathRegExp *regexp.Regexp

func init() {
}

func NewDJSON(v ...int) *DJSON {

	dj := DJSON{
		JsonType: JSON_NULL,
	}

	if len(v) == 1 {
		switch v[0] {
		case JSON_OBJECT:
			dj.Object = NewObject()
			dj.JsonType = JSON_OBJECT
		case JSON_ARRAY:
			dj.Array = NewArray()
			dj.JsonType = JSON_ARRAY
		case JSON_STRING:
			dj.JsonType = JSON_STRING
		case JSON_INT:
			dj.JsonType = JSON_INT
		case JSON_FLOAT:
			dj.JsonType = JSON_FLOAT
		case JSON_BOOL:
			dj.JsonType = JSON_BOOL
		}
	}

	return &dj
}

func (m *DJSON) SetAsObject() *DJSON {
	m.Object = NewObject()
	m.Array = nil
	m.JsonType = JSON_OBJECT

	return m
}

func (m *DJSON) SetAsArray() *DJSON {
	m.Array = NewArray()
	m.Object = nil
	m.JsonType = JSON_ARRAY

	return m
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

func (m *DJSON) Put(v ...interface{}) *DJSON {

	if len(v) == 0 {
		return m
	}

	if len(v) == 2 {

		if key, ok := v[0].(string); ok {
			m.PutAsObject(key, v[1])
		}

		return m
	}

	if v[0] == nil {
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_NULL
		return m
	}

	if IsInTypes(v[0], "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		m.Int, _ = getIntBase(v)
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_INT
		return m
	}

	if IsInTypes(v[0], "float32", "float64") {
		m.Float, _ = getFloatBase(v[0])
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_FLOAT
		return m
	}

	if IsInTypes(v[0], "bool") {
		m.Bool, _ = getBoolBase(v[0])
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_BOOL
		return m
	}

	if IsInTypes(v[0], "string") {
		m.String, _ = getStringBase(v[0])
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_STRING
		return m
	}

	switch t := v[0].(type) {
	case map[string]interface{}:
		if m.JsonType == JSON_OBJECT {
			for key := range t {
				m.Object.Put(key, t[key])
			}
		} else {
			m.Object = ConverMapToObject(t)
			m.Array = nil
			m.JsonType = JSON_OBJECT
		}
	case Object:
		if m.JsonType == JSON_OBJECT {
			for key := range map[string]interface{}(t) {
				m.Object.Put(key, t[key])
			}
		} else {
			m.Object = ConverMapToObject(t)
			m.Array = nil
			m.JsonType = JSON_OBJECT
		}
	case *DO:
		if m.JsonType == JSON_OBJECT {
			for key := range t.Map {
				m.Object.Put(key, t.Map[key])
			}
		} else {
			m.Object = t
			m.Array = nil
			m.JsonType = JSON_OBJECT
		}
	case DO:
		if m.JsonType == JSON_OBJECT {
			for key := range t.Map {
				m.Object.Put(key, t.Map[key])
			}
		} else {
			m.Object = &t
			m.Array = nil
			m.JsonType = JSON_OBJECT
		}
	case []interface{}:
		if m.JsonType == JSON_ARRAY {
			m.Array.Put(t)
		} else {
			m.Array = ConvertSliceToArray(t)
			m.Object = nil
			m.JsonType = JSON_ARRAY
		}

	case Array:
		if m.JsonType == JSON_ARRAY {
			m.Array.Put([]interface{}(t))
		} else {
			m.Array = ConvertSliceToArray(t)
			m.Object = nil
			m.JsonType = JSON_ARRAY
		}
	case *DA:
		if m.JsonType == JSON_ARRAY {
			m.Array.Put(t.Element)
		} else {
			m.Array = t
			m.Object = nil
			m.JsonType = JSON_ARRAY
		}
	case DA:
		if m.JsonType == JSON_ARRAY {
			m.Array.Put(t.Element)
		} else {
			m.Array = &t
			m.Object = nil
			m.JsonType = JSON_ARRAY
		}
	case DJSON:
		m = &t
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

func (m *DJSON) Remove(key interface{}) *DJSON {
	switch tkey := key.(type) {
	case string:
		if m.JsonType == JSON_OBJECT {
			m.Object.Remove(tkey)
		}
	case int:
		if m.JsonType == JSON_ARRAY {
			m.Array.Remove(tkey)
		}
	}

	return m
}

func (m *DJSON) GetAsInterface(key ...interface{}) interface{} {
	if len(key) == 0 {
		switch m.JsonType {
		case JSON_NULL:
			return nil
		case JSON_STRING:
			return m.String
		case JSON_BOOL:
			return m.Bool
		case JSON_INT:
			return m.Int
		case JSON_FLOAT:
			return m.Float
		case JSON_OBJECT:
			return m.Object
		case JSON_ARRAY:
			return m.Array
		}

		return nil
	} else {

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				if obj, ok := m.Object.Get(tkey); ok {
					return obj
				}
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				if arr, ok := m.Array.Get(tkey); ok {
					return arr
				}
			}
		}
	}

	return nil
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

		eVal := reflect.ValueOf(element)

		switch t := element.(type) {
		case nil:
			r.JsonType = JSON_NULL
		case string:
			r.String = t
			r.JsonType = JSON_STRING
		case bool:
			r.Bool = t
			r.JsonType = JSON_BOOL
		case uint8, uint16, uint32, uint64, uint:
			intVal := int64(eVal.Uint())
			r.Int = intVal
			r.JsonType = JSON_INT
		case int8, int16, int32, int64, int:
			intVal := eVal.Int()
			r.Int = intVal
			r.JsonType = JSON_INT
		case float32, float64:
			floatVal := eVal.Float()
			r.Float = floatVal
			r.JsonType = JSON_FLOAT
		case DA:
			r.Array = &t
			r.JsonType = JSON_ARRAY
		case DO:
			r.Object = &t
			r.JsonType = JSON_OBJECT
		case *DA:
			r.Array = t
			r.JsonType = JSON_ARRAY
		case *DO:
			r.Object = t
			r.JsonType = JSON_OBJECT
		default:
			return nil, false
		}

		return r, true
	}
}

// The DJSON as return shared Object.

func (m *DJSON) GetAsObject(key ...interface{}) (*DJSON, bool) {

	if m.JsonType != JSON_OBJECT && m.JsonType != JSON_ARRAY {
		return nil, false
	}

	if len(key) == 0 {
		if m.JsonType == JSON_OBJECT {
			return m, true
		}
	} else {

		var ok bool
		var newObject *DO

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				newObject, ok = m.Object.GetAsObject(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				newObject, ok = m.Array.GetAsObject(tkey)
			}
		}

		if !ok {
			return nil, false
		}

		if newObject != nil {
			return &DJSON{
				Object:   newObject,
				Array:    nil,
				JsonType: JSON_OBJECT,
			}, true
		}
	}

	return nil, false
}

// The DJSON as return shared Array.

func (m *DJSON) GetAsArray(key ...interface{}) (*DJSON, bool) {

	if m.JsonType != JSON_OBJECT && m.JsonType != JSON_ARRAY {
		return nil, false
	}

	if len(key) == 0 {
		if m.JsonType == JSON_ARRAY {
			return m, true
		}
	} else {

		var ok bool
		var newArray *DA

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				newArray, ok = m.Object.GetAsArray(tkey)
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				newArray, ok = m.Array.GetAsArray(tkey)
			}
		}

		if !ok {
			return nil, false
		}

		if newArray != nil {
			return &DJSON{
				Object:   nil,
				Array:    newArray,
				JsonType: JSON_ARRAY,
			}, true
		}

	}

	return nil, false
}

func (m *DJSON) GetAsInt(key ...interface{}) int64 {

	if len(key) == 0 {

		switch m.JsonType {
		case JSON_ARRAY, JSON_OBJECT, JSON_NULL:
			return 0
		case JSON_BOOL:
			if m.Bool {
				return 1
			}
			return 0
		case JSON_STRING:
			if iVal, err := strconv.ParseInt(m.String, 10, 64); err == nil {
				return iVal
			}
			return 0
		case JSON_INT:
			return m.Int
		case JSON_FLOAT:
			return int64(m.Float)
		}

	} else {

		var hasDefaultValue bool
		var defaultValue int64

		if len(key) >= 2 {
			defaultValue, hasDefaultValue = key[1].(int64)
		}

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				if m.Object.HasKey(tkey) {
					if iVal, ok := m.Object.GetAsInt(tkey); ok {
						return iVal
					}
				}
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				if m.Array.Size() > tkey {
					if iVal, ok := m.Array.GetAsInt(tkey); ok {
						return iVal
					}
				}
			}
		}

		if hasDefaultValue {
			return defaultValue
		}
	}

	return 0 // zero value
}

func (m *DJSON) GetAsBool(key ...interface{}) bool {
	if len(key) == 0 {

		switch m.JsonType {
		case JSON_NULL, JSON_FLOAT, JSON_ARRAY, JSON_OBJECT:
			return false
		case JSON_STRING:
			if strings.EqualFold(m.String, "true") {
				return true
			}
			return false
		case JSON_INT:
			return m.Int == 1
		case JSON_BOOL:
			return m.Bool
		}

	} else {

		var hasDefaultValue bool
		var defaultValue bool

		if len(key) >= 2 {
			defaultValue, hasDefaultValue = key[1].(bool)
		}

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {

				if m.Object.HasKey(tkey) {
					if bVal, ok := m.Object.GetAsBool(tkey); ok {
						return bVal
					}
				}
			}
		case int:

			if m.Array.Size() > tkey {
				if bVal, ok := m.Array.GetAsBool(tkey); ok {
					return bVal
				}
			}
		}

		if hasDefaultValue {
			return defaultValue
		}

	}

	return false // zero value
}

func (m *DJSON) GetAsFloat(key ...interface{}) float64 {
	if len(key) == 0 {

		switch m.JsonType {
		case JSON_NULL, JSON_ARRAY, JSON_OBJECT:
			return 0
		case JSON_STRING:
			if fVal, err := strconv.ParseFloat(m.String, 64); err == nil {
				return fVal
			}
			return 0
		case JSON_BOOL:
			if m.Bool {
				return 1
			}
			return 0
		case JSON_INT:
			return float64(m.Int)
		case JSON_FLOAT:
			return m.Float
		}

	} else {

		var hasDefaultValue bool
		var defaultValue float64

		if len(key) >= 2 {
			defaultValue, hasDefaultValue = key[1].(float64)
		}

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {

				if m.Object.HasKey(tkey) {
					if fVal, ok := m.Object.GetAsFloat(tkey); ok {
						return fVal
					}
				}

			}
		case int:
			if m.JsonType == JSON_ARRAY {
				if m.Array.Size() > tkey {
					if fVal, ok := m.Array.GetAsFloat(tkey); ok {
						return fVal
					}
				}
			}
		}

		if hasDefaultValue {
			return defaultValue
		}
	}

	return 0 // zero value
}

func (m *DJSON) GetAsString(key ...interface{}) string {

	if len(key) == 0 {

		return m.ToString()

	} else {

		var hasDefaultValue bool
		var defaultValue string

		if len(key) >= 2 {
			defaultValue, hasDefaultValue = key[1].(string)
		}

		switch tkey := key[0].(type) {
		case string:
			if m.JsonType == JSON_OBJECT {
				if m.Object.HasKey(tkey) {
					return m.Object.GetAsString(tkey)
				} else {
					if hasDefaultValue {
						return defaultValue
					}
				}
			} else {
				return tkey // maybe default
			}
		case int:
			if m.JsonType == JSON_ARRAY {
				if m.Array.Size() > tkey {
					return m.Array.GetAsString(tkey)
				} else {
					if hasDefaultValue {
						return defaultValue
					}
				}
			} else {
				return strconv.Itoa(tkey) // maybe default
			}
		}

	}

	return "" // zero value
}

func (m *DJSON) ToString() string {

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

	return "" // zero value
}

func (m *DJSON) ReplaceAt(k interface{}, v interface{}) *DJSON {
	switch tkey := k.(type) {
	case string:
		if m.JsonType == JSON_OBJECT {

			if m.Object.HasKey(tkey) {
				m.Object.Put(tkey, v)
			}

		}
	case int:
		if m.JsonType == JSON_ARRAY {
			if m.Array.Size() > tkey {
				m.Array.ReplaceAt(tkey, v)
			}
		}
	}

	return m
}

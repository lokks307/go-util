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
	var err error
	XPathRegExp, err = regexp.Compile(`\[(\"[a-zA-Z0-9]+\"|[0-9]+)\]`)
	if err != nil {
		XPathRegExp = nil
	}
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
	if v == nil {
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_NULL
		return m
	}

	if IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		m.Int, _ = getIntBase(v)
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_INT
		return m
	}

	if IsInTypes(v, "float32", "float64") {
		m.Float, _ = getFloatBase(v)
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_FLOAT
		return m
	}

	if IsInTypes(v, "bool") {
		m.Bool, _ = getBoolBase(v)
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_BOOL
		return m
	}

	if IsInTypes(v, "string") {
		m.String, _ = getStringBase(v)
		m.Array = nil
		m.Object = nil
		m.JsonType = JSON_STRING
		return m
	}

	switch t := v.(type) {
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

func (m *DJSON) GetRaw(key ...interface{}) interface{} {
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

func (m *DJSON) HasKey(key interface{}) bool {
	switch tkey := key.(type) {
	case string:
		if m.JsonType == JSON_OBJECT {
			return m.Object.HasKey(tkey)
		}
	case int:
		if m.JsonType == JSON_ARRAY {
			return tkey >= 0 && m.Array.Size() > tkey
		}
	}

	return false
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

	if m.JsonType == JSON_STRING || m.JsonType == JSON_INT || m.JsonType == JSON_FLOAT {
		return nil, false
	}

	if len(key) == 0 {
		if m.JsonType == JSON_OBJECT {
			return m, true
		}
	} else {

		var ok bool
		var nObj *DO

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

// The DJSON as return shared Array.

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
		var nArr *DA

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

	return ""
}

func (m *DJSON) Size() int {
	return m.Length()
}

func (m *DJSON) Length() int {
	if m.JsonType == JSON_NULL {
		return 0
	}

	if m.JsonType == JSON_ARRAY {
		return m.Array.Length()
	}

	if m.JsonType == JSON_OBJECT {
		return m.Object.Length()
	}

	return 1
}

func (m *DJSON) getTypeSimple(key interface{}) string {

	switch m.JsonType {
	case JSON_OBJECT:
		if key, tok := key.(string); tok {
			if typeStr, ok := m.Object.GetType(key); ok {
				return typeStr
			}
		}
	case JSON_ARRAY:
		if idx, tok := key.(int); tok {
			if typeStr, ok := m.Array.GetType(idx); ok {
				return typeStr
			}
		}
	}

	return ""
}

func (m *DJSON) isSameType(key interface{}, inTypeStr string) bool {

	return m.getTypeSimple(key) == inTypeStr
}

func (m *DJSON) IsBool(key ...interface{}) bool {
	if len(key) == 0 {
		return m.JsonType == JSON_BOOL
	}

	return m.isSameType(key[0], "bool")
}

func (m *DJSON) IsInt(key ...interface{}) bool {
	if len(key) == 0 {
		return m.JsonType == JSON_INT
	}

	return m.isSameType(key[0], "int")
}

func (m *DJSON) IsNumeric(key ...interface{}) bool {
	if len(key) == 0 {
		return m.JsonType == JSON_FLOAT || m.JsonType == JSON_INT
	}

	return m.isSameType(key[0], "int") || m.isSameType(key[0], "float")
}

func (m *DJSON) IsFloat(key ...interface{}) bool {
	if len(key) == 0 {
		return m.JsonType == JSON_FLOAT
	}

	return m.isSameType(key[0], "float")
}

func (m *DJSON) IsString(key ...interface{}) bool {
	if len(key) == 0 {
		return m.JsonType == JSON_STRING
	}

	return m.isSameType(key[0], "string")
}

func (m *DJSON) IsNull(key ...interface{}) bool {
	if len(key) == 0 {
		return m.JsonType == JSON_NULL
	}

	return m.isSameType(key[0], "null")
}

func (m *DJSON) IsObject(key ...interface{}) bool {
	if len(key) == 0 {
		return m.JsonType == JSON_OBJECT
	}

	return m.isSameType(key[0], "object")
}

func (m *DJSON) IsArray(key ...interface{}) bool {
	if len(key) == 0 {
		return m.JsonType == JSON_ARRAY
	}

	return m.isSameType(key[0], "array")
}

func (m *DJSON) GetType(key ...interface{}) string {
	if len(key) == 0 {
		switch m.JsonType {
		case JSON_NULL:
			return "null"
		case JSON_OBJECT:
			return "object"
		case JSON_ARRAY:
			return "array"
		case JSON_STRING:
			return "string"
		case JSON_INT:
			return "int"
		case JSON_FLOAT:
			return "float"
		case JSON_BOOL:
			return "bool"
		}

		return ""
	}

	return m.getTypeSimple(key[0])
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

func (m *DJSON) GetTypePath(path string) string {
	var pathType string

	_ = m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			pathType, _ = da.GetType(idx)
		},
		func(do *DO, key string, v interface{}) {
			pathType, _ = do.GetType(key)
		},
	)

	return pathType
}

func (m *DJSON) RemovePath(path string) error {
	return m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			da.Remove(idx)
		},
		func(do *DO, key string, v interface{}) {
			do.Remove(key)
		},
	)
}

func (m *DJSON) UpdatePath(path string, val interface{}) error {
	return m.doPathFunc(path, val,
		func(da *DA, idx int, v interface{}) {
			da.ReplaceAt(idx, v)
		},
		func(do *DO, key string, v interface{}) {
			do.Put(key, v)
		},
	)
}

func (m *DJSON) doPathFunc(path string, val interface{},
	arrayTaskFunc func(da *DA, idx int, v interface{}),
	objectTaskFunc func(do *DO, key string, v interface{})) error {

	if XPathRegExp == nil {
		return unavailableError
	}

	matches := XPathRegExp.FindAllStringSubmatch(path, -1)

	pathLen := len(matches)

	if pathLen == 0 {
		return invalidPathError
	}

	jsonMode := m.JsonType
	dObject := m.Object
	dArray := m.Array

	for idx := range matches {

		kstr := matches[idx][1]

		kidx, err := strconv.Atoi(kstr)
		if err != nil {
			if strings.HasPrefix(kstr, `"`) && strings.HasSuffix(kstr, `"`) {
				kstr = strings.TrimRight(strings.TrimLeft(kstr, `"`), `"`)
			} else if strings.HasPrefix(kstr, `'`) && strings.HasSuffix(kstr, `'`) {
				kstr = strings.TrimRight(strings.TrimLeft(kstr, `'`), `'`)
			}

			if jsonMode != JSON_OBJECT {
				return invalidPathError
			}

			if dObject == nil {
				return invalidPathError
			}

			if idx == pathLen-1 {
				objectTaskFunc(dObject, kstr, val)
				return nil
			} else {
				if _, ok := dObject.Map[kstr]; !ok {
					return invalidPathError
				}

				switch t := dObject.Map[kstr].(type) {
				case *DO:
					dObject = t
					dArray = nil
					jsonMode = JSON_OBJECT
				case *DA:
					dObject = nil
					dArray = t
					jsonMode = JSON_ARRAY
				default:
					return invalidPathError
				}
			}

		} else {
			if jsonMode != JSON_ARRAY {
				return invalidPathError
			}

			if dArray == nil {
				return invalidPathError
			}

			for dArray.Size() <= kidx {
				dArray.PushBack(0)
			}

			if idx == pathLen-1 {
				arrayTaskFunc(dArray, kidx, val)
				return nil
			} else {
				switch t := dArray.Element[kidx].(type) {
				case *DO:
					dObject = t
					dArray = nil
					jsonMode = JSON_OBJECT
				case *DA:
					dObject = nil
					dArray = t
					jsonMode = JSON_ARRAY
				default:
					return invalidPathError
				}
			}
		}
	}

	return invalidPathError
}

func (m *DJSON) AutoFields(s interface{}) {
	target := reflect.ValueOf(s)
	elements := target.Elem()

	for i := 0; i < elements.NumField(); i++ {
		mValue := elements.Field(i)
		mType := elements.Type().Field(i)
		tag := mType.Tag.Get("json")

		if !mValue.CanSet() {
			continue
		}

		switch mType.Type.String() {
		case "int", "int8", "int16", "int32", "int64":
			mValue.SetInt(m.GetAsInt(tag))
		case "uint", "uint8", "uint16", "uint32", "uint64":
			mValue.SetUint(uint64(m.GetAsInt(tag)))
		case "float32", "float64":
			mValue.SetFloat(m.GetAsFloat(tag))
		case "string":
			mValue.SetString(m.GetAsString(tag))
		case "bool":
			mValue.SetBool(m.GetAsBool(tag))
		}
	}
}

package djson

import (
	"encoding/json"

	"github.com/volatiletech/null/v8"
)

type DO struct {
	Map map[string]interface{}
}

func NewObject() *DO {
	return &DO{
		Map: make(map[string]interface{}),
	}
}

func (m *DO) Put(key string, value interface{}) *DO {
	if IsBaseType(value) {
		m.Map[key] = value
		return m
	}

	switch t := value.(type) {
	case null.String:
		m.Map[key] = t.String
	case null.Bool:
		m.Map[key] = t.Bool
	case null.Int:
		m.Map[key] = t.Int
	case null.Int8:
		m.Map[key] = t.Int8
	case null.Int16:
		m.Map[key] = t.Int16
	case null.Int32:
		m.Map[key] = t.Int32
	case null.Int64:
		m.Map[key] = t.Int64
	case null.Uint:
		m.Map[key] = t.Uint
	case null.Uint8:
		m.Map[key] = t.Uint8
	case null.Uint16:
		m.Map[key] = t.Uint16
	case null.Uint32:
		m.Map[key] = t.Uint32
	case null.Uint64:
		m.Map[key] = t.Uint64
	case null.Float32:
		m.Map[key] = t.Float32
	case null.Float64:
		m.Map[key] = t.Float64
	case DO:
		m.Map[key] = &t
	case DA:
		m.Map[key] = &t
	case *DO:
		m.Map[key] = t
	case *DA:
		m.Map[key] = t
	case map[string]interface{}:
		m.Map[key] = ConverMapToObject(t)
	case []interface{}:
		m.Map[key] = ConvertSliceToArray(t)
	case Object:
		m.Map[key] = ConverMapToObject(t)
	case Array:
		m.Map[key] = ConvertSliceToArray(t)
	case DJSON:
		m.Map[key] = t.GetAsInterface()
	case *DJSON:
		m.Map[key] = t.GetAsInterface()
	case nil:
		m.Map[key] = nil
	}

	return m
}

func (m *DO) PutAsArray(key string, array ...interface{}) *DO {
	nArray := NewArray()
	nArray.Put(array)
	m.Put(key, nArray)
	return m
}

func (m *DO) Append(obj map[string]interface{}) *DO {
	for k, v := range obj {
		m.Put(k, v)
	}

	return m
}

func (m *DO) HasKey(key string) bool {
	_, ok := m.Map[key]
	return ok
}

func (m *DO) GetAsString(key string) string {
	value, ok := m.Map[key]
	if !ok {
		return ""
	}

	switch t := value.(type) {
	case DO:
		return t.ToString()
	case DA:
		return t.ToString()
	case *DO:
		return t.ToString()
	case *DA:
		return t.ToString()
	case nil:
		return "null"
	}

	if str, ok := getStringBase(m.Map[key]); ok {
		return str
	}

	return ""
}

func (m *DO) GetAsString2(key string) (string, bool) {
	value, ok := m.Map[key]
	if !ok {
		return "", false
	}

	switch t := value.(type) {
	case DO:
		return t.ToString(), true
	case DA:
		return t.ToString(), true
	case *DO:
		return t.ToString(), true
	case *DA:
		return t.ToString(), true
	case nil:
		return "null", true
	}

	return getStringBase(m.Map[key])
}

func (m *DO) Get(key string) (interface{}, bool) {
	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	return value, true
}

func (m *DO) GetType(key string) (string, bool) {
	value, ok := m.Map[key]
	if !ok {
		return "", false
	}

	switch value.(type) {
	case DA, *DA:
		return "array", true
	case DO, *DO:
		return "object", true
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return "int", true
	case float32, float64:
		return "float", true
	case string:
		return "string", true
	case bool:
		return "bool", true
	case nil:
		return "null", true
	}

	return "", false
}

func (m *DO) GetAsBool(key string) (bool, bool) {
	value, ok := m.Map[key]
	if !ok {
		return false, false
	}

	if boolVal, ok := getBoolBase(value); ok {
		return boolVal, true
	}

	return false, false
}

func (m *DO) GetAsFloat(key string) (float64, bool) {
	value, ok := m.Map[key]
	if !ok {
		return 0, false
	}

	if floatVal, ok := getFloatBase(value); ok {
		return floatVal, true
	}

	return 0, false
}

func (m *DO) GetAsInt(key string) (int64, bool) {
	value, ok := m.Map[key]
	if !ok {
		return 0, false
	}

	if intVal, ok := getIntBase(value); ok {
		return intVal, true
	}

	return 0, false
}

func (m *DO) GetAsObject(key string) (*DO, bool) {
	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case DO:
		return &t, true
	case *DO:
		return t, true
	case **DO:
		return *t, true
	}

	return nil, false
}

func (m *DO) GetAsArray(key string) (*DA, bool) {

	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case DA:
		return &t, true
	case *DA:
		return t, true
	case **DA:
		return *t, true
	}

	return nil, false

}

func (m *DO) Remove(keys ...string) *DO {
	for idx := range keys {
		delete(m.Map, keys[idx])
	}
	return m
}

func (m *DO) ToStringPretty() string {
	jsonByte, _ := json.MarshalIndent(ConverObjectToMap(m), "", "   ")
	return string(jsonByte)
}

func (m *DO) ToString() string {
	jsonByte, _ := json.Marshal(ConverObjectToMap(m))
	return string(jsonByte)
}

func (m *DO) Length() int {
	return len(m.Map)
}

func (m *DO) Size() int {
	return len(m.Map)
}

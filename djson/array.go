package djson

import (
	"encoding/json"
)

type DA struct {
	Element []interface{}
}

func NewArray() *DA {
	return &DA{
		Element: make([]interface{}, 0),
	}
}

func (m *DA) PushBack(values interface{}) *DA {
	return m.Insert(m.Size(), values)
}

func (m *DA) PushFront(values interface{}) *DA {
	return m.Insert(0, values)
}

func (m *DA) ReplaceAt(idx int, value interface{}) *DA {
	if idx >= m.Size() || idx < 0 {
		return m
	}

	if IsBaseType(value) {
		m.Element[idx] = value
		return m
	}

	switch t := value.(type) {
	case *DA:
		m.Element[idx] = t
	case *DO:
		m.Element[idx] = t
	case DA:
		m.Element[idx] = &t
	case DO:
		m.Element[idx] = &t
	case map[string]interface{}:
		m.Element[idx] = ConverMapToObject(t)
	case []interface{}:
		m.Element[idx] = ConvertSliceToArray(t)
	case Object:
		m.Element[idx] = ConverMapToObject(t)
	case Array:
		m.Element[idx] = ConvertSliceToArray(t)
	case DJSON:
		m.Element[idx] = t.GetAsInterface()
	case *DJSON:
		m.Element[idx] = t.GetAsInterface()
	case nil:
		m.Element[idx] = nil
	}

	return m
}

func (m *DA) Insert(idx int, value interface{}) *DA {
	if idx > m.Size() || idx < 0 {
		return m
	}

	if idx == m.Size() { // back
		m.Element = append(m.Element, nil)
	} else {
		m.Element = append(m.Element[:idx+1], m.Element[idx:]...)
	}

	return m.ReplaceAt(idx, value)
}

func (m *DA) Put(value interface{}) *DA {

	switch t := value.(type) {
	case Array:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []interface{}:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	default:
		m.Insert(m.Size(), value)
	}

	return m
}

func (m *DA) Size() int {
	return len(m.Element)
}

func (m *DA) Length() int {
	return len(m.Element)
}

func (m *DA) Remove(idx int) *DA {
	if idx >= m.Size() || idx < 0 {
		return m
	}

	m.Element = append(m.Element[:idx], m.Element[idx+1:]...)
	return m
}

func (m *DA) Get(idx int) (interface{}, bool) {
	if idx >= m.Size() || idx < 0 {
		return nil, false
	}

	return m.Element[idx], true
}

func (m *DA) GetType(idx int) (string, bool) {
	if idx >= m.Size() || idx < 0 {
		return "", false
	}

	switch m.Element[idx].(type) {
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

func (m *DA) GetAsBool(idx int) (bool, bool) {
	if idx >= m.Size() || idx < 0 {
		return false, false
	}

	if boolVal, ok := getBoolBase(m.Element[idx]); ok {
		return boolVal, true
	}

	return false, false
}

func (m *DA) GetAsFloat(idx int) (float64, bool) {
	if idx >= m.Size() || idx < 0 {
		return 0, false
	}

	if floatVal, ok := getFloatBase(m.Element[idx]); ok {
		return floatVal, true
	}

	return 0, false
}

func (m *DA) GetAsInt(idx int) (int64, bool) {
	if idx >= m.Size() || idx < 0 {
		return 0, false
	}

	if intVal, ok := getIntBase(m.Element[idx]); ok {
		return intVal, true
	}

	return 0, false
}

func (m *DA) GetAsObject(idx int) (*DO, bool) {
	if idx >= m.Size() || idx < 0 {
		return nil, false
	}

	switch t := m.Element[idx].(type) {
	case DO:
		return &t, true
	case *DO:
		return t, true
	}

	return nil, false
}

func (m *DA) GetAsArray(idx int) (*DA, bool) {
	if idx >= m.Size() || idx < 0 {
		return nil, false
	}

	switch t := m.Element[idx].(type) {
	case DA:
		return &t, true
	case *DA:
		return t, true
	}

	return nil, false
}

func (m *DA) GetAsString(idx int) string {
	if idx >= m.Size() || idx < 0 {
		return ""
	}

	switch t := m.Element[idx].(type) {
	case DA:
		return t.ToString()
	case *DA:
		return t.ToString()
	case DO:
		return t.ToString()
	case *DO:
		return t.ToString()
	case nil:
		return "null"
	}

	if str, ok := getStringBase(m.Element[idx]); ok {
		return str
	}

	return ""
}

func (m *DA) GetAsString2(idx int) (string, bool) {
	if idx >= m.Size() || idx < 0 {
		return "", false
	}

	switch t := m.Element[idx].(type) {
	case DA:
		return t.ToString(), true
	case *DA:
		return t.ToString(), true
	case DO:
		return t.ToString(), true
	case *DO:
		return t.ToString(), true
	case nil:
		return "null", true
	}

	return getStringBase(m.Element[idx])
}

func (m *DA) ToStringPretty() string {
	jsonByte, _ := json.MarshalIndent(ConvertArrayToSlice(m), "", "   ")
	return string(jsonByte)
}

func (m *DA) ToString() string {
	jsonByte, _ := json.Marshal(ConvertArrayToSlice(m))
	return string(jsonByte)
}

package djson

import "encoding/json"

type A struct {
	Element []interface{}
}

func NewArray() *A {
	return &A{
		Element: make([]interface{}, 0),
	}
}

func (m *A) PushBack(values interface{}) *A {
	return m.Insert(m.Size(), values)
}

func (m *A) PushFront(values interface{}) *A {
	return m.Insert(0, values)
}

func (m *A) Insert(idx int, value interface{}) *A {
	if idx > m.Size() || idx < 0 {
		return m
	}

	if idx == m.Size() { // back
		m.Element = append(m.Element, nil)
	} else {
		m.Element = append(m.Element[:idx+1], m.Element[idx:]...)
	}

	if IsBaseType(value) {
		m.Element[idx] = value
		return m
	}

	switch tValue := value.(type) {
	case *A:
		m.Element[idx] = tValue
	case *O:
		m.Element[idx] = tValue
	case A:
		m.Element[idx] = &tValue
	case O:
		m.Element[idx] = &tValue
	case map[string]interface{}:
		m.Element[idx] = ConverMapToObject(tValue)
	case []interface{}:
		m.Element[idx] = ConvertSliceToArray(tValue)
	case _Object:
		m.Element[idx] = ConverMapToObject(tValue)
	case _Array:
		m.Element[idx] = ConvertSliceToArray(tValue)
	case nil:
		m.Element[idx] = nil
	}

	return m
}

func (m *A) Put(values ...interface{}) *A {

	for idx := range values {
		m.Insert(m.Size(), values[idx])
	}

	return m
}

func (m *A) Size() int {
	return len(m.Element)
}

func (m *A) Length() int {
	return len(m.Element)
}

func (m *A) Remove(idx int) *A {
	if idx >= m.Size() || idx < 0 {
		return m
	}

	m.Element = append(m.Element[:idx], m.Element[idx+1:]...)
	return m
}

func (m *A) Get(idx int) (interface{}, bool) {
	if idx >= m.Size() || idx < 0 {
		return nil, false
	}

	return m.Element[idx], true
}

func (m *A) GetAsBool(idx int) bool {
	if idx >= m.Size() || idx < 0 {
		return false
	}

	if boolVal, ok := getBoolBase(m.Element[idx]); ok {
		return boolVal
	}

	return false
}

func (m *A) GetAsFloat(idx int) float64 {
	if idx >= m.Size() || idx < 0 {
		return 0
	}

	if floatVal, ok := getFloatBase(m.Element[idx]); ok {
		return floatVal
	}

	return 0
}

func (m *A) GetAsInt(idx int) int64 {
	if idx >= m.Size() || idx < 0 {
		return 0
	}

	if intVal, ok := getIntBase(m.Element[idx]); ok {
		return intVal
	}

	return 0
}

func (m *A) GetAsObject(idx int) (*O, bool) {
	if idx >= m.Size() || idx < 0 {
		return nil, false
	}

	switch t := m.Element[idx].(type) {
	case O:
		return &t, true
	case *O:
		return t, true
	}

	return nil, false
}

func (m *A) GetAsArray(idx int) (*A, bool) {
	if idx >= m.Size() || idx < 0 {
		return nil, false
	}

	switch t := m.Element[idx].(type) {
	case A:
		return &t, true
	case *A:
		return t, true
	}

	return nil, false
}

func (m *A) GetAsString(idx int) string {
	if idx >= m.Size() || idx < 0 {
		return ""
	}

	str, ok := getStringBase(m.Element[idx])
	if !ok {
		switch t := m.Element[idx].(type) {
		case A:
			return t.ToString()
		case *A:
			return t.ToString()
		case O:
			return t.ToString()
		case *O:
			return t.ToString()
		}
	}

	return str
}

func (m *A) ToStringPretty() string {
	jsonByte, _ := json.MarshalIndent(ConvertArrayToSlice(m), "", "   ")
	return string(jsonByte)
}

func (m *A) ToString() string {
	jsonByte, _ := json.Marshal(ConvertArrayToSlice(m))
	return string(jsonByte)
}

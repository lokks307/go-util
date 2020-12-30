package djson

import "encoding/json"

type O struct {
	Map map[string]interface{}
}

func NewObject() *O {
	return &O{
		Map: make(map[string]interface{}),
	}
}

func (m *O) Put(key string, value interface{}) *O {
	if IsBaseType(value) {
		m.Map[key] = value
		return m
	}

	switch t := value.(type) {
	case O:
		m.Map[key] = &t
	case A:
		m.Map[key] = &t
	case *O:
		m.Map[key] = t
	case *A:
		m.Map[key] = t
	case map[string]interface{}:
		m.Map[key] = ConverMapToObject(t)
	case []interface{}:
		m.Map[key] = ConvertSliceToArray(t)
	case _Object:
		m.Map[key] = ConverMapToObject(t)
	case _Array:
		m.Map[key] = ConvertSliceToArray(t)
	case nil:
		m.Map[key] = t
	}

	return m
}

func (m *O) PutAsArray(key string, array ...interface{}) *O {
	nArray := NewArray()
	nArray.Put(array)
	m.Map[key] = nArray
	return m
}

func (m *O) Append(obj map[string]interface{}) *O {
	for k, v := range obj {
		m.Put(k, v)
	}

	return m
}

func (m *O) GetAsString(key string) string {
	value, ok := m.Map[key]
	if !ok {
		return ""
	}

	switch t := value.(type) {
	case O:
		return t.ToString()
	case A:
		return t.ToString()
	case *O:
		return t.ToString()
	case *A:
		return t.ToString()
	case nil:
		return "null"
	}

	str, ok := getStringBase(m.Map[key])
	if !ok {
		return ""
	}

	return str
}

func (m *O) Get(key string) (interface{}, bool) {
	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	return value, true
}

func (m *O) GetAsBool(key string) bool {
	value, ok := m.Map[key]
	if !ok {
		return false
	}

	if boolVal, ok := getBoolBase(value); ok {
		return boolVal
	}

	return false
}

func (m *O) GetAsFloat(key string) float64 {
	value, ok := m.Map[key]
	if !ok {
		return 0
	}

	if floatVal, ok := getFloatBase(value); ok {
		return floatVal
	}

	return 0
}

func (m *O) GetAsInt(key string) int64 {
	value, ok := m.Map[key]
	if !ok {
		return 0
	}

	if intVal, ok := getIntBase(value); ok {
		return intVal
	}

	return 0
}

func (m *O) GetAsObject(key string) (*O, bool) {
	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case O:
		return &t, true
	case *O:
		return t, true
	case **O:
		return *t, true
	}

	return nil, false
}

func (m *O) GetAsArray(key string) (*A, bool) {

	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case A:
		return &t, true
	case *A:
		return t, true
	case **A:
		return *t, true
	}

	return nil, false

}

func (m *O) Remove(keys ...string) *O {
	for idx := range keys {
		delete(m.Map, keys[idx])
	}
	return m
}

func (m *O) ToStringPretty() string {
	jsonByte, _ := json.MarshalIndent(ConverObjectToMap(m), "", "   ")
	return string(jsonByte)
}

func (m *O) ToString() string {
	jsonByte, _ := json.Marshal(ConverObjectToMap(m))
	return string(jsonByte)
}

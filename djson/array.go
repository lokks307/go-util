package djson

import (
	"encoding/json"
	"reflect"
	"sort"

	"github.com/volatiletech/null/v8"
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

	if n, ok := value.(json.Number); ok {
		if i, err := n.Int64(); err == nil {
			m.Element[idx] = i
			return m
		}
		if f, err := n.Float64(); err == nil {
			m.Element[idx] = f
			return m
		}
	}

	switch t := value.(type) {
	case null.String:
		m.Element[idx] = t.String
	case null.Bool:
		m.Element[idx] = t.Bool
	case null.Int:
		m.Element[idx] = t.Int
	case null.Int8:
		m.Element[idx] = t.Int8
	case null.Int16:
		m.Element[idx] = t.Int16
	case null.Int32:
		m.Element[idx] = t.Int32
	case null.Int64:
		m.Element[idx] = t.Int64
	case null.Uint:
		m.Element[idx] = t.Uint
	case null.Uint8:
		m.Element[idx] = t.Uint8
	case null.Uint16:
		m.Element[idx] = t.Uint16
	case null.Uint32:
		m.Element[idx] = t.Uint32
	case null.Uint64:
		m.Element[idx] = t.Uint64
	case null.Float32:
		m.Element[idx] = t.Float32
	case null.Float64:
		m.Element[idx] = t.Float64
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
		idx = m.Size()
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
	case *DA:
		for idx := range t.Element {
			m.Insert(m.Size(), t.Element[idx])
		}
	case Array:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []interface{}:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []int:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []int8:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []int16:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []int32:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []int64:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []uint:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []uint8:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []uint16:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []uint32:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []uint64:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []float32:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []float64:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []bool:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []string:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.String:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Bool:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Float32:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Float64:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Int:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Int8:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Int16:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Int32:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Int64:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Uint:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Uint8:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Uint16:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Uint32:
		for idx := range t {
			m.Insert(m.Size(), t[idx])
		}
	case []null.Uint64:
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

func (m *DA) SortObject(isAsc bool, key string) bool {
	numElement := len(m.Element)

	if numElement == 0 {
		return false
	}

	var elemType string

	for i := range m.Element {
		do, ok := m.GetAsObject(i)
		if !ok {
			return false
		}

		kv, ok := do.Get(key)
		if !ok {
			return false
		}

		eachType := reflect.TypeOf(kv).String()
		if elemType == "" {
			elemType = eachType
		} else {
			if elemType != eachType {
				return false
			}
		}
	}

	if isAsc {

		sort.Slice(m.Element, func(i, j int) bool {

			ido, _ := m.Element[i].(*DO)
			jdo, _ := m.Element[j].(*DO)

			switch elemType {
			case "string":

				iRune := []rune(ido.GetAsString(key))
				jRune := []rune(jdo.GetAsString(key))

				lenToInspect := len(iRune)
				if len(jRune) < lenToInspect {
					lenToInspect = len(jRune)
				}

				for k := 0; k < lenToInspect; k++ {
					if iRune[k] == jRune[k] {
						continue
					}

					if iRune[k] < jRune[k] {
						return true
					}

					return false
				}

				return len(iRune) < len(jRune)

			case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
				iInt, _ := ido.GetAsInt(key)
				jInt, _ := jdo.GetAsInt(key)
				return iInt < jInt
			case "float32", "float64":
				iFloat, _ := ido.GetAsFloat(key)
				jFloat, _ := jdo.GetAsFloat(key)
				return iFloat < jFloat
			case "bool":
				jBool, _ := jdo.GetAsBool(key)
				return jBool
			default:
				return true
			}
		})

	} else {
		sort.Slice(m.Element, func(i, j int) bool {

			ido, _ := m.Element[i].(*DO)
			jdo, _ := m.Element[j].(*DO)

			switch elemType {
			case "string":

				iRune := []rune(ido.GetAsString(key))
				jRune := []rune(jdo.GetAsString(key))

				lenToInspect := len(iRune)
				if len(jRune) < lenToInspect {
					lenToInspect = len(jRune)
				}

				for k := 0; k < lenToInspect; k++ {
					if iRune[k] == jRune[k] {
						continue
					}

					if iRune[k] > jRune[k] {
						return true
					}

					return false
				}

				return len(iRune) > len(jRune)

			case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
				iInt, _ := ido.GetAsInt(key)
				jInt, _ := jdo.GetAsInt(key)
				return iInt > jInt
			case "float32", "float64":
				iFloat, _ := ido.GetAsFloat(key)
				jFloat, _ := jdo.GetAsFloat(key)
				return iFloat > jFloat
			case "bool":
				iBool, _ := ido.GetAsBool(key)
				return iBool
			default:
				return false
			}
		})
	}

	return true

}

func (m *DA) Sort(isAsc bool) bool {

	numElement := len(m.Element)

	if numElement == 0 {
		return false
	}

	var elemType string
	for i := range m.Element {
		eachType := reflect.TypeOf(m.Element[i]).String()
		if elemType == "" {
			elemType = eachType
		} else {
			if elemType != eachType {
				return false
			}
		}
	}

	if elemType == "string" {

		tmpElment := make([]string, numElement)
		for i := range m.Element {
			tmpElment[i] = m.Element[i].(string)
		}

		sort.Strings(tmpElment)

		if isAsc {
			for i := range m.Element {
				m.Element[i] = tmpElment[i]
			}
		} else {
			for i := range m.Element {
				m.Element[numElement-i-1] = tmpElment[i]
			}
		}

		return true

	} else {

		if isAsc {
			sort.Slice(m.Element, func(i, j int) bool {
				switch elemType {
				case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
					iInt, _ := m.GetAsInt(i)
					jInt, _ := m.GetAsInt(j)
					return iInt < jInt
				case "float32", "float64":
					iFloat, _ := m.GetAsFloat(i)
					jFloat, _ := m.GetAsFloat(j)
					return iFloat < jFloat
				case "bool":
					jBool, _ := m.GetAsBool(j)
					return jBool
				default:
					return true
				}
			})
		} else {
			sort.Slice(m.Element, func(i, j int) bool {
				switch elemType {
				case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
					iInt, _ := m.GetAsInt(i)
					jInt, _ := m.GetAsInt(j)
					return iInt > jInt
				case "float32", "float64":
					iFloat, _ := m.GetAsFloat(i)
					jFloat, _ := m.GetAsFloat(j)
					return iFloat < jFloat
				case "bool":
					iBool, _ := m.GetAsBool(i)
					return iBool
				default:
					return false
				}
			})
		}

		return true
	}
}

func (m *DA) SortAsc() bool {
	return m.Sort(true)
}

func (m *DA) SortDesc() bool {
	return m.Sort(false)
}

func (m *DA) Equal(t *DA) bool {
	if m.Size() != t.Size() {
		return false
	}

	for i := range m.Element {

		mtype := reflect.TypeOf(m.Element[i]).String()
		ttype := reflect.TypeOf(t.Element[i]).String()

		if mtype != ttype {
			return false
		}

		switch mtype {
		case "string":
			if m.Element[i].(string) != t.Element[i].(string) {
				return false
			}
		case "bool":
			if m.Element[i].(bool) != t.Element[i].(bool) {
				return false
			}
		case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
			mInt, _ := m.GetAsInt(i)
			tInt, _ := t.GetAsInt(i)
			if mInt != tInt {
				return false
			}
		case "float32", "float64":
			mFloat, _ := m.GetAsFloat(i)
			tFloat, _ := t.GetAsFloat(i)
			if mFloat != tFloat {
				return false
			}
		case "*djson.DO":
			mdo := m.Element[i].(*DO)
			tdo := t.Element[i].(*DO)

			if !mdo.Equal(tdo) {
				return false
			}
		case "*djson.DA":
			mda := m.Element[i].(*DA)
			tda := t.Element[i].(*DA)

			if !mda.Equal(tda) {
				return false
			}
		case "*djson.DJSON":
			mjson := m.Element[i].(*DJSON)
			tjson := t.Element[i].(*DJSON)

			if !mjson.Equal(tjson) {
				return false
			}
		}
	}

	return true
}

func (m *DA) Clone() *DA {

	t := NewArray()

	t.Element = make([]interface{}, m.Size())

	for i := range m.Element {

		mtype := reflect.TypeOf(m.Element[i]).String()

		switch mtype {
		case "string":
			t.Element[i] = m.Element[i].(string)
		case "bool":
			t.Element[i] = m.Element[i].(bool)
		case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
			t.Element[i], _ = m.GetAsInt(i)
		case "float32", "float64":
			t.Element[i], _ = m.GetAsFloat(i)
		case "*djson.DO":
			mdo := m.Element[i].(*DO)
			t.Element[i] = mdo.Clone()
		case "*djson.DA":
			mda := m.Element[i].(*DA)
			t.Element[i] = mda.Clone()
		case "*djson.DJSON":
			mdjson := m.Element[i].(*DJSON)
			t.Element[i] = mdjson.Clone()
		}
	}

	return t
}

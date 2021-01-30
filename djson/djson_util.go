package djson

import (
	"reflect"
)

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

func (m *DJSON) toFieldsValue(val reflect.Value) {

	for i := 0; i < val.NumField(); i++ {
		eachVal := val.Field(i)
		eachType := val.Type().Field(i)
		eachTag := eachType.Tag.Get("json")

		if !eachVal.CanSet() || !m.HasKey(eachTag) {
			continue
		}

		eachKind := eachType.Type.Kind()

		if eachKind == reflect.Struct {

			switch eachType.Type.String() {
			case "null.String":
				eachVal.FieldByName("String").SetString(m.GetAsString(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Bool":
				eachVal.FieldByName("Bool").SetBool(m.GetAsBool(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Float32":
				eachVal.FieldByName("Float32").SetFloat(m.GetAsFloat(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Float64":
				eachVal.FieldByName("Float64").SetFloat(m.GetAsFloat(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int":
				eachVal.FieldByName("Int").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int8":
				eachVal.FieldByName("Int8").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int16":
				eachVal.FieldByName("Int16").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int32":
				eachVal.FieldByName("Int32").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int64":
				eachVal.FieldByName("Int64").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint":
				eachVal.FieldByName("Uint").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint8":
				eachVal.FieldByName("Uint8").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint16":
				eachVal.FieldByName("Uint16").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint32":
				eachVal.FieldByName("Uint32").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint64":
				eachVal.FieldByName("Uint64").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			default:

				if oJson, ok := m.GetAsObject(eachTag); ok {
					oJson.toFieldsValue(eachVal)
				}

			}

		} else {

			switch eachType.Type.String() {
			case "int", "int8", "int16", "int32", "int64":
				eachVal.SetInt(m.GetAsInt(eachTag))
			case "uint", "uint8", "uint16", "uint32", "uint64":
				eachVal.SetUint(uint64(m.GetAsInt(eachTag)))
			case "float32", "float64":
				eachVal.SetFloat(m.GetAsFloat(eachTag))
			case "string":
				eachVal.SetString(m.GetAsString(eachTag))
			case "bool":
				eachVal.SetBool(m.GetAsBool(eachTag))
			}
		}
	}
}

func (m *DJSON) ToFields(st interface{}) {
	target := reflect.ValueOf(st)
	elements := target.Elem()
	m.toFieldsValue(elements)
}

func (m *DJSON) fromFiledsValue(val reflect.Value) {

	for i := 0; i < val.NumField(); i++ {
		eachVal := val.Field(i)
		eachType := val.Type().Field(i)
		eachTag := eachType.Tag.Get("json")

		eachKind := eachType.Type.Kind()

		if eachKind == reflect.Struct {

			switch eachType.Type.String() {
			case "null.String":
				m.Put(eachTag, eachVal.FieldByName("String").String())
			case "null.Bool":
				m.Put(eachTag, eachVal.FieldByName("Bool").Bool())
			case "null.Float32":
				m.Put(eachTag, eachVal.FieldByName("Float32").Float())
			case "null.Float64":
				m.Put(eachTag, eachVal.FieldByName("Float64").Float())
			case "null.Int":
				m.Put(eachTag, eachVal.FieldByName("Int").Int())
			case "null.Int8":
				m.Put(eachTag, eachVal.FieldByName("Int8").Int())
			case "null.Int16":
				m.Put(eachTag, eachVal.FieldByName("Int16").Int())
			case "null.Int32":
				m.Put(eachTag, eachVal.FieldByName("Int32").Int())
			case "null.Int64":
				m.Put(eachTag, eachVal.FieldByName("Int64").Int())
			case "null.Uint":
				m.Put(eachTag, eachVal.FieldByName("Uint").Uint())
			case "null.Uint8":
				m.Put(eachTag, eachVal.FieldByName("Uint8").Uint())
			case "null.Uint16":
				m.Put(eachTag, eachVal.FieldByName("Uint16").Uint())
			case "null.Uint32":
				m.Put(eachTag, eachVal.FieldByName("Uint32").Uint())
			case "null.Uint64":
				m.Put(eachTag, eachVal.FieldByName("Uint64").Uint())
			default:
				sJson := NewDJSON()
				sJson.fromFiledsValue(eachVal)
				m.Put(eachTag, sJson)
			}

		} else {
			switch eachType.Type.String() {
			case "int", "int8", "int16", "int32", "int64":
				m.Put(eachTag, eachVal.Int())
			case "uint", "uint8", "uint16", "uint32", "uint64":
				m.Put(eachTag, eachVal.Uint())
			case "float32", "float64":
				m.Put(eachTag, eachVal.Float())
			case "string":
				m.Put(eachTag, eachVal.String())
			case "bool":
				m.Put(eachTag, eachVal.Bool())
			}
		}
	}
}

func (m *DJSON) FromFields(st interface{}) {
	baseValue := reflect.ValueOf(st)

	kind := baseValue.Type().Kind()

	if kind == reflect.Array || kind == reflect.Slice {

		m.JsonType = JSON_ARRAY
		m.fromFiledsValue(baseValue)

	} else if kind == reflect.Struct {

		m.JsonType = JSON_OBJECT
		m.fromFiledsValue(baseValue)

	} else if kind == reflect.Map {

		m.JsonType = JSON_OBJECT

		for _, e := range baseValue.MapKeys() {
			eachKey, ok := e.Interface().(string)
			if !ok {
				continue
			}

			eachVal := baseValue.MapIndex(e)

			switch eachVal.Type().String() {
			case "int", "int8", "int16", "int32", "int64":
				m.Put(eachKey, eachVal.Int())
			case "uint", "uint8", "uint16", "uint32", "uint64":
				m.Put(eachKey, eachVal.Uint())
			case "float32", "float64":
				m.Put(eachKey, eachVal.Float())
			case "string":
				m.Put(eachKey, eachVal.String())
			case "bool":
				m.Put(eachKey, eachVal.Bool())
			case "nil":
				m.Put(eachKey, nil)
			}
		}
	}
}

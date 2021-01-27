package djson

import "reflect"

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

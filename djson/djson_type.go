package djson

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
	if IsEmptyArg(key) {
		return m.JsonType == JSON_BOOL
	}

	return m.isSameType(key[0], "bool")
}

func (m *DJSON) IsInt(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m.JsonType == JSON_INT
	}

	return m.isSameType(key[0], "int")
}

func (m *DJSON) IsNumeric(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m.JsonType == JSON_FLOAT || m.JsonType == JSON_INT
	}

	return m.isSameType(key[0], "int") || m.isSameType(key[0], "float")
}

func (m *DJSON) IsFloat(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m.JsonType == JSON_FLOAT
	}

	return m.isSameType(key[0], "float")
}

func (m *DJSON) IsString(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m.JsonType == JSON_STRING
	}

	return m.isSameType(key[0], "string")
}

func (m *DJSON) IsNull(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m.JsonType == JSON_NULL
	}

	return m.isSameType(key[0], "null")
}

func (m *DJSON) IsObject(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m.JsonType == JSON_OBJECT
	}

	return m.isSameType(key[0], "object")
}

func (m *DJSON) IsArray(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m.JsonType == JSON_ARRAY
	}

	return m.isSameType(key[0], "array")
}

func (m *DJSON) GetType(key ...interface{}) string {
	if IsEmptyArg(key) {
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

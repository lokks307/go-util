package djson

func (m *DJSON) GetAsObjectPath(path string) (*DJSON, bool) {

	retJson := NewDJSON()

	err := m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if obj, ok := da.GetAsObject(idx); ok {
				retJson.Object = obj
				retJson.JsonType = JSON_OBJECT
			}
		},
		func(do *DO, key string, v interface{}) {
			if obj, ok := do.GetAsObject(key); ok {
				retJson.Object = obj
				retJson.JsonType = JSON_OBJECT
			}
		},
	)

	if err != nil || retJson.JsonType != JSON_OBJECT {
		return nil, false
	}

	return retJson, true
}

func (m *DJSON) GetAsArrayPath(path string) (*DJSON, bool) {

	retJson := NewDJSON()

	err := m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if arr, ok := da.GetAsArray(idx); ok {
				retJson.Array = arr
				retJson.JsonType = JSON_ARRAY
			}
		},
		func(do *DO, key string, v interface{}) {
			if arr, ok := do.GetAsArray(key); ok {
				retJson.Array = arr
				retJson.JsonType = JSON_ARRAY
			}
		},
	)

	if err != nil || retJson.JsonType != JSON_ARRAY {
		return nil, false
	}

	return retJson, true
}

func (m *DJSON) GetAsFloatPath(path string, defFloat ...float64) float64 {
	var retFloat float64
	var ok bool

	err := m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			retFloat, ok = da.GetAsFloat(idx)
		},
		func(do *DO, key string, v interface{}) {
			retFloat, ok = do.GetAsFloat(key)
		},
	)

	if err == nil && ok {
		return retFloat
	}

	if len(defFloat) > 0 {
		return defFloat[0]
	}

	return 0
}

func (m *DJSON) GetAsIntPath(path string, defInt ...int64) int64 {
	var retInt int64
	var ok bool

	err := m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			retInt, ok = da.GetAsInt(idx)
		},
		func(do *DO, key string, v interface{}) {
			retInt, ok = do.GetAsInt(key)
		},
	)

	if err == nil && ok {
		return retInt
	}

	if len(defInt) > 0 {
		return defInt[0]
	}

	return 0
}

func (m *DJSON) GetAsBoolPath(path string, defBool ...bool) bool {

	var retBool bool
	var ok bool

	err := m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			retBool, ok = da.GetAsBool(idx)
		},
		func(do *DO, key string, v interface{}) {
			retBool, ok = do.GetAsBool(key)
		},
	)

	if err == nil && ok {
		return retBool
	}

	if len(defBool) > 0 {
		return defBool[0]
	}

	return false
}

func (m *DJSON) GetAsStringPath(path string) string {
	var retStr string

	_ = m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			retStr = da.GetAsString(idx)
		},
		func(do *DO, key string, v interface{}) {
			retStr = do.GetAsString(key)
		},
	)

	return retStr
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

func (m *DJSON) PutPath(path string, okey interface{}, oval interface{}) error {
	return m.doPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if oIdx, ok := okey.(int); ok {
				da.Insert(oIdx, oval)
			}
		},
		func(do *DO, key string, v interface{}) {
			if oKey, ok := okey.(string); ok {
				do.Put(oKey, oval)
			}
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

func (m *DJSON) doPathFuncCore(
	arrayTaskFunc func(da *DA, idx int, v interface{}),
	objectTaskFunc func(do *DO, key string, v interface{}),
	val interface{}, token ...interface{}) error {

	jsonMode := m.JsonType
	dObject := m.Object
	dArray := m.Array

	tokenLen := len(token)

	for idx := range token {
		switch tkey := token[idx].(type) {
		case string:
			if jsonMode != JSON_OBJECT || dObject == nil {
				return invalidPathError
			}

			if idx == tokenLen-1 {
				objectTaskFunc(dObject, tkey, val)
				return nil
			} else {
				if _, ok := dObject.Map[tkey]; !ok {
					return invalidPathError
				}

				switch t := dObject.Map[tkey].(type) {
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
		case int:
			if jsonMode != JSON_ARRAY || dArray == nil {
				return invalidPathError
			}

			for dArray.Size() <= tkey {
				dArray.PushBack(0)
			}

			if idx == tokenLen-1 {
				arrayTaskFunc(dArray, tkey, val)
				return nil
			} else {
				switch t := dArray.Element[tkey].(type) {
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
		default:
			return invalidPathError
		}
	}

	return invalidPathError

}

func (m *DJSON) doPathFunc(path string, val interface{},
	arrayTaskFunc func(da *DA, idx int, v interface{}),
	objectTaskFunc func(do *DO, key string, v interface{})) error {
	return m.doPathFuncCore(arrayTaskFunc, objectTaskFunc, val, PathTokenizer(path)...)
}

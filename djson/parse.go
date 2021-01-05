package djson

import (
	"encoding/json"
	"errors"
)

func ParseToObject(doc string) (*DO, error) {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(doc), &data)
	if err != nil {
		return nil, errors.New("not Object")
	}

	return ParseObject(data), nil

}

func ParseToArray(doc string) (*DA, error) {
	var data []interface{}

	err := json.Unmarshal([]byte(doc), &data)
	if err != nil {
		return nil, errors.New("not Array")
	}

	return ParseArray(data), nil
}

func ParseObject(data map[string]interface{}) *DO {
	obj := NewObject()
	for k, v := range data {
		if IsBaseType(v) {
			//log.Println("ParseObject > BaseType")
			obj.Put(k, v)
			continue
		}

		switch tValue := v.(type) {
		case []interface{}: // Array
			//log.Println("ParseObject > []interface{}")
			nArr := ParseArray(tValue)
			obj.Put(k, nArr)
		case map[string]interface{}: // Object
			//log.Println("ParseObject > map[string]interface{}")
			nObj := ParseObject(tValue)
			obj.Put(k, nObj)
		case nil: // null
			obj.Put(k, nil)
		}
	}

	return obj
}

func ParseArray(data []interface{}) *DA {
	arr := NewArray()

	for idx := range data {
		if IsBaseType(data[idx]) {
			//log.Println("ParseArray > BaseType")
			arr.Put(data[idx])
			continue
		}

		switch tValue := data[idx].(type) {
		case []interface{}: // Array
			//log.Println("ParseArray > []interface{}")
			nArr := ParseArray(tValue)
			arr.Put(nArr)
		case map[string]interface{}: // Object
			//log.Println("ParseArray > map[string]interface{}")
			nObj := ParseObject(tValue)
			arr.Put(nObj)
		case nil: // null
			arr.Put(nil)
		}
	}

	return arr
}

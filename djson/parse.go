package djson

import (
	"encoding/json"
	"fmt"
)

func ParseToObject(doc string) O {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(doc), &data)
	if err != nil {
		fmt.Println(err)
		return Object()
	}

	return ParseObject(data)

}

func ParseObject(data map[string]interface{}) O {
	obj := Object()
	for k, v := range data {
		if IsBaseType(data[k]) {
			obj.Put(k, v)
			continue
		}

		switch t := data[k].(type) {
		case []interface{}: // Array
			arr := ParseArray(t)
			obj.Put(k, arr)
		case map[string]interface{}: // Object
			nObj := ParseObject(t)
			obj.Put(k, nObj)
		case nil: // null
			obj.Put(k, nil)
		}
	}

	return obj
}

func ParseArray(data []interface{}) *A {
	arr := Array()

	for idx := range data {
		if IsBaseType(data[idx]) {
			arr.Put(data[idx])
			continue
		}

		switch t := data[idx].(type) {
		case []interface{}: // Array
			nArr := ParseArray(t)
			arr.Put(nArr)
		case map[string]interface{}: // Object
			nObj := ParseObject(t)
			arr.Put(nObj)
		case nil: // null
			arr.Put(nil)
		}
	}

	return arr
}

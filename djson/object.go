package djson

import (
	"errors"
	"fmt"
	"reflect"
)

type O map[string]interface{}

func Object() O {
	return make(map[string]interface{})
}

func (this O) Put(key string, value interface{}) O {
	this[key] = value
	return this
}

func (this O) PutAsArray(key string, array ...interface{}) O {
	newArray := Array()
	newArray.Put(array)

	this[key] = newArray
	return this
}

func (this O) Append(obj map[string]interface{}) O {
	for k, v := range obj {
		this[k] = v
	}

	return this
}

func (this O) GetAsString(key string) string {
	switch t := this[key].(type) {
	case map[string]interface{}:
		return _string(t)
	case []interface{}:
		return "Array"
	case O:
		return t.String()
	case A:
		return t.String()
	case nil:
		return "null"
	}

	str, ok := _stringBase(this[key])
	if !ok {
		return "Object"
	}

	return str
}

func (this O) Get(key string) interface{} {
	return this[key]
}

func (this O) GetObject(key string) (value O, err error) {
	switch t := this[key].(type) {
	case map[string]interface{}:
		return t, nil
	case O:
		return t, nil
	}

	return nil, errors.New(fmt.Sprintf("Casting error. Interface is %s, not jsongo.object", reflect.TypeOf(this[key])))
}

func (this O) GetArray(key string) (newArray *A, err error) {
	newArray = Array()

	if IsBaseType(this[key]) {
		return nil, errors.New("Casting error. Interface is base type, not array")
	}

	switch arr := this[key].(type) {
	case []interface{}:
		newArray.Put(arr)
		return newArray, nil
	case []string:
		newArray.Put(arr)
		return newArray, nil
	case []bool:
		newArray.Put(arr)
		return newArray, nil
	case []int:
		newArray.Put(arr)
		return newArray, nil
	case []uint:
		newArray.Put(arr)
		return newArray, nil
	case []int8:
		newArray.Put(arr)
		return newArray, nil
	case []uint8:
		newArray.Put(arr)
		return newArray, nil
	case []int16:
		newArray.Put(arr)
		return newArray, nil
	case []uint16:
		newArray.Put(arr)
		return newArray, nil
	case []int32:
		newArray.Put(arr)
		return newArray, nil
	case []uint32:
		newArray.Put(arr)
		return newArray, nil
	case []int64:
		newArray.Put(arr)
		return newArray, nil
	case []uint64:
		newArray.Put(arr)
		return newArray, nil
	case []float32:
		newArray.Put(arr)
		return newArray, nil
	case []float64:
		newArray.Put(arr)
		return newArray, nil
	case A:
		return &arr, nil
	}

	return nil, errors.New(fmt.Sprintf("Casting error. Interface is %s, not jsongo.A or []interface{}", reflect.TypeOf(this[key])))
}

func (this O) Remove(keys ...string) O {
	for idx := range keys {
		delete(this, keys[idx])
	}
	return this
}

func (this O) Indent() string {
	return indent(this)
}

func (this O) String() string {
	return _string(this)
}

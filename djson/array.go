package djson

import (
	"errors"
	"fmt"
	"reflect"
)

type A []interface{}

func Array() *A {
	return &A{make([]interface{}, 0)}
}

func (this *A) PushBack(values interface{}) *A {
	return this.Put(values)
}

func (this *A) PushFront(values interface{}) *A {
	return this.Insert(0, values)
}

func (this *A) Put(values ...interface{}) *A {
	fmt.Println("Array.Put() - ", reflect.TypeOf(values[0]).String(), values[0])
	*this = append((*this), values...)
	fmt.Println((*this))
	return this
}

func (this *A) Indent() string {
	return indent((*this))
}

func (this *A) String() string {
	return _string((*this))
}

func (this *A) Size() int {
	return len((*this))
}

func (this *A) Length() int {
	return len((*this))
}

func (this *A) Remove(idx int) *A {
	if idx >= (*this).Size() || idx < 0 {
		return this
	}

	(*this) = append((*this)[:idx], (*this)[idx+1:]...)
	return this
}

func (this *A) Insert(idx int, value interface{}) *A {
	if idx > (*this).Size() || idx < 0 {
		return this
	}

	(*this) = append((*this)[:idx+1], (*this)[idx:]...)
	(*this)[idx] = value
	return this
}

func (this *A) Get(idx int) interface{} {
	if idx >= (*this).Size() || idx < 0 {
		return nil
	}

	return (*this)[idx]
}

func (this *A) GetAsObject(idx int) O {
	if idx >= (*this).Size() || idx < 0 {
		return nil
	}

	switch t := (*this)[idx].(type) {
	case O:
		return t
	case map[string]interface{}:
		return O(t)
	}

	return nil
}

func (this *A) GetAsString(idx int) string {

	fmt.Println("GetAsString(", idx, ")")

	if idx >= (*this).Size() || idx < 0 {
		fmt.Println((*this))
		return ""
	}

	str, ok := _stringBase((*this)[idx])
	if !ok {
		return "Object"
	}

	return str
}

func (this *A) OfString() (values []string, err error) {
	for _, value := range *this {
		if reflect.TypeOf(value).String() != "string" {
			return nil, errors.New(fmt.Sprintf("Value is %s, not a string.", reflect.TypeOf(value)))
		}

		values = append(values, value.(string))
	}

	return values, nil
}

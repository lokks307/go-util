package djson

import (
	"fmt"
	"testing"
)

func TestParseJson(t *testing.T) {
	doc := `{
		"name": "Anderson",
		"age": 50,
		"address" : [
			"Seoul", "Korea"
		],
		"family" : {
			"father": "John",
			"mother": "Jane"
		}
	}`

	obj := ParseToObject(doc)
	if obj != nil {
		fmt.Println(obj.GetAsString("name"))
		arr, err := obj.GetArray("address")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(arr.GetAsString(0))
	}
}

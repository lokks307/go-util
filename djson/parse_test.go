package djson

import (
	"log"
	"testing"
)

func TestParseJson(t *testing.T) {
	jsonDoc := `{
		"name": null,
		"age": 50,
		"address" : [
			"Seoul", "Korea"
		],
		"family" : {
			"father": "John",
			"mother": "Jane"
		}
	}`

	obj, err := ParseToObject(jsonDoc)
	if err == nil {
		log.Println(obj.GetAsString("name"))
		arr, ok := obj.GetAsArray("address")
		if !ok {
			log.Fatal("no such key")
		}
		log.Println(arr.GetAsString(1))
	}
}

func TestParseArray(t *testing.T) {

	jsonDoc := `[
		{
			"name":"Ricardo Longa",
			"idade":28,
			"skills":[
				"Golang","Android"
			]
		},
		{
			"name":"Hery Victor",
			"idade":32,
			"skills":[
				"Golang",
				"Java"
			]
		}
	]`

	arr, err := ParseToArray(jsonDoc)
	if err == nil {
		obj, ok := arr.GetAsObject(0)
		if !ok {
			log.Fatal("no such index")
		}
		log.Println(obj.GetAsString("name"))
	} else {
		log.Fatal("not array")
	}

}

func TestParseDJSON(t *testing.T) {
	jsonDoc := `[
		{
			"name":"Ricardo Longa",
			"idade":28,
			"skills":[
				"Golang","Android"
			]
		},
		{
			"name":"Hery Victor",
			"idade":32,
			"skills":[
				"Golang",
				"Java"
			]
		}
	]`

	aJson := NewDJSON().Parse(jsonDoc)

	bJson, ok := aJson.GetAsObject(1)
	if !ok {
		log.Fatal("not object")
	}

	log.Println(bJson.GetAsInt("skills"))
	log.Println(bJson.GetAsString())
}

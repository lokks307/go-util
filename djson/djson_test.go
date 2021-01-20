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

func TestPutDJSON(t *testing.T) {
	aJson := NewDJSON().Put(
		Array{
			Object{
				"name":  "Ricardo Longa",
				"idade": 28,
				"skills": Array{
					"Golang",
					"Android",
				},
			},
			Object{
				"name":  "Hery Victor",
				"idade": 32,
				"skills": Array{
					"Golang",
					"Java",
				},
			},
		},
	)

	bJson, ok := aJson.GetAsObject(1)
	if !ok {
		log.Fatal("not object")
	}

	log.Println(bJson.GetAsInt("skills"))
	log.Println(bJson.GetAsString())

	log.Println(bJson.HasKey("name"))
}

func TestPutDArrayDJSON(t *testing.T) {
	aJson := NewDJSON()
	aJson.Put(
		Array{
			Array{
				1, 2, 3, 4,
			},
			Array{
				5, 6, 7, 8,
			},
		},
	)

	bJson, ok := aJson.Get(1)
	if !ok {
		log.Fatal("not array")
	}

	log.Println(bJson.GetAsInt(1))
	log.Println(bJson.GetAsString())

	log.Println(bJson.HasKey(1))
}

func TestArrayAppendDJSON(t *testing.T) {
	aJson := NewDJSON()
	aJson.Put(
		Array{
			1, 2, 3, 4,
		},
	)
	aJson.Put( // append array
		Array{
			5, 6, 7, 8,
		},
	)

	log.Println(aJson.GetAsString())
}

func TestObjectAppendDJSON(t *testing.T) {
	aJson := NewDJSON()
	aJson.Put(
		Object{
			"name": "Hery Victor",
		},
	)
	aJson.Put( // append
		Object{
			"idade": 32,
		},
	)

	aJson.Put("not appended")

	log.Println(aJson.GetAsString())
}

package djson

import (
	"log"
	"testing"
)

func TestParseJson(t *testing.T) {
	jsonDoc := `{
		"name": null,
		"age": 9223372036854775807,
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
		log.Println(obj.GetAsString("age"))
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

func TestUpdateDJSON(t *testing.T) {
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
	bJson.Put(Object{"hobbies": Array{"game"}})
	_ = bJson.UpdatePath(`["hobbies"][1]`, "running")
	_ = bJson.UpdatePath(`["hobbies"][0]`, "art")

	log.Println(aJson.GetAsString())
}

func TestHandleDJSON(t *testing.T) {
	jsonDoc := `{
		"name":"Ricardo Longa",
		"idade":28,
		"skills":[
			"Golang","Android"
		]
	}`

	mJson := NewDJSON().Parse(jsonDoc)

	aJson := NewDJSON()
	// aJson.Put("name", mJson.GetAsInterface("name"))
	// aJson.Put("idade", mJson.GetAsInterface("idade"))

	aJson.PutAsObject("name", mJson.GetAsInterface("name"))
	aJson.PutAsObject("idade", mJson.GetAsInterface("idade"))

	log.Println(aJson.ToString())

}

func TestFastDeclare(t *testing.T) {
	dJson := NewObjectJSON(
		"key", "value", "key2", "value2",
	)

	log.Println(dJson.ToString())

	aJson := NewArrayJSON(
		1, 2, 3, 4, 5, 6, 7,
	)

	log.Println(aJson.ToString())

	tJson := NewDJSON(JSON_ARRAY)
	tJson.Put(1)

	log.Println(tJson.ToString())

	pJson := NewDJSON()
	pJson.Put(1)

	log.Println(pJson.ToString())
}

func TestWrapArray(t *testing.T) {
	mJson := NewDJSON(JSON_ARRAY)

	bb := []int{1, 2, 3, 4, 5, 6, 7, 8}

	mJson.Put(bb)

	log.Println(mJson.ToString())
}

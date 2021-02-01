package djson

import (
	"log"
	"testing"
)

func TestGetType(t *testing.T) {

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

	mJson := NewDJSON().Parse(jsonDoc)

	log.Println(mJson.GetType())
	log.Println(mJson.GetType(1))
	log.Println(mJson.IsObject(1))
	log.Println(mJson.GetTypePath(`[1]["name"]`))
	log.Println(mJson.GetTypePath(`[1]["idade"]`))
	log.Println(mJson.GetTypePath(`[1]["skills"]`))

	log.Println(mJson.GetAsFloat(1, 0.7))
}

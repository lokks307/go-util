package djson

import (
	"log"
	"regexp"
	"testing"
)

func TestParsePath(t *testing.T) {

	XPathRegExp, err := regexp.Compile(`\[(\"[a-zA-Z0-9]+\"|[0-9]+)\]`)

	if err != nil {
		log.Fatal(err)
	}

	paths := XPathRegExp.FindAllStringSubmatch(`["name"][1]`, -1)
	log.Println(paths)
}

func TestPutPath(t *testing.T) {
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

	err := aJson.UpdatePath(`[1]["name"]`, Object{
		"first":  "kim",
		"family": "kim",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(aJson.GetAsString())

	err = aJson.UpdatePath(`[1]["name"]["first"]`, "seo")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(aJson.GetAsString())

	err = aJson.RemovePath(`[1]["name"]["family"]`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(aJson.GetAsString())

	err = aJson.RemovePath(`[1]["name"]`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(aJson.GetAsString())

	err = aJson.RemovePath(`[1]`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(aJson.GetAsString())
}

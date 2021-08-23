package djson

import (
	"log"
	"testing"
)

func TestValidator1(t *testing.T) {

	dv := NewValidator()
	dv.Compile(`{
		"type": "OBJECT",
		"object": {
			"name": {
				"type": "STRING",
				"min": 4,
				"max": 25
			},
			"skills": {
				"type": "ARRAY",
				"array": [
					{
						"type": "HEX",
						"min": 6,
						"max": 6
					}
				]
			}
		}
	}`)

	asjon := NewDJSON().Parse(`{"name": "wakeupbb", "skill": ["040809","aaaaaa"]}`)

	if dv.IsValid(asjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

}

func TestValidator2(t *testing.T) {

	dv := NewValidator()
	dv.Compile(`[{
		"type": "OBJECT",
		"object": {
			"name": {
				"type": "STRING",
				"min": 4,
				"max": 25
			},
			"skill": {
				"required": true,
				"type": "ARRAY",
				"array": [
					{
						"type": "HEX",
						"min": 6,
						"max": 6
					}
				]
			}
		}
	},"HEX"]`)

	asjon := NewDJSON().Parse(`FF112345`)

	if dv.IsValid(asjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

}

func TestValidator3(t *testing.T) {

	dv := NewValidator()
	dv.Compile(`{
		"type": "OBJECT",
		"object": {
			"name": {
				"type": "OBJECT",
				"object": {
					"age": ["INT","HEX"],
					"home": "STRING"
				}
			}
		}
	}`)

	asjon := NewDJSON().Parse(`{"name":{"age":"11y","home":"aaa"}}`)

	if dv.IsValid(asjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

}

func TestValidator4(t *testing.T) {

	dv := NewValidator()
	dv.Compile(`{
		"type": "OBJECT",
		"object": {
			"name": "HEX256.IF.EXIST",
			"skill": "HEX256.IF.EXIST",
			"home": "BOOL"
		}
	}`)

	asjon := NewDJSON().Parse(`{"name":"f0e42af202cb3fd1c45b973b3597951cd4a205b9d620217e84b09a6fd0cc96bf","skill":"","home":true}`)

	if dv.IsValid(asjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

}

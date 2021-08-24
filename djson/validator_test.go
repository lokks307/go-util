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
			"skill": {
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

	asjon := NewDJSON().Parse(`{"name": "wakeupbb", "skill": ["040809","aaaaab"]}`)

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
	dv.Compile(`
	{
		"type": "OBJECT",
		"object": {
			"service_name": "STRING",
			"service_id": "HEX256.IF.EXIST",
			"call_address": "HEX256.IF.EXIST",
			"medi": "HEX256.IF.EXIST",
			"medi2": "HEX256.IF.EXIST",
			"patient_id": "NONEMPTY.STRING",
			"patient_tel": "TELEPHONE",
			"patient_name": "NONEMPTY.STRING",
			"patient_color": "STRING",
			"patient_note": "STRING",
			"patient_sms": "BOOL",
			"patient_suptel": {
				"type": "ARRAY",
				"array": {
					"type": "OBJECT",
					"object": {
						"name": {
							"type":"STRING",
							"required":true
						},
						"tel": {
							"type":"TELEPHONE",
							"required":true
						}
					}
				}
			}
		}
	}`)

	asjon := NewDJSON().Parse(`{
		"service_name": "",
		"service_id": "",
		"call_address": "",
		"medi": "",
		"medi2": "",
		"patient_id": "99999",
		"patient_tel": "010-1111-1111",
		"patient_name": "삼색이",
		"patient_color": "",
		"patient_note": "",
		"patient_sms": true,
		"patient_suptel": [{"name":"father","tel":"010-6666-6666"}]
	}`)

	if dv.IsValid(asjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

}

func TestValidator5(t *testing.T) {

	dv := NewValidator()
	dv.Compile(`
	{
		"type": "OBJECT",
		"object": {}
	}`)

	asjon := NewDJSON().Parse(`{"name":"top"}`)
	bsjon := NewDJSON().Parse(`[]`)

	if dv.IsValid(asjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

	if dv.IsValid(bsjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

}

func TestValidator6(t *testing.T) {

	dv := NewValidator()
	dv.Compile(`
	[{
		"type": "OBJECT",
		"object": {}
	},"HEX"]`)

	asjon := NewDJSON().Parse(`{"name":"top"}`)
	bsjon := NewDJSON().Parse(`BB112233`)

	if dv.IsValid(asjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

	if dv.IsValid(bsjon) {
		log.Println("valid")
	} else {
		log.Println("not valid")
	}

}

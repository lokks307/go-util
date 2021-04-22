package djson

import (
	"log"
	"testing"

	"github.com/volatiletech/null/v8"
)

func TestToFieldTag(t *testing.T) {
	type User struct {
		Id    string      `json:"id"`
		Name  string      `json:"name"`
		Email null.String `json:"email"`
	}

	var user User

	mJson := NewDJSON().Put(
		Object{
			"id":    "id-1234",
			"name":  "Ricardo Longa",
			"email": "longa@test.com",
		},
	)

	mJson.ToFields(&user, "id", "email")

	log.Println(user)
}

func TestFromFieldTag(t *testing.T) {
	type User struct {
		Id    string      `json:"id"`
		Name  string      `json:"name"`
		Email null.String `json:"email"`
	}

	var user = User{
		Id:   "id-1234",
		Name: "Ricardo Longa",
		Email: null.String{
			String: "longa@test.com",
			Valid:  true,
		},
	}

	mJson := NewDJSON()
	mJson.FromFields(user)

	log.Println(mJson.ToString())
}

func TestFromFieldMapTest(t *testing.T) {

	type Name struct {
		First  string `json:"first"`
		Family string `json:"family"`
	}

	user := make(map[string]interface{})

	user["id"] = "id-1234"
	user["name"] = Name{
		First:  "Ricardo",
		Family: "Longa",
	}

	user["email"] = null.String{
		String: "longa@test.com",
		Valid:  true,
	}

	mJson := NewDJSON()
	mJson.FromFields(user, "name.first", "email")

	log.Println(mJson.ToString())
}

func TestSortingArray(t *testing.T) {
	mJson := NewArrayJSON(5, 6, 7, 8, 1, 2, 3, 4)
	if ok := mJson.SortAsc(); !ok {
		log.Fatal("sorting test failed")
	}

	log.Println(mJson.ToString())

	tJson := NewDJSON().Put(Object{
		"d": "aaa",
		"a": Array{
			5, 6, 7, 8, 1, 2, 3, 4,
		},
	})

	if ok := tJson.SortAsc("a"); !ok {
		log.Fatal("sorting test failed")
	}

	log.Println(tJson.ToString())

	err := tJson.SortDescPath(`["a"]`)

	if err != nil {
		log.Fatal("sorting path failed")
	}

	log.Println(tJson.ToString())

	pJson := NewDJSON().Put(
		Array{
			Object{
				"k": "1",
				"v": "1",
			},
			Object{
				"k": "22",
				"v": "2",
			},
			Object{
				"k": "4444",
				"v": "4",
			},
			Object{
				"k": "333",
				"v": "3",
			},
		},
	)

	pJson.SortObjectArrayAsc("k")

	p2Json := NewDJSON().Put(
		Object{
			"kk": Array{
				Object{
					"k": null.String{
						String: "1",
					},
					"v": "1",
				},
				Object{
					"k": null.String{
						String: "22",
					},
					"v": "2",
				},
				Object{
					"k": "9",
					"v": "4",
				},
				Object{
					"k": "1",
					"v": "3",
				},
			},
		},
	)

	p2Json.SortObjectArrayDescPath(`[kk]`, "k")

	log.Println(pJson.ToString())
	log.Println(p2Json.ToString())

}

func TestReflectType(t *testing.T) {

	aJson := NewDJSON().Put(Object{
		"name": "yu",
		"skill": Array{
			"running", "playing",
		},
	})

	bJson := NewDJSON().Put(Object{
		"skill": Array{
			"running", "playing",
		},
		"name": "yu",
	})

	if aJson.Equal(bJson) {
		log.Println("jsons are the same")
	} else {
		log.Println("jsons are not the same")
	}
}

func TestClone(t *testing.T) {
	aJson := NewDJSON().Put(Object{
		"name": "yu",
		"skill": Array{
			"running", "playing",
		},
	})

	bJson := aJson.Clone()

	aJson.UpdatePath(`[skill][2]`, "swimming")

	bJson.PutAsObject("name", "not you")
	bJson.UpdatePath(`[skill][2]`, "studying")

	log.Println(aJson.ToString())
	log.Println(bJson.ToString())
}

func TestFind(t *testing.T) {
	aJson := NewDJSON().Put(Array{
		Object{
			"name":  "1",
			"skill": "apple",
		},
		Object{
			"name":  "2",
			"skill": "banana",
		},
	})

	log.Println(aJson.ToString())

	bJson := aJson.Find("name", "1")
	if bJson == nil {
		log.Fatalln("find failed")
	}

	log.Println(bJson.ToString())
}

func TestAppend(t *testing.T) {
	aJson := NewDJSON().Put(Array{
		Object{
			"name":  "1",
			"skill": "apple",
		},
		Object{
			"name":  "2",
			"skill": "banana",
		},
	})

	bJson := NewDJSON().Put(Array{
		Object{
			"name":  "3",
			"skill": "apple",
		},
		Object{
			"name":  "4",
			"skill": "banana",
		},
	})

	log.Println(aJson.Append(bJson).ToString())
}

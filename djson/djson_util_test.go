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

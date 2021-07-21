package djson

import (
	"log"
	"testing"
)

func TestTokenizer(t *testing.T) {

	log.Println(PathTokenizer(`["aa"][1][b_b]`))  // [aa 1 b_b]
	log.Println(PathTokenizer(`["a'a"][1][b]b]`)) // [a'a 1 b]
}

func TestParse(t *testing.T) {
	doc := `[[1,2,3]]`
	tdjson := NewDJSON().Parse(doc)
	log.Println(tdjson.GetAsIntPath(`[0][0]`))
}

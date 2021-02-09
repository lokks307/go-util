package djson

import (
	"log"
	"testing"
)

func TestTokenizer(t *testing.T) {

	log.Println(PathTokenizer(`["aa"][1][b_b]`))  // [aa 1 b_b]
	log.Println(PathTokenizer(`["a'a"][1][b]b]`)) // [a'a 1 b]
}

package kakao

import "testing"

func TestLocalAPIParserGenerator(t *testing.T) {
	_, err := LocalAPIParserGenerator("starbucks")()
	if err != nil {
		t.Fatal(err)
	}
}

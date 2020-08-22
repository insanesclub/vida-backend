package kakao

import (
	"fmt"
	"log"
	"testing"
)

func TestGenerator(t *testing.T) {
	parser := LocalAPIParserGenerator("starbucks")
	result, err := parser()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
	for _, url := range result.Places() {
		err = MapParserGenerator(url)()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

package kakao

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maengsanha/vida-backend/model/kakao"
)

// LocalAPIParserGenerator generates a kakao local API parser.
func LocalAPIParserGenerator(query string) func() (kakao.LocalAPIResult, error) {
	return func() (result kakao.LocalAPIResult, _ error) {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://dapi.kakao.com/v2/local/search/keyword.json?query=%s", query), nil)
		if err != nil {
			return result, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("KakaoAK %s", kakao.REST_API_KEY))
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			return result, err
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&result)
		return result, err
	}
}

// MapParserGenerator generates a kakao map parser.
func MapParserGenerator(url string) func() error {
	return func() error {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		fmt.Println(url)
		fmt.Println(resp.Body)
		defer resp.Body.Close()
		return err
	}
}

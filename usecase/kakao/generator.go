package kakao

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benbjohnson/phantomjs"
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

func CrawlWithPhantom(url string) error {
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		return err
	}
	defer phantomjs.DefaultProcess.Close()
	page, err := phantomjs.CreateWebPage()
	if err != nil {
		return err
	}
	if err = page.Open(url); err != nil {
		return err
	}

	return nil
}

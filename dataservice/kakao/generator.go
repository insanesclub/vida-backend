package kakao

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/maengsanha/vida-backend/model/kakao/local"
)

const max_death_cnt = 3

// MapParserGenerator generates a kakao map local API parser.
func MapParserGenerator(query string) func() (local.Response, error) {
	var (
		is_end bool
		page   = 1
	)
	return func() (response local.Response, _ error) {
		if is_end {
			return response, fmt.Errorf("is_end: %t", is_end)
		}
		var death_cnt int

		for death_cnt < max_death_cnt { // retry up to maximum death count
			req, err := http.NewRequest("GET", fmt.Sprintf("https://dapi.kakao.com/v2/local/search/keyword.json?query=%s&page=%d", url.QueryEscape(query), page), nil)
			if err != nil {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			req.Header.Set("Authorization", fmt.Sprintf("KakaoAK %s", REST_API_KEY))
			client := new(http.Client)
			resp, err := client.Do(req)
			if err != nil {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			defer resp.Body.Close()
			if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			defer func() {
				is_end = response.Meta.IsEnd
				page++
			}()
			return response, nil
		}
		// aborted when the maximum death count is exceeded
		return response, fmt.Errorf("death count exceeded %d", death_cnt)
	}
}

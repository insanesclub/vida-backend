package kakao

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maengsanha/vida-backend/model/kakao"
)

// LocalAPIParserGenerator generates a kakao local API parser.
func LocalAPIParserGenerator(query string) func() (kakao.MapInfo, error) {
	var (
		is_end bool
		page   = 1
	)
	return func() (info kakao.MapInfo, _ error) {
		if is_end {
			return info, fmt.Errorf("is_end: %t", is_end)
		}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://dapi.kakao.com/v2/local/search/keyword.json?query=%s&page=%d", query, page), nil)
		if err != nil {
			return info, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("KakaoAK %s", kakao.REST_API_KEY))
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			return info, err
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&info)
		defer func() {
			is_end = info.Meta.IsEnd
			page++
		}()
		return info, err
	}
}

package instagram

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maengsanha/vida-backend/model/instagram"
)

const max_death_cnt = 3

// PageParserGenerator generates an Instagram page parser.
func PageParserGenerator(tag string) func() (instagram.Tagpage, error) {
	var (
		has_next_page = true
		end_cursor    string
	)

	return func() (page instagram.Tagpage, err error) {
		if !has_next_page {
			err = fmt.Errorf("has_next_page: %t", has_next_page)
			return
		}
		var death_cnt int

		for death_cnt < max_death_cnt { // retry up to maximum death count
			var resp = new(http.Response)
			resp, err = http.Get(fmt.Sprintf("http://www.instagram.com/explore/tags/%s?__a=1&max_id=%s", tag, end_cursor))
			if err != nil {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			defer resp.Body.Close()

			if err = json.NewDecoder(resp.Body).Decode(&page); err != nil {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			has_next_page = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.HasNextPage
			end_cursor = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.EndCursor
			return
		}
		// aborted when the maximum death count is exceeded
		err = fmt.Errorf("death count exceeded %d", death_cnt)
		return
	}
}

// PostParserGenerator generates an Instagram post parser.
func PostParserGenerator(shortcode string) func() (instagram.Post, error) {
	return func() (post instagram.Post, err error) {
		var death_cnt int

		for death_cnt < max_death_cnt { // retry up to maximum death count
			var resp = new(http.Response)
			resp, err = http.Get(fmt.Sprintf("http://www.instagram.com/p/%s?__a=1", shortcode))
			if err != nil {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			if resp.StatusCode == http.StatusNotFound { // give up in the event of PAGE NOT FOUND
				err = fmt.Errorf("failed with status code %d", resp.StatusCode)
				return
			}
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			defer resp.Body.Close()

			if err = json.NewDecoder(resp.Body).Decode(&post); err != nil {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			return
		}
		// aborted when the maximum death count is exceeded
		err = fmt.Errorf("death count exceeded %d", death_cnt)
		return
	}
}

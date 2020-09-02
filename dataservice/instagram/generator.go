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

	return func() (page instagram.Tagpage, _ error) {
		if !has_next_page {
			return page, fmt.Errorf("has_next_page: %t", has_next_page)
		}
		var death_cnt int

		for death_cnt < max_death_cnt { // retry up to maximum death count
			var resp = new(http.Response)
			resp, err := http.Get(fmt.Sprintf("https://www.instagram.com/explore/tags/%s/?__a=1&max_id=%s", tag, end_cursor))
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
			defer func() {
				has_next_page = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.HasNextPage
				end_cursor = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.EndCursor
			}()
			return page, nil
		}
		// aborted when the maximum death count is exceeded
		return page, fmt.Errorf("death count exceeded %d", death_cnt)
	}
}

// PostParserGenerator generates an Instagram post parser.
func PostParserGenerator(shortcode string) func() (instagram.Post, error) {
	return func() (post instagram.Post, _ error) {
		var death_cnt int

		for death_cnt < max_death_cnt { // retry up to maximum death count
			resp, err := http.Get(fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1", shortcode))
			if err != nil {
				death_cnt++
				continue // RETRY INSTRUCTION
			}
			if resp.StatusCode == http.StatusNotFound { // give up in the event of PAGE NOT FOUND
				return post, fmt.Errorf("failed with status code %d", resp.StatusCode)
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
			return post, nil
		}
		// aborted when the maximum death count is exceeded
		return post, fmt.Errorf("death count exceeded %d", death_cnt)
	}
}

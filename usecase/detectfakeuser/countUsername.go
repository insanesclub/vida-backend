package detectfakeuser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/haxana/vida-backend/model/instagram"
)

const MAX_DEATH_CNT = 3

// Usernames counts user names from the search results on Instagram with a given query.
func Usernames(query string) map[string]int {
	usernames := make(chan string)
	var syncer sync.WaitGroup
	syncer.Add(1)
	go parsePage(query, usernames, &syncer)
	go func() {
		syncer.Wait()
		close(usernames)
	}()

	var u_map sync.Map

	for username := range usernames {
		syncer.Add(1)
		go func(username string) {
			defer syncer.Done()
			if cnt, loaded := u_map.LoadOrStore(username, 1); loaded {
				u_map.Store(username, cnt.(int)+1)
			}
		}(username)
	}
	syncer.Wait()

	usrnames := make(map[string]int)

	u_map.Range(func(name, cnt interface{}) bool {
		usrnames[name.(string)] = cnt.(int)
		return true
	})

	return usrnames
}

func parsePage(query string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()
	var (
		has_next_page = true
		end_cursor    string
	)
	for has_next_page {
		url := fmt.Sprintf("https://www.instagram.com/explore/tags/%s/?__a=1&max_id=%s", query, end_cursor)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
			// error handling
		}
		defer resp.Body.Close()

		var page instagram.Tagpage
		if err = json.NewDecoder(resp.Body).Decode(&page); err != nil {
			// error handling
		}
		has_next_page = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.HasNextPage
		end_cursor = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.EndCursor

		for _, shortcode := range page.Shortcodes() {
			syncer.Add(1)
			go parseUsername(shortcode, ch, syncer)
		}
	}
}

func parseUsername(shortcode string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()

	url := fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1", shortcode)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
		// error handling
	}
	defer resp.Body.Close()

	var post instagram.Post
	if err = json.NewDecoder(resp.Body).Decode(&post); err != nil {
		// error handling
	}
	ch <- post.GraphQL.ShortcodeMedia.Owner.Username
}

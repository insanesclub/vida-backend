package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/haxana/vida-backend/model/instagram"
)

const (
	_ int32 = iota
	done
	_
	max_death_cnt
)

// Uploads counts user names from the search results on Instagram with a given query.
func Uploads(tag string) map[string]int {
	var usernames = make(chan string)
	var syncer sync.WaitGroup
	syncer.Add(1)
	go parsePage(tag, usernames, &syncer)
	go func() {
		syncer.Wait()
		close(usernames)
	}()
	var uploads = make(map[string]int)

	for username := range usernames {
		uploads[username]++
	}
	return uploads
}

func parsePage(tag string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()
	var (
		flag       int32
		end_cursor string
	)
	for atomic.LoadInt32(&flag) != done {
		var death_cnt int32
		for death_cnt < max_death_cnt { // retry up to maximum death count
			resp, err := http.Get(fmt.Sprintf("https://www.instagram.com/explore/tags/%s/?__a=1&max_id=%s", tag, end_cursor))
			if err != nil {
				death_cnt++
			} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
				death_cnt++
			} else {
				defer resp.Body.Close()
				var page = new(instagram.Tagpage)
				if err := json.NewDecoder(resp.Body).Decode(page); err != nil {
					death_cnt++
				} else {
					for _, shortcode := range page.Shortcodes() {
						syncer.Add(1)
						go parseUsername(shortcode, &flag, ch, syncer)
					}
					time.Sleep(0) // force context switching
					end_cursor = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.EndCursor
					if !page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.HasNextPage {
						return
					}
					break
				}
			}
		}
		if death_cnt > max_death_cnt { // aborted when the maximum death count is exceeded
			break
		}
	}
}

func parseUsername(shortcode string, flag *int32, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()
	var death_cnt int32
	for death_cnt < max_death_cnt {
		resp, err := http.Get(fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1", shortcode))
		if err != nil {
			death_cnt++
		} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
			death_cnt++
		} else {
			defer resp.Body.Close()
			var post = new(instagram.Post)
			if err := json.NewDecoder(resp.Body).Decode(post); err != nil {
				death_cnt++
			} else if !post.PostedToday() {
				atomic.StoreInt32(flag, done)
				return
			} else {
				ch <- post.GraphQL.ShortcodeMedia.Owner.Username
				return
			}
		}
	}
}

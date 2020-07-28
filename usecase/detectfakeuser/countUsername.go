package detectfakeuser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/haxana/vida-backend/model/instagram"
)

const max_death_cnt = 3

type atomicBool struct {
	flag int32
}

func (b *atomicBool) set(value bool) {
	var i int32
	if value {
		i = 1
	}
	atomic.StoreInt32(&b.flag, i)
}

func (b *atomicBool) get() bool {
	return atomic.LoadInt32(&b.flag) == 1
}

// Usernames counts user names from the search results on Instagram with a given query.
func Usernames(tag string) map[string]int {
	var usernames = make(chan string)
	var syncer sync.WaitGroup
	syncer.Add(1)
	go parsePage(tag, usernames, &syncer)
	go func() {
		syncer.Wait()
		close(usernames)
	}()
	var user_map = make(map[string]int)

	for username := range usernames {
		user_map[username]++
	}
	return user_map
}

func parsePage(tag string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()
	var (
		has_next_page = &atomicBool{flag: 1}
		end_cursor    string
	)
	for has_next_page.get() {
		var death_cnt int
		for death_cnt < max_death_cnt {
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
					has_next_page.set(page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.HasNextPage)
					end_cursor = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.EndCursor
					for _, shortcode := range page.Shortcodes() {
						syncer.Add(1)
						go parseUsername(shortcode, has_next_page, ch, syncer)
					}
					break
				}
			}
		}
		if death_cnt >= max_death_cnt {
			break
		}
	}
}

func parseUsername(shortcode string, has_next_page *atomicBool, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()
	var death_cnt int
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
				has_next_page.set(false)
				return
			} else {
				ch <- post.GraphQL.ShortcodeMedia.Owner.Username
				return
			}
		}
	}
}

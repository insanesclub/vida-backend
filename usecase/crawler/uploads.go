package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/haxana/vida-backend/model/instagram"
)

const max_death_cnt = 3

// Uploads counts user names from the search results on Instagram with a given tag.
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
		has_next_page = true
		end_cursor    string
	)
	for has_next_page {
		var death_cnt int
		for death_cnt < max_death_cnt { // retry up to maximum death count
			/*
				logging for test
			*/
			var url = fmt.Sprintf("https://www.instagram.com/explore/tags/%s/?__a=1&max_id=%s", tag, end_cursor)
			logging(fmt.Sprintf("requesting to %s", url))

			// resp, err := http.Get(fmt.Sprintf("https://www.instagram.com/explore/tags/%s/?__a=1&max_id=%s", tag, end_cursor))
			var request, err = http.NewRequest("GET", url, nil)
			if err != nil {
				/*
					logging for test
				*/
				logging(fmt.Sprintf("http.NewRequest: %v", err))

				death_cnt++
				continue
			}
			request.Header.Set("Accept", "application/json")
			var client = new(http.Client)
			resp, err := client.Do(request)

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
				/*
					logging for test
				*/
				logging(fmt.Sprintf("request failed with status code %d", resp.StatusCode))

				death_cnt++
				continue
			}
			if err != nil {
				/*
					logging for test
				*/
				logging(fmt.Sprintf("http.Get: %v", err))

				death_cnt++
				continue
			}
			defer resp.Body.Close()
			var page = new(instagram.Tagpage)
			// if err := json.NewDecoder(resp.Body).Decode(page); err != nil {
			// 	/*
			// 		logging for test
			// 	*/
			// 	logging(fmt.Sprintf("json: %v on requesting to %s", err, url))

			// 	death_cnt++
			// 	continue
			// }
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				/*
					logging for test
				*/
				logging(fmt.Sprintf("ioutil.ReadAll: %v while parsing from %s", err, url))

				death_cnt++
				continue
			}
			if err = json.Unmarshal(body, page); err != nil {
				/*
					logging for test
				*/
				logging(fmt.Sprintf("json.Unmarshal: %v while parsing from %s", err, url))
				logging(string(body))

				death_cnt++
				continue
			}
			var shortcodes = page.FilterByTimestamp()
			for _, shortcode := range shortcodes {
				syncer.Add(1)
				go parseUsername(shortcode, ch, syncer)
			}
			time.Sleep(0) // force context switching
			has_next_page = (len(shortcodes) == len(page.GraphQL.Hashtag.EdgeHashtagToMedia.Edges)) &&
				page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.HasNextPage
			end_cursor = page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.EndCursor
			break
		}
		if death_cnt >= max_death_cnt { // aborted when the maximum death count is exceeded
			/*
				logging for test
			*/
			logging(fmt.Sprintf("death count exceeded: %d", death_cnt))

			break
		}
	}
}

func parseUsername(shortcode string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()
	var death_cnt int
	for death_cnt < max_death_cnt {
		var url = fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1", shortcode)
		/*
			logging for test
		*/
		logging(fmt.Sprintf("requesting to %s", url))

		// resp, err := http.Get(fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1", shortcode))

		var request, err = http.NewRequest("GET", url, nil)
		if err != nil {
			/*
				logging for test
			*/
			logging(fmt.Sprintf("http.NewRequest: %v", err))

			death_cnt++
			continue
		}
		request.Header.Set("Accept", "application/json")
		var client = new(http.Client)
		resp, err := client.Do(request)

		if resp.StatusCode == http.StatusNotFound { // give up in the event of PAGE NOT FOUND
			return
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
			/*
				logging for test
			*/
			logging(fmt.Sprintf("request failed with status code %d", resp.StatusCode))

			death_cnt++
			continue
		}
		if err != nil {
			/*
				logging for test
			*/
			logging(fmt.Sprintf("http.Get: %v", err))

			death_cnt++
			continue
		}
		defer resp.Body.Close()
		var post = new(instagram.Post)
		// if err := json.NewDecoder(resp.Body).Decode(post); err != nil {
		// 	/*
		// 		logging for test
		// 	*/
		// 	logging(fmt.Sprintf("json: %v on requesting to %s", err, url))

		// 	death_cnt++
		// 	continue
		// }
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			/*
				logging for test
			*/
			logging(fmt.Sprintf("ioutil.ReadAll: %v while parsing from %s", err, url))

			death_cnt++
			continue
		}
		if err = json.Unmarshal(body, post); err != nil {
			/*
				logging for test
			*/
			logging(fmt.Sprintf("json.Unmarshal: %v while parsing from %s", err, url))
			logging(string(body))

			death_cnt++
			continue
		}
		ch <- post.GraphQL.ShortcodeMedia.Owner.Username
		return
	}
}

/*
	logging for test
*/
func logging(message interface{}) {
	file, err := os.OpenFile("nohup.out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println(message)
}

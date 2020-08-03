package crawl

import (
	"sync"
	"time"

	"github.com/haxana/vida-backend/usecase/instagram"
)

// Uploads returns a map that contains the users' upload frequency information with a given tag.
func Uploads(tag string) map[string]int {
	var usernames = make(chan string)
	var syncer = new(sync.WaitGroup)
	syncer.Add(1)
	go parsePage(tag, usernames, syncer)
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
		parser        = instagram.PageParserGenerator(tag)
	)

	for has_next_page {
		var page, err = parser()
		if err != nil {
			return
		}
		var shortcodes = page.FilterByTimestamp()
		for _, shortcode := range shortcodes {
			syncer.Add(1)
			go parseUsername(shortcode, ch, syncer)
		}
		time.Sleep(0) // force context switching
		has_next_page = (len(page.GraphQL.Hashtag.EdgeHashtagToMedia.Edges) == len(shortcodes)) &&
			page.GraphQL.Hashtag.EdgeHashtagToMedia.PageInfo.HasNextPage
	}
}

func parseUsername(shortcode string, ch chan<- string, syncer *sync.WaitGroup) {
	defer syncer.Done()

	var post, err = instagram.PostParserGenerator(shortcode)()
	if err != nil {
		return
	}
	ch <- post.GraphQL.ShortcodeMedia.Owner.Username
}

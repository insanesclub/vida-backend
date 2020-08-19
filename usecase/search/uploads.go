package search

import (
	"sync"
	"time"

	"github.com/maengsanha/vida-backend/usecase/instagram"
)

// Uploads counts the frequency of users uploading posts with a given tag.
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
		goNext = true
		parser = instagram.PageParserGenerator(tag)
	)

	for goNext {
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
		goNext = (len(page.GraphQL.Hashtag.EdgeHashtagToMedia.Edges) == len(shortcodes)) &&
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

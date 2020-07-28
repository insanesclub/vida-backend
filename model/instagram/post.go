package instagram

import (
	"errors"
	"strings"
	"time"
)

// Post represents an Instagram post.
type Post struct {
	GraphQL struct {
		ShortcodeMedia struct {
			Location struct {
				Name    string `json:"name"`
				Address string `json:"address_json"`
			} `json:"location"`
			Owner struct {
				// account name of user
				Username string `json:"username"`
				// profile name of user
				FullName                 string `json:"full_name"`
				EdgeOwnerToTimelineMedia struct {
					// number of posts
					Count int `json:"count"`
				} `json:"edge_owner_to_timeline_media"`
				EdgeFollowedBy struct {
					// number of followers
					Count int `json:"count"`
				}
			} `json:"owner"`
			TakenAtTimestamp int64 `json:"taken_at_timestamp"`
		} `json:"shortcode_media"`
	} `json:"graphql"`
}

// Geotag parses geotag from p.
func (p Post) Geotag(tag string) string {
	if err := p.filterByLocation(tag); err != nil {
		return ""
	}
	return p.GraphQL.ShortcodeMedia.Location.Name
}

func (p Post) filterByLocation(tag string) error {
	if !strings.Contains(p.GraphQL.ShortcodeMedia.Location.Address, tag) {
		return errors.New("address mismatch")
	}
	if p.GraphQL.ShortcodeMedia.Location.Name == tag {
		return errors.New("name mismatch")
	}
	return nil
}

// PostedToday checks whether p is posted today.
func (p Post) PostedToday() bool {
	var now = time.Now().Unix()
	return now-now%86400 <= p.GraphQL.ShortcodeMedia.TakenAtTimestamp &&
		p.GraphQL.ShortcodeMedia.TakenAtTimestamp < now-now%86400+86400
}

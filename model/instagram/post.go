package instagram

import (
	"errors"
	"strings"
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
func (p Post) Geotag(query string) string {
	if err := p.filter(query); err != nil {
		return ""
	}
	return p.GraphQL.ShortcodeMedia.Location.Name
}

func (p Post) filter(query string) error {
	if !strings.Contains(p.GraphQL.ShortcodeMedia.Location.Address, query) {
		return errors.New("address mismatch")
	}
	if p.GraphQL.ShortcodeMedia.Location.Name == query {
		return errors.New("name mismatch")
	}
	return nil
}

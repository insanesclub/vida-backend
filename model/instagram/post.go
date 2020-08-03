package instagram

import "strings"

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
		} `json:"shortcode_media"`
	} `json:"graphql"`
}

// FilterByLocation returns geotag of p.
func (p Post) FilterByLocation(tag string) (geotag string) {
	if strings.Contains(p.GraphQL.ShortcodeMedia.Location.Address, tag) && p.GraphQL.ShortcodeMedia.Location.Name != tag {
		geotag = p.GraphQL.ShortcodeMedia.Location.Name
	}
	return
}

package instagram

import "strings"

// Post represents an Instagram post.
type Post struct {
	GraphQL struct {
		ShortcodeMedia struct {
			Location struct {
				Name        string `json:"name"`
				AddressJSON string `json:"address_json"`
			} `json:"location"`
			Owner struct {
				Username                 string `json:"username"`  // account name of user
				FullName                 string `json:"full_name"` // profile name of user
				EdgeOwnerToTimelineMedia struct {
					Count int `json:"count"` // number of posts
				} `json:"edge_owner_to_timeline_media"`
				EdgeFollowedBy struct {
					Count int `json:"count"` // number of followers
				} `json:"edge_followed_by"`
			} `json:"owner"`
		} `json:"shortcode_media"`
	} `json:"graphql"`
}

// FilterByLocation returns geotag of p if.
func (p Post) FilterByLocation(tag string) (_ string) {
	if strings.Contains(p.GraphQL.ShortcodeMedia.Location.AddressJSON,
		p.GraphQL.ShortcodeMedia.Location.Name) {
		return
	}
	return p.GraphQL.ShortcodeMedia.Location.Name
}

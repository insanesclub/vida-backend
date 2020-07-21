// Package instagram defines Instagram object models.
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
		return errors.New("Address mismatch")
	}
	if p.GraphQL.ShortcodeMedia.Location.Name == query {
		return errors.New("Name mismatch")
	}
	return nil
}

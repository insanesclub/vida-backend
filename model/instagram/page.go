package instagram

// Tagpage represents an Instagram page source.
type Tagpage struct {
	GraphQL struct {
		Hashtag struct {
			EdgeHashtagToMedia struct {
				PageInfo struct {
					HasNextPage bool   `json:"has_next_page"`
					EndCursor   string `json:"end_cursor,omitempty"`
				} `json:"page_info"`
				Edges []struct {
					Node struct {
						Shortcode string `json:"shortcode"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_hashtag_to_media"`
		} `json:"hashtag"`
	} `json:"graphql"`
}

// Shortcodes parses URLs from t.
func (t Tagpage) Shortcodes() (shortcodes []string) {
	for _, edge := range t.GraphQL.Hashtag.EdgeHashtagToMedia.Edges {
		shortcodes = append(shortcodes, edge.Node.Shortcode)
	}
	return
}

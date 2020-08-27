package instagram

import "time"

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
						Shortcode        string `json:"shortcode"`
						TakenAtTimestamp int64  `json:"taken_at_timestamp"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_hashtag_to_media"`
		} `json:"hashtag"`
	} `json:"graphql"`
}

// FilterByTimestamp returns the shortcodes of the posts posted today on t.
func (t Tagpage) FilterByTimestamp() (shortcodes []string) {
	now := time.Now().Unix()
	for _, edge := range t.GraphQL.Hashtag.EdgeHashtagToMedia.Edges {
		if now-now%86400 <= edge.Node.TakenAtTimestamp &&
			edge.Node.TakenAtTimestamp < now-now%86400+86400 {
			shortcodes = append(shortcodes, edge.Node.Shortcode)
		}
	}
	return
}

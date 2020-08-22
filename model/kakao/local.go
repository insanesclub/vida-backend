package kakao

// LocalAPIResult represents the results of a place search provided by kakao local API.
// See https://developers.kakao.com/docs/latest/ko/local/dev-guide#search-by-keyword
type LocalAPIResult struct {
	Documents []struct {
		ID                string `json:"id"`
		PlaceName         string `json:"place_name"`
		CategoryName      string `json:"category_name"`
		CategoryGroupCode string `json:"category_group_code"`
		CategoryGroupName string `json:"category_group_name"`
		Phone             string `json:"phone"`
		Address           string `json:"address"`
		RoadAddressName   string `json:"road_adress_name"`
		CoordX            string `json:"x"`
		CoordY            string `json:"y"`
		PlaceURL          string `json:"place_url"`
		Distance          string `json:"distance"`
	} `json:"documents"`
	Meta struct {
		IsEnd         bool `json:"is_end"`
		PageableCount int  `json:"pageable_count"`
		SameName      struct {
			Keyword        string   `json:"keyword"`
			Region         []string `json:"region"`
			SelectedRegion string   `json:"selected_region"`
		} `json:"same_name"`
		TotalCount int `json:"total_count"`
	} `json:"meta"`
}

// Places returns place URLs in l.
func (l LocalAPIResult) Places() []string {
	places := make([]string, len(l.Documents))
	for idx, doc := range l.Documents {
		places[idx] = doc.PlaceURL
	}
	return places
}

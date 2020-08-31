package kakao

// MapInfo represents the results of a place search provided by kakao local API.
// See https://developers.kakao.com/docs/latest/ko/local/dev-guide#search-by-keyword
type MapInfo struct {
	Documents []Document `json:"documents"`
	Meta      struct {
		IsEnd bool `json:"is_end"`
	} `json:"meta"`
}

type Document struct {
	ID                string `json:"id"`
	PlaceName         string `json:"place_name"`
	CategoryName      string `json:"category_name"`
	CategoryGroupCode string `json:"category_group_code"`
	CategoryGroupName string `json:"category_group_name"`
	Phone             string `json:"phone"`
	AddressName       string `json:"address_name"`
	RoadAddressName   string `json:"road_address_name"`
}

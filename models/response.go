package models

// PageResponse ...
type PageResponse struct {
	Records     int64       `json:"records"`
	PageLimit   int64       `json:"pageLimit"`
	CurrentPage int64       `json:"currentPage"`
	Data        interface{} `json:"data"`
}

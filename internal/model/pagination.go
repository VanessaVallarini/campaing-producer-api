package model

// Page is the current page number
// Size is how many items there is in the current page
// Total items are the total amount of items in the whole table
type Paging struct {
	Page       int `json:"page"`
	Size       int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
}

package beans

type PaginationBean struct {
	TotalItems int64         `json:"totalItems"`
	Offset     int           `json:"offset"`
	Items      []interface{} `json:"items"`
}

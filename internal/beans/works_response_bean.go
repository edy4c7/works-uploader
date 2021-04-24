package beans

import "time"

type WorksResponseBean struct {
	ID     uint64 `json:"id"`
	Title  string `json:"title"`
	Author struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	Description  string    `json:"description"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	ContentURL   string    `json:"contentUrl"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

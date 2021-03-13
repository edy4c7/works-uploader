package entities

import "time"

type Work struct {
	ID                uint64
	Title             string
	Author            string
	Description       string
	ThumbnailFileName string
	ContentFileName   string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

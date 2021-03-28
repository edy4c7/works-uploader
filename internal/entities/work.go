package entities

import (
	"time"

	"github.com/edy4c7/darkpot-school-works/internal/common/constants"
)

type Work struct {
	ID           uint64
	Type         constants.WorkType
	Title        string
	Author       string
	Description  string
	ThumbnailURL string
	ContentURL   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

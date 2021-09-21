package entities

import (
	"time"

	"github.com/edy4c7/works-uploader/internal/common/constants"
	"gorm.io/gorm"
)

type Work struct {
	ID           uint64
	Type         constants.WorkType
	Title        string `size:"40"`
	Author       string
	Description  string `size:"200"`
	ThumbnailURL string
	ContentURL   string
	Version      uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

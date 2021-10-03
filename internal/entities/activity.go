package entities

import (
	"time"

	"github.com/edy4c7/works-uploader/internal/common/constants"
)

type Activity struct {
	ID        uint64
	Type      constants.ActivityType
	UserID    string `json:"-"`
	User      *User
	WorkID    uint64 `json:"-"`
	Work      *Work
	CreatedAt time.Time
}

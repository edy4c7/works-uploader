package entities

import (
	"time"

	"github.com/edy4c7/darkpot-school-works/internal/common/constants"
)

type Activity struct {
	ID        uint64
	Type      constants.ActivityType
	User      string
	Work      *Work
	CreatedAt time.Time
}

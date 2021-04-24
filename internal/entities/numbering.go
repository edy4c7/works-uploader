package entities

import "time"

type Numbering struct {
	Type      int
	Number    uint64
	CreatedAt time.Time `firestore:"CreatedAt,serverTimestamp"`
	UpdatedAt time.Time `firestore:"UpdatedAt,serverTimestamp"`
}

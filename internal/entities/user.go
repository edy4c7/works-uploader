package entities

import "time"

type User struct {
	ID        string `gorm:"primaryKey" json:"-"`
	Name      string
	Nickname  string
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

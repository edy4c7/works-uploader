package entities

import "time"

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Nickname  string
	Picture   string
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

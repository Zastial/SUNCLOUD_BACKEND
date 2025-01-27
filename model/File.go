package model

import "time"

type File struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Path      string    `gorm:"not null;unique" json:"path"`
	Size      int64     `gorm:"not null" json:"size"`
	Type      string    `gorm:"not null" json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

package structs

import "time"

type Comment struct {
	ID        int64
	PostID    int64
	Post      Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"default:null"`
	Author    string
	Content   string
}

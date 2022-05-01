package structs

import (
	"time"
)

type Post struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"default:null"`
	Author    string
	Content   string
	Like      int
	Dislike   int
}

package model

import (
	"time"
)

type Project struct {
	Id          int       `json:"id" bson:"id"`
	Name        string    `bson:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index" bson:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;index" bson:"updatedAt"`
}

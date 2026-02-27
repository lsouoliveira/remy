package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Content   string `gorm:"type:text;not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

func CreateNote(content string) *Note {
	return &Note{
		Content: content,
	}
}

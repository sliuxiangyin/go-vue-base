package models

import (
	"time"

	"gorm.io/gorm"
)

// PhoneticDictionary 发音词典（全局共享）
type PhoneticDictionary struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Word string `gorm:"type:varchar(100);uniqueIndex;not null;comment:单词" json:"word"`
	IPA  string `gorm:"type:varchar(100);not null;comment:国际音标" json:"ipa"`
}

// TableName 指定表名
func (PhoneticDictionary) TableName() string {
	return "phonetic_dictionary"
}

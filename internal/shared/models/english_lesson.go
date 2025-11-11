package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// EnglishLesson 英语课程模型
type EnglishLesson struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title     string  `gorm:"type:varchar(255);not null;comment:文章或片段标题" json:"title"`
	AudioURL  string  `gorm:"type:varchar(512);comment:音频文件路径" json:"audio_url"`
	Duration  float64 `gorm:"type:float;comment:音频时长（秒）" json:"duration"`
	ContentEN string  `gorm:"type:text;comment:原英文内容" json:"content_en"`
	ContentZH string  `gorm:"type:text;comment:中文翻译" json:"content_zh"`

	// JSON 字段
	SemanticJSON      SemanticChunks `gorm:"type:json;comment:语义意群与时间戳匹配结果" json:"semantic_json"`
	WordTimestampJSON WordTimestamps `gorm:"type:json;comment:Whisper输出的逐词时间戳" json:"word_timestamp_json"`
	PhoneticJSON      Phonetics      `gorm:"type:json;comment:发音词典" json:"phonetic_json"`
	Tags              Tags           `gorm:"type:json;comment:分类标签" json:"tags"`

	Level    int8 `gorm:"type:tinyint;comment:难度等级（1~5）" json:"level"`
	IsPublic bool `gorm:"type:tinyint;default:0;comment:是否公开" json:"is_public"`
}

// TableName 指定表名
func (EnglishLesson) TableName() string {
	return "english_lessons"
}

// SemanticChunk 语义意群结构
type SemanticChunk struct {
	Text  string  `json:"text"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// SemanticChunks 语义意群数组，实现 sql.Scanner 和 driver.Valuer 接口
type SemanticChunks []SemanticChunk

func (s SemanticChunks) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *SemanticChunks) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// WordTimestamp 逐词时间戳结构
type WordTimestamp struct {
	Word       string  `json:"word"`
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Confidence float64 `json:"confidence,omitempty"`
}

// WordTimestamps 逐词时间戳数组
type WordTimestamps []WordTimestamp

func (w WordTimestamps) Value() (driver.Value, error) {
	if w == nil {
		return nil, nil
	}
	return json.Marshal(w)
}

func (w *WordTimestamps) Scan(value interface{}) error {
	if value == nil {
		*w = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, w)
}

// Phonetic 发音信息结构
type Phonetic struct {
	Word string `json:"word"`
	IPA  string `json:"ipa"`
}

// Phonetics 发音信息数组
type Phonetics []Phonetic

func (p Phonetics) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

func (p *Phonetics) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// Tags 标签数组
type Tags []string

func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}

func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, t)
}

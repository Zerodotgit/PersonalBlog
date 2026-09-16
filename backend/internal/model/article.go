package model

import "time"

const (
	StatusDraft    = 0
	StatusPublish  = 1
)

// Article 表示一篇博客文章，content 存储 Markdown 原文。
type Article struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"size:255;not null;index" json:"title"`
	Summary    string    `gorm:"size:500" json:"summary"`
	Content    string    `gorm:"type:mediumtext;not null" json:"content"`
	CoverImage string    `gorm:"size:500" json:"cover_image"`
	Tags       string    `gorm:"size:255" json:"tags"`
	Status     int8      `gorm:"size:1;default:1;index" json:"status"`
	ViewCount  int64     `gorm:"default:0" json:"view_count"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Summary 对外返回列表/详情时使用的精简结构，避免在大列表中携带正文。
type ArticleListItem struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	CoverImage string    `json:"cover_image"`
	Tags       string    `json:"tags"`
	Status     int8      `json:"status"`
	ViewCount  int64     `json:"view_count"`
	CreatedAt  time.Time `json:"created_at"`
}

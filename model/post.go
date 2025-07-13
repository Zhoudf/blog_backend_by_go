package model

import (
	"gorm.io/gorm"
)

// Post 模型定义博客文章信息
type Post struct {
	gorm.Model
	Title    string    `gorm:"size:100;not null" json:"title"`              // 文章标题，不为空
	Content  string    `gorm:"type:text;not null" json:"content"`           // 文章内容，不为空
	UserID   uint      `gorm:"not null" json:"user_id"`                     // 关联用户ID，不为空
	User     User      `gorm:"foreignKey:UserID" json:"user"`               // 关联用户信息
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"` // 文章的评论列表
}

func (Post) TableName() string {
	return "gin_posts"
}

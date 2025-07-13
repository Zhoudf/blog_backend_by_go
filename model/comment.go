package model

import (
	"gorm.io/gorm"
)

// Comment 模型定义文章评论信息
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"`       // 评论内容，不为空
	UserID  uint   `gorm:"not null" json:"user_id"`                 // 关联用户ID，不为空
	PostID  uint   `gorm:"not null" json:"post_id"`                 // 关联文章ID，不为空
	User    User   `gorm:"foreignKey:UserID" json:"user"`           // 关联用户信息
	Post    Post   `gorm:"foreignKey:PostID" json:"post,omitempty"` // 关联文章信息
}

func (Comment) TableName() string {
	return "gin_comments"
}

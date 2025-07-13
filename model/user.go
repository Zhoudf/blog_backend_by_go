package model

import (
	"gorm.io/gorm"
)

// User 模型定义用户信息
type User struct {
	gorm.Model
	Username string    `gorm:"size:50;uniqueIndex;not null" json:"username"` // 用户名，唯一且不为空
	Password string    `gorm:"not null" json:"-"`                            // 密码，不为空，JSON序列化时忽略
	Email    string    `gorm:"size:100;uniqueIndex;not null" json:"email"`   // 邮箱，唯一且不为空
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`     // 用户发布的文章
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`  // 用户发表的评论
}

func (User) TableName() string {
	return "gin_users"
}

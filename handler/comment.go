package handler

import (
	"net/http"
	"strconv"

	"github.com/Zhoudf/blog_backend_by_go/config"
	"github.com/Zhoudf/blog_backend_by_go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCommentRequest 创建评论请求参数
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 获取文章ID - 修改为使用统一的:id参数
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID不能为空"})
		return
	}

	// 验证文章是否存在
	var post model.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败"})
		}
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 转换postID为uint
	postIDUint, err := strconv.ParseUint(postID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 创建评论
	comment := model.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postIDUint),
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论创建失败"})
		return
	}

	// 查询创建的评论，包含用户信息
	if err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username")
	}).First(&comment, comment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论查询失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": comment,
	})
}

// GetComments 获取文章评论列表
func GetComments(c *gin.Context) {
	// 获取文章ID - 修改为使用统一的:id参数
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID不能为空"})
		return
	}

	// 验证文章是否存在
	var post model.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败"})
		}
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	// 查询评论列表，包含用户信息，按创建时间倒序
	var comments []model.Comment
	var total int64

	// 获取总数
	config.DB.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&total)

	// 获取分页数据
	if err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username")
	}).Where("post_id = ?", postID).Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"pagination": gin.H{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/Zhoudf/blog_backend_by_go/config"
	"github.com/Zhoudf/blog_backend_by_go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePostRequest 创建文章请求参数
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required"`
}

// UpdatePostRequest 更新文章请求参数
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=100"`
	Content string `json:"content" binding:"omitempty"`
}

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 创建文章
	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文章创建失败"})
		return
	}

	// 查询创建的文章，包含用户信息
	if err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username") // 只查询用户ID和用户名，保护隐私
	}).First(&post, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post":    post,
	})
}

// GetPosts 获取文章列表
func GetPosts(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 查询文章列表，包含用户信息，分页
	var posts []model.Post
	var total int64

	// 获取总数
	config.DB.Model(&model.Post{}).Count(&total)

	// 获取分页数据
	if err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username")
	}).Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetPost 获取文章详情
func GetPost(c *gin.Context) {
	// 获取文章ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID不能为空"})
		return
	}

	// 查询文章，包含用户信息
	var post model.Post
	if err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username")
	}).First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 获取文章ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID不能为空"})
		return
	}

	// 查询文章
	var post model.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败"})
		}
		return
	}

	// 检查是否为文章作者
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限更新此文章"})
		return
	}

	// 绑定更新数据
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 更新文章
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文章更新失败"})
		return
	}

	// 查询更新后的文章，包含用户信息
	if err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Username")
	}).First(&post, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post":    post,
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 获取文章ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID不能为空"})
		return
	}

	// 查询文章
	var post model.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文章查询失败"})
		}
		return
	}

	// 检查是否为文章作者
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文章"})
		return
	}

	// 删除文章
	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文章删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}

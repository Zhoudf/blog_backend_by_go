package handler

import (
	"github.com/Zhoudf/blog_backend_by_go/config"
	"github.com/Zhoudf/blog_backend_by_go/model"
	"github.com/Zhoudf/blog_backend_by_go/utils"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的邮箱格式"})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	} else if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	} else if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新用户
	newUser := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户注册成功",
		"user": gin.H{
			"id":       newUser.ID,
			"username": newUser.Username,
			"email":    newUser.Email,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 查询用户
	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		}
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

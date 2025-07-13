# 博客系统后端

一个使用Go语言、Gin框架和GORM库开发的个人博客系统后端，支持用户认证、文章管理和评论功能。

## 功能特点

- 用户认证：注册、登录（JWT认证）
- 文章管理：创建、读取、更新、删除文章
- 评论功能：发表评论、获取文章评论列表
- 错误处理：统一错误响应格式
- 日志记录：请求日志和错误日志

## 技术栈

- 语言：Go 1.16+
- Web框架：Gin
- ORM：GORM
- 数据库：MySQL
- 认证：JWT
- 密码加密：bcrypt

## 项目结构

```
blog-backend/
├── main.go           # 程序入口
├── go.mod            # 依赖管理
├── go.sum            # 依赖校验
├── README.md         # 项目说明
├── model/            # 数据模型
│   ├── user.go       # 用户模型
│   ├── post.go       # 文章模型
│   └── comment.go    # 评论模型
├── handler/          # 处理器
│   ├── auth.go       # 用户认证处理器
│   ├── post.go       # 文章处理器
│   └── comment.go    # 评论处理器
├── middleware/       # 中间件
│   ├── auth.go       # JWT认证中间件
│   ├── logger.go     # 日志中间件
│   └── recovery.go   # 错误恢复中间件
├── config/           # 配置
│   └── database.go   # 数据库配置
├── utils/            # 工具函数
│   ├── jwt.go        # JWT工具
│   └── password.go   # 密码工具
└── routes/           # 路由
    └── router.go     # 路由配置
```

## 环境要求

- Go 1.16+
- MySQL 5.7+ 或 8.0+

## 安装步骤

### 1. 克隆代码

```bash
git clone https://github.com/Zhoudf/blog_backend_by_go.git
cd blog_backend_by_go
```

### 2. 初始化依赖

```bash
go mod tidy
```

### 3. 数据库配置

#### 3.1 创建数据库

```sql
CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 3.2 配置数据库连接

可以通过环境变量配置数据库连接：

```bash
export DB_DSN="用户名:密码@tcp(地址:端口)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
```

或者直接修改`config/database.go`文件中的默认DSN。

### 4. 配置JWT密钥

```bash
export JWT_SECRET="your-secret-key"  # 替换为自己的密钥
```

## 启动方法

### 开发环境

```bash
go run main.go
```

### 生产环境

#### 构建可执行文件

```bash
go build -o blog-backend
```

#### 运行

```bash
./blog-backend
```

默认端口为8080，可以通过环境变量修改：

```bash
export PORT=8000
./blog-backend
```

## 测试方法
以下文件可以导入到postman
blog_test.postman_collection.json


## API接口文档

### 认证接口

#### 注册用户

- URL: `/api/auth/register`
- 方法: POST
- 请求体:
  ```json
  {
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
  }
  ```
- 响应: 201 Created
  ```json
  {
    "message": "用户注册成功",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
  ```

#### 用户登录

- URL: `/api/auth/login`
- 方法: POST
- 请求体:
  ```json
  {
    "username": "testuser",
    "password": "password123"
  }
  ```
- 响应: 200 OK
  ```json
  {
    "message": "登录成功",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
  ```

### 文章接口

#### 获取文章列表

- URL: `/api/posts?page=1&page_size=10`
- 方法: GET
- 响应: 200 OK
  ```json
  {
    "posts": [
      {
        "ID": 1,
        "CreatedAt": "2023-07-01T10:00:00Z",
        "UpdatedAt": "2023-07-01T10:00:00Z",
        "DeletedAt": null,
        "title": "文章标题",
        "content": "文章内容...",
        "user_id": 1,
        "user": {
          "ID": 1,
          "Username": "testuser"
        }
      }
    ],
    "pagination": {
      "total": 10,
      "page": 1,
      "page_size": 10,
      "total_pages": 1
    }
  }
  ```

#### 获取文章详情

- URL: `/api/posts/:id`
- 方法: GET
- 响应: 200 OK
  ```json
  {
    "post": {
      "ID": 1,
      "CreatedAt": "2023-07-01T10:00:00Z",
      "UpdatedAt": "2023-07-01T10:00:00Z",
      "DeletedAt": null,
      "title": "文章标题",
      "content": "文章内容...",
      "user_id": 1,
      "user": {
        "ID": 1,
        "Username": "testuser"
      }
    }
  }
  ```

#### 创建文章

- URL: `/api/posts`
- 方法: POST
- 请求头: `Authorization: Bearer {token}`
- 请求体:
  ```json
  {
    "title": "新文章标题",
    "content": "新文章内容..."
  }
  ```
- 响应: 201 Created
  ```json
  {
    "message": "文章创建成功",
    "post": {
      "ID": 2,
      "CreatedAt": "2023-07-01T11:00:00Z",
      "UpdatedAt": "2023-07-01T11:00:00Z",
      "DeletedAt": null,
      "title": "新文章标题",
      "content": "新文章内容...",
      "user_id": 1,
      "user": {
        "ID": 1,
        "Username": "testuser"
      }
    }
  }
  ```

#### 更新文章

- URL: `/api/posts/:id`
- 方法: PUT
- 请求头: `Authorization: Bearer {token}`
- 请求体:
  ```json
  {
    "title": "更新后的标题",
    "content": "更新后的内容..."
  }
  ```
- 响应: 200 OK
  ```json
  {
    "message": "文章更新成功",
    "post": {
      "ID": 1,
      "CreatedAt": "2023-07-01T10:00:00Z",
      "UpdatedAt": "2023-07-01T12:00:00Z",
      "DeletedAt": null,
      "title": "更新后的标题",
      "content": "更新后的内容...",
      "user_id": 1,
      "user": {
        "ID": 1,
        "Username": "testuser"
      }
    }
  }
  ```

#### 删除文章

- URL: `/api/posts/:id`
- 方法: DELETE
- 请求头: `Authorization: Bearer {token}`
- 响应: 200 OK
  ```json
  {
    "message": "文章删除成功"
  }
  ```

### 评论接口

#### 获取文章评论列表

- URL: `/api/posts/:post_id/comments?page=1&page_size=20`
- 方法: GET
- 响应: 200 OK
  ```json
  {
    "comments": [
      {
        "ID": 1,
        "CreatedAt": "2023-07-01T13:00:00Z",
        "UpdatedAt": "2023-07-01T13:0０:00Z",
        "DeletedAt": null,
        "content": "这是一条评论",
        "user_id": 1,
        "post_id": 1,
        "user": {
          "ID": 1,
          "Username": "testuser"
        }
      }
    ],
    "pagination": {
      "total": 5,
      "page": 1,
      "page_size": 20,
      "total_pages": 1
    }
  }
  ```

#### 创建评论

- URL: `/api/posts/:post_id/comments`
- 方法: POST
- 请求头: `Authorization: Bearer {token}`
- 请求体:
  ```json
  {
    "content": "这是一条新评论"
  }
  ```
- 响应: 201 Created  
  ```json
  {
    "message": "评论创建成功",
    "comment": {
      "ID": 2,
      "CreatedAt": "2023-07-01T14:00:00Z",
      "UpdatedAt": "2023-07-01T14:00:00Z",
      "DeletedAt": null,
      "content": "这是一条新评论",
      "user_id": 1,
      "post_id": 1,
      "user": {
        "ID": 1,
        "Username": "testuser"
      }
    }
  }
  ```

## 错误响应格式

所有错误响应将采用统一格式：

```json
{
  "error": "错误描述信息"
}
```

常见的HTTP状态码：
- 400: 请求参数错误
- 401: 未认证或认证失败
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器内部错误
博客系统API测试用例
测试环境
基础URL: http://localhost:8080/api
测试工具: Postman
数据库: MySQL
测试用例设计
1. 用户认证接口
1.1 用户注册
用例ID	接口URL	方法	请求头	请求体	预期状态码	预期响应
AUTH-001	/auth/register	POST	Content-Type: application/json	{"username":"testuser","password":"Password123","email":"test@example.com"}	201	包含"用户注册成功"消息和用户信息
AUTH-002	/auth/register	POST	Content-Type: application/json	{"username":"testuser","password":"123","email":"invalid-email"}	400	包含"无效的请求参数"错误
AUTH-003	/auth/register	POST	Content-Type: application/json	{"username":"testuser","password":"Password123","email":"test@example.com"}	409	包含"用户名已存在"错误
1.2 用户登录
用例ID	接口URL	方法	请求头	请求体	预期状态码	预期响应
AUTH-004	/auth/login	POST	Content-Type: application/json	{"username":"testuser","password":"Password123"}	200	包含"登录成功"消息、token和用户信息
AUTH-005	/auth/login	POST	Content-Type: application/json	{"username":"nonexistent","password":"Password123"}	401	包含"用户名或密码错误"错误
AUTH-006	/auth/login	POST	Content-Type: application/json	{"username":"testuser","password":"WrongPassword"}	401	包含"用户名或密码错误"错误
2. 文章管理接口
2.1 创建文章
用例ID	接口URL	方法	请求头	请求体	预期状态码	预期响应
POST-001	/posts	POST	Content-Type: application/json<br>Authorization: Bearer {valid_token}	{"title":"测试文章标题","content":"这是测试文章的内容，至少需要10个字符。"}	201	包含"文章创建成功"消息和文章信息
POST-002	/posts	POST	Content-Type: application/json	{"title":"无认证创建文章","content":"这应该失败"}	401	包含"未认证"错误
POST-003	/posts	POST	Content-Type: application/json<br>Authorization: Bearer {valid_token}	{"title":"短标题","content":"短"}	400	包含"无效的请求参数"错误
2.2 获取文章列表
用例ID	接口URL	方法	请求头	请求参数	预期状态码	预期响应
POST-004	/posts	GET	无	page=1&page_size=10	200	包含文章列表和分页信息
POST-005	/posts	GET	无	page=999&page_size=10	200	包含空文章列表和分页信息
2.3 获取文章详情
用例ID	接口URL	方法	请求头	路径参数	预期状态码	预期响应
POST-006	/posts/{id}	GET	无	id=1 (存在的文章ID)	200	包含文章详情
POST-007	/posts/{id}	GET	无	id=999 (不存在的文章ID)	404	包含"文章不存在"错误
2.4 更新文章
用例ID	接口URL	方法	请求头	路径参数 & 请求体	预期状态码	预期响应
POST-008	/posts/{id}	PUT	Content-Type: application/json<br>Authorization: Bearer {owner_token}	id=1 (自己的文章)<br>{"title":"更新后的标题","content":"更新后的内容"}	200	包含"文章更新成功"消息和更新后的文章信息
POST-009	/posts/{id}	PUT	Content-Type: application/json<br>Authorization: Bearer {other_token}	id=1 (他人的文章)<br>{"title":"尝试更新他人文章","content":"这应该失败"}	403	包含"没有权限更新此文章"错误
POST-010	/posts/{id}	PUT	Content-Type: application/json<br>Authorization: Bearer {valid_token}	id=999 (不存在的文章)<br>{"title":"更新不存在的文章","content":"这应该失败"}	404	包含"文章不存在"错误
2.5 删除文章
用例ID	接口URL	方法	请求头	路径参数	预期状态码	预期响应
POST-011	/posts/{id}	DELETE	Authorization: Bearer {owner_token}	id=1 (自己的文章)	200	包含"文章删除成功"消息
POST-012	/posts/{id}	DELETE	Authorization: Bearer {other_token}	id=1 (他人的文章)	403	包含"没有权限删除此文章"错误
POST-013	/posts/{id}	DELETE	Authorization: Bearer {valid_token}	id=999 (不存在的文章)	404	包含"文章不存在"错误
3. 评论接口
3.1 创建评论
用例ID	接口URL	方法	请求头	路径参数 & 请求体	预期状态码	预期响应
COMMENT-001	/posts/{id}/comments	POST	Content-Type: application/json<br>Authorization: Bearer {valid_token}	id=1 (存在的文章ID)<br>{"content":"这是一条测试评论"}	201	包含"评论创建成功"消息和评论信息
COMMENT-002	/posts/{id}/comments	POST	Content-Type: application/json	id=1<br>{"content":"无认证创建评论"}	401	包含"未认证"错误
COMMENT-003	/posts/{id}/comments	POST	Content-Type: application/json<br>Authorization: Bearer {valid_token}	id=999 (不存在的文章ID)<br>{"content":"对不存在文章的评论"}	404	包含"文章不存在"错误
3.2 获取文章评论列表
用例ID	接口URL	方法	请求头	路径参数 & 请求参数	预期状态码	预期响应
COMMENT-004	/posts/{id}/comments	GET	无	id=1 (存在的文章ID)<br>page=1&page_size=20	200	包含评论列表和分页信息
COMMENT-005	/posts/{id}/comments	GET	无	id=999 (不存在的文章ID)	404	包含"文章不存在"错误


Postman集合导出方法
在Postman中创建一个新集合"Blog Backend API Test"
按照上述测试用例添加请求
点击集合菜单，选择"Export"
选择"Collection v2.1"格式导出JSON文件
测试注意事项
测试前确保服务器已启动且数据库连接正常
测试用户注册后，使用该用户登录获取token用于后续测试
测试更新和删除文章时，需要先创建文章并记录文章ID
测试权限控制时，需要创建至少两个不同的用户账号
测试完成后，可以使用以下SQL清理测试数据：
DELETE FROM comments WHERE user_id IN (SELECT id FROM users WHERE username = 'testuser');
DELETE FROM posts WHERE user_id IN (SELECT id FROM users WHERE username = 'testuser');
DELETE FROM users WHERE username = 'testuser';
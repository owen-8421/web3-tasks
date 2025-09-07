# **Go-Blog-Backend 个人博客系统后端**

项目采用 [Gin](https://github.com/gin-gonic/gin) 作为 Web 框架，[GORM](https://gorm.io/) 作为 ORM 库与 MySQL 数据库交互，并使用 JWT 实现用户认证。

## **✨ 主要功能**

* **用户认证**:  
  * 用户注册（密码自动加密）  
  * 用户登录  
  * 基于 JWT 的状态管理和接口授权  
* **文章管理 (CRUD)**:  
  * 创建文章（需要认证）  
  * 读取文章列表（公开）  
  * 读取单篇文章详情（公开）  
  * 更新文章（仅限作者）  
  * 删除文章（仅限作者）  
* **评论功能**:  
  * 创建评论（需要认证）  
  * 读取某篇文章下的所有评论（公开）

## **📂 项目结构**
.  
├── cmd/server/main.go          \# 🚀 程序主入口  
├── internal  
│   ├── api  
│   │   ├── handler             \# 📦 (Handler层) 负责处理HTTP请求和响应  
│   │   ├── middleware          \# 🛡️ (中间件) 如JWT认证  
│   │   └── router.go           \# 🛣️ (路由) 定义所有API端点  
│   ├── model                   \# 🏗️ (模型层) 定义数据库表结构  
│   ├── pkg                     \# 🧰 (公共包) 存放配置、数据库、JWT等工具  
│   ├── repository              \# 🗄️ (数据访问层) 封装对数据库的直接操作  
│   └── service                 \# 🧠 (业务逻辑层) 处理核心业务逻辑  
├── setting.yaml                \# ⚙️ 项目配置文件  
└── go.mod                      \# Go模块管理



## **🚀 快速开始**

### **1\. 环境准备**

* 安装 [Go](https://golang.org/) (建议版本 \>= 1.18)  
* 安装 [MySQL](https://www.mysql.com/) 数据库  
* 在 MySQL 中创建一个数据库，例如 blog\_db



### **2\. 快速开始**

### 2.1 安装项目依赖  
go mod tidy

### 2.2 启动服务  
go run ./cmd/server/main.go

启动成功后，你将看到服务在指定端口（例如 :9080）上运行的日志。

## **📖 API 接口文档**

所有接口的基础路径为 /api/v1。

### **用户认证**

| 功能 | 方法 | 路径 | 认证 | 请求体 |
| :---- | :---- | :---- | :---- | :---- |
| 用户注册 | POST | /register | 否 | {"username": "user1", "password": "password123"} |
| 用户登录 | POST | /login | 否 | {"username": "user1", "password": "password123"} |

### **文章管理**

| 功能 | 方法 | 路径 | 认证 | 请求体 |
| :---- | :---- | :---- | :---- | :---- |
| 创建文章 | POST | /posts | **是** | {"title": "文章标题", "content": "文章内容"} |
| 获取文章列表 | GET | /posts | 否 | N/A |
| 获取单篇文章 | GET | /posts/:id | 否 | N/A |
| 更新文章 | PUT | /posts/:id | **是** (仅作者) | {"title": "更新的标题", "content": "更新的内容"} |
| 删除文章 | DELETE | /posts/:id | **是** (仅作者) | N/A |

### **评论管理**

| 功能 | 方法 | 路径 | 认证 | 请求体 |
| :---- | :---- | :---- | :---- | :---- |
| 创建评论 | POST | /posts/:id/comments | **是** | {"content": "这是我的评论内容"} |
| 获取文章评论 | GET | /posts/:id/comments | 否 | N/A |



package task3

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

/*
题目1：模型定义

假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。

要求 ：

使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。

编写Go代码，使用Gorm创建这些模型对应的数据库表。
*/

// User 用户模型
type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100);not null"`
	Posts        []Post // 一个用户可以有多篇文章 (一对多)
	ArticleCount int    `gorm:"default:0"` // 新增字段：文章数量统计
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title         string    `gorm:"type:varchar(200);not null"`
	Content       string    `gorm:"type:text;not null"`
	UserID        uint      // 外键
	Comments      []Comment // 一篇文章可以有多条评论 (一对多)
	CommentStatus string    `gorm:"type:varchar(50);default:'无评论'"` // 新增字段：评论状态
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	PostID  uint   // 外键
}

func test3() {
	// 连接到gorm_blog数据库
	db, err := gorm.Open(sqlite.Open("gorm_blog.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	fmt.Println("数据库连接成功！")

	// 自动迁移模型，Gorm会创建表、缺失的外键、约束、列和索引。
	// 它只会添加缺失的东西，不会删除或更改现有的列/索引。

	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	fmt.Println("数据库表结构创建/更新成功！")
}

/*
题目2：关联查询
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。
*/

func SearchComment(db *gorm.DB, userName string) {
	// --- 1. 查询某个用户发布的所有文章及其对应的评论信息 ---，userName传参张三
	fmt.Println("--- 1. 查询用户'张三'的所有文章及评论 ---")
	var user User
	// 使用 Preload 来预加载关联数据
	// Preload("Posts") 会加载该用户的所有文章
	// Preload("Posts.Comments") 会在加载文章的同时，加载每篇文章的评论
	if err := db.Preload("Posts.Comments").First(&user, "name = ?", userName).Error; err != nil {
		log.Printf("查询用户失败: %v", err)
	} else {
		fmt.Printf("用户名: %s\n", user.Name)
		for _, post := range user.Posts {
			fmt.Printf("  文章: %s\n", post.Title)
			if len(post.Comments) > 0 {
				for _, comment := range post.Comments {
					fmt.Printf("    - 评论: %s\n", comment.Content)
				}
			} else {
				fmt.Println("    - (暂无评论)")
			}
		}
	}
	fmt.Println("\n" + "--------------------------------------" + "\n")
}

func SearchDocs(db *gorm.DB) {
	// --- 2. 查询评论数量最多的文章信息 ---
	fmt.Println("--- 2. 查询评论数量最多的文章 ---")
	var mostCommentedPost Post
	// 使用子查询来找到评论数最多的文章ID
	// 或者使用 Joins, Group 和 Order
	err := db.Model(&Post{}).
		Select("posts.*, count(comments.id) as comment_count").
		Joins("left join comments on comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count desc").
		First(&mostCommentedPost).Error

	if err != nil {
		log.Printf("查询评论最多的文章失败: %v", err)
	} else {
		fmt.Printf("评论最多的文章是: '%s'\n", mostCommentedPost.Title)
		// 你可以进一步查询这篇文章的评论数量
		var count int64
		db.Model(&Comment{}).Where("post_id = ?", mostCommentedPost.ID).Count(&count)
		fmt.Printf("它共有 %d 条评论。\n", count)
	}
}

/*
继续使用博客系统的模型。

要求 ：

为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/
// --- 钩子函数定义 ---

// AfterCreate 在创建文章后触发
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Printf("[钩子触发] Post AfterCreate: 文章 '%s' 已创建, 更新用户ID %d 的文章数。\n", p.Title, p.UserID)
	// 更新对应用户的文章数量
	// 使用 gorm.Expr 来执行 SQL 表达式，实现原子更新
	return tx.Model(&User{}).Where("id = ?", p.UserID).
		UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error
}

// AfterDelete 在删除评论后触发
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Printf("[钩子触发] Comment AfterDelete: 评论ID %d 已删除, 检查文章ID %d 的评论数。\n", c.ID, c.PostID)
	var count int64
	// 统计关联文章还剩下多少评论
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		return err
	}

	fmt.Printf("[钩子检查] 文章ID %d 剩余评论数: %d\n", c.PostID, count)

	// 如果评论数量为0，则更新文章状态
	if count == 0 {
		fmt.Printf("[钩子操作] 评论数为0，更新文章ID %d 的状态为 '无评论'。\n", c.PostID)
		return tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "无评论").Error
	}
	// 如果还有评论，可以更新状态为“有评论”
	return tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "有评论").Error
}

// BeforeCreate 在创建评论前触发 (为了演示方便，在创建评论时主动更新文章状态)
func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Printf("[钩子触发] Comment BeforeCreate: 准备为文章ID %d 创建评论, 更新其状态。\n", c.PostID)
	return tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "有评论").Error
}

func test5() {
	db, err := gorm.Open(sqlite.Open("gorm_hooks_test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	// 每次运行时都重新创建表，确保环境干净
	db.Migrator().DropTable(&User{}, &Post{}, &Comment{})
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	fmt.Println("--- 1. 测试 Post 的 AfterCreate 钩子 ---")
	// 创建一个用户
	user := User{Name: "王五"}
	db.Create(&user)
	fmt.Printf("初始状态: 用户 '%s' 的文章数为 %d\n", user.Name, user.ArticleCount)

	// 创建一篇新文章，这会触发Post的AfterCreate钩子
	post := Post{Title: "探索Gorm钩子函数", Content: "钩子函数很有用...", UserID: user.ID}
	db.Create(&post)

	// 重新从数据库查询用户信息，验证 article_count 是否已更新
	var updatedUser User
	db.First(&updatedUser, user.ID)
	fmt.Printf("操作后状态: 用户 '%s' 的文章数变为 %d\n", updatedUser.Name, updatedUser.ArticleCount)
	fmt.Println("\n" + "--------------------------------------" + "\n")

	fmt.Println("--- 2. 测试 Comment 的 AfterDelete 钩子 ---")
	// 创建一条评论，触发 Comment 的 BeforeCreate 钩子，更新文章状态
	comment := Comment{Content: "这是唯一一条评论", PostID: post.ID}
	db.Create(&comment)

	// 查询文章状态
	var postWithComment Post
	db.First(&postWithComment, post.ID)
	fmt.Printf("初始状态: 文章 '%s' 的评论状态是 '%s'\n", postWithComment.Title, postWithComment.CommentStatus)

	// 删除这条评论，这会触发Comment的AfterDelete钩子
	db.Delete(&comment)

	// 再次查询文章状态，验证是否已更新为 "无评论"
	var postWithoutComment Post
	db.First(&postWithoutComment, post.ID)
	fmt.Printf("操作后状态: 文章 '%s' 的评论状态变为 '%s'\n", postWithoutComment.Title, postWithoutComment.CommentStatus)
}
